package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zodileap/taurus_go/notify"
)

const defaultBaseURL = "https://api.telegram.org"

var _ notify.Sender = (*Client)(nil)

// Option 定义 Telegram 客户端选项。
type Option func(*Client)

// Client Telegram 通知客户端。
type Client struct {
	token      string
	chatID     string
	httpClient *http.Client
	baseURL    string
}

type requestPlan struct {
	method      string
	contentType string
	body        *bytes.Buffer
}

type telegramEnvelope struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

type telegramMessage struct {
	MessageID int64 `json:"message_id"`
}

type inputMedia struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// New 创建 Telegram 通知客户端。
func New(token string, chatID string, opts ...Option) *Client {
	c := &Client{
		token:      token,
		chatID:     chatID,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    defaultBaseURL,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}

// WithHTTPClient 设置自定义 HTTP 客户端。
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithBaseURL 设置 Telegram API 基础地址。
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		baseURL = strings.TrimSpace(baseURL)
		if baseURL != "" {
			c.baseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

// Send 函数式发送 Telegram 通知。
func Send(ctx context.Context, token string, chatID string, n notify.Notification, opts ...Option) (*notify.Result, error) {
	return New(token, chatID, opts...).Send(ctx, n)
}

// Send 发送通知。
func (c *Client) Send(ctx context.Context, n notify.Notification) (*notify.Result, error) {
	text, parseMode, err := c.validate(n)
	if err != nil {
		return nil, err
	}
	plan, err := c.buildRequestPlan(n, text, parseMode)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(plan.method), bytes.NewReader(plan.body.Bytes()))
	if err != nil {
		return nil, Err_0400010002.Sprintf(err)
	}
	req.Header.Set("Content-Type", plan.contentType)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, Err_0400010004.Sprintf(err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, Err_0400010007.Sprintf(err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, Err_0400010005.Sprintf(resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	return decodeResult(raw)
}

func (c *Client) endpoint(method string) string {
	baseURL := strings.TrimRight(c.baseURL, "/")
	return fmt.Sprintf("%s/bot%s/%s", baseURL, c.token, method)
}

func (c *Client) validate(n notify.Notification) (string, string, error) {
	if strings.TrimSpace(c.token) == "" {
		return "", "", Err_0400010001.Sprintf("token is empty")
	}
	if strings.TrimSpace(c.chatID) == "" {
		return "", "", Err_0400010001.Sprintf("chat_id is empty")
	}

	text := renderText(n.Title, n.Body)
	if text == "" && len(n.Attachments) == 0 {
		return "", "", Err_0400010001.Sprintf("notification body and attachments are both empty")
	}
	parseMode, err := normalizeFormat(n.Format)
	if err != nil {
		return "", "", err
	}
	for _, attachment := range n.Attachments {
		if err := validateAttachment(attachment); err != nil {
			return "", "", err
		}
	}
	if len(n.Attachments) > 1 {
		if err := validateMediaGroup(n.Attachments); err != nil {
			return "", "", err
		}
	}
	return text, parseMode, nil
}

func (c *Client) buildRequestPlan(n notify.Notification, text string, parseMode string) (*requestPlan, error) {
	switch len(n.Attachments) {
	case 0:
		return c.buildMessagePlan(n, text, parseMode)
	case 1:
		return c.buildSingleAttachmentPlan(n, n.Attachments[0], text, parseMode)
	default:
		return c.buildMediaGroupPlan(n, text, parseMode)
	}
}

func (c *Client) buildMessagePlan(n notify.Notification, text string, parseMode string) (*requestPlan, error) {
	payload := map[string]any{
		"chat_id": c.chatID,
		"text":    text,
	}
	if parseMode != "" {
		payload["parse_mode"] = parseMode
	}
	if n.Silent {
		payload["disable_notification"] = true
	}
	if n.ProtectContent {
		payload["protect_content"] = true
	}
	return newJSONPlan("sendMessage", payload)
}

func (c *Client) buildSingleAttachmentPlan(n notify.Notification, attachment notify.Attachment, text string, parseMode string) (*requestPlan, error) {
	method, field, err := methodAndFieldForAttachment(attachment.Kind)
	if err != nil {
		return nil, err
	}
	caption := attachment.Caption
	if caption == "" {
		caption = text
	}

	if usesLocalUpload(attachment) {
		return newMultipartPlan(method, func(writer *multipart.Writer) error {
			if err := writer.WriteField("chat_id", c.chatID); err != nil {
				return Err_0400010002.Sprintf(err)
			}
			if caption != "" {
				if err := writer.WriteField("caption", caption); err != nil {
					return Err_0400010002.Sprintf(err)
				}
				if parseMode != "" {
					if err := writer.WriteField("parse_mode", parseMode); err != nil {
						return Err_0400010002.Sprintf(err)
					}
				}
			}
			if n.Silent {
				if err := writer.WriteField("disable_notification", "true"); err != nil {
					return Err_0400010002.Sprintf(err)
				}
			}
			if n.ProtectContent {
				if err := writer.WriteField("protect_content", "true"); err != nil {
					return Err_0400010002.Sprintf(err)
				}
			}
			return writeUpload(writer, field, attachment)
		})
	}

	payload := map[string]any{
		"chat_id": c.chatID,
		field:     remoteAttachmentSource(attachment),
	}
	if caption != "" {
		payload["caption"] = caption
		if parseMode != "" {
			payload["parse_mode"] = parseMode
		}
	}
	if n.Silent {
		payload["disable_notification"] = true
	}
	if n.ProtectContent {
		payload["protect_content"] = true
	}
	return newJSONPlan(method, payload)
}

func (c *Client) buildMediaGroupPlan(n notify.Notification, text string, parseMode string) (*requestPlan, error) {
	allRemote := true
	for _, attachment := range n.Attachments {
		if usesLocalUpload(attachment) {
			allRemote = false
			break
		}
	}
	if allRemote {
		media := make([]inputMedia, 0, len(n.Attachments))
		for idx, attachment := range n.Attachments {
			item := inputMedia{
				Type:  string(attachment.Kind),
				Media: remoteAttachmentSource(attachment),
			}
			caption := mediaGroupCaption(idx, text, attachment.Caption)
			if caption != "" {
				item.Caption = caption
				if parseMode != "" {
					item.ParseMode = parseMode
				}
			}
			media = append(media, item)
		}
		payload := map[string]any{
			"chat_id": c.chatID,
			"media":   media,
		}
		if n.Silent {
			payload["disable_notification"] = true
		}
		if n.ProtectContent {
			payload["protect_content"] = true
		}
		return newJSONPlan("sendMediaGroup", payload)
	}

	return newMultipartPlan("sendMediaGroup", func(writer *multipart.Writer) error {
		if err := writer.WriteField("chat_id", c.chatID); err != nil {
			return Err_0400010002.Sprintf(err)
		}
		if n.Silent {
			if err := writer.WriteField("disable_notification", "true"); err != nil {
				return Err_0400010002.Sprintf(err)
			}
		}
		if n.ProtectContent {
			if err := writer.WriteField("protect_content", "true"); err != nil {
				return Err_0400010002.Sprintf(err)
			}
		}

		media := make([]inputMedia, 0, len(n.Attachments))
		for idx, attachment := range n.Attachments {
			item := inputMedia{Type: string(attachment.Kind)}
			if usesLocalUpload(attachment) {
				item.Media = "attach://" + mediaFieldName(idx)
			} else {
				item.Media = remoteAttachmentSource(attachment)
			}
			caption := mediaGroupCaption(idx, text, attachment.Caption)
			if caption != "" {
				item.Caption = caption
				if parseMode != "" {
					item.ParseMode = parseMode
				}
			}
			media = append(media, item)
		}

		data, err := json.Marshal(media)
		if err != nil {
			return Err_0400010002.Sprintf(err)
		}
		if err := writer.WriteField("media", string(data)); err != nil {
			return Err_0400010002.Sprintf(err)
		}
		for idx, attachment := range n.Attachments {
			if !usesLocalUpload(attachment) {
				continue
			}
			if err := writeUpload(writer, mediaFieldName(idx), attachment); err != nil {
				return err
			}
		}
		return nil
	})
}

func newJSONPlan(method string, payload any) (*requestPlan, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, Err_0400010002.Sprintf(err)
	}
	return &requestPlan{
		method:      method,
		contentType: "application/json",
		body:        bytes.NewBuffer(data),
	}, nil
}

func newMultipartPlan(method string, build func(writer *multipart.Writer) error) (*requestPlan, error) {
	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	if err := build(writer); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, Err_0400010002.Sprintf(err)
	}
	return &requestPlan{
		method:      method,
		contentType: writer.FormDataContentType(),
		body:        body,
	}, nil
}

func decodeResult(raw []byte) (*notify.Result, error) {
	var envelope telegramEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, Err_0400010007.Sprintf(err)
	}
	if !envelope.OK {
		return nil, Err_0400010006.Sprintf(envelope.ErrorCode, envelope.Description)
	}

	data := bytes.TrimSpace(envelope.Result)
	if len(data) == 0 {
		return nil, Err_0400010007.Sprintf("missing result")
	}

	var ids []string
	switch data[0] {
	case '{':
		var message telegramMessage
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, Err_0400010007.Sprintf(err)
		}
		ids = []string{strconv.FormatInt(message.MessageID, 10)}
	case '[':
		var messages []telegramMessage
		if err := json.Unmarshal(data, &messages); err != nil {
			return nil, Err_0400010007.Sprintf(err)
		}
		ids = make([]string, 0, len(messages))
		for _, message := range messages {
			ids = append(ids, strconv.FormatInt(message.MessageID, 10))
		}
	default:
		return nil, Err_0400010007.Sprintf("unexpected result shape")
	}
	return &notify.Result{
		Provider:   "telegram",
		MessageIDs: ids,
		Raw:        raw,
	}, nil
}

func renderText(title string, body string) string {
	switch {
	case title != "" && body != "":
		return title + "\n\n" + body
	case title != "":
		return title
	default:
		return body
	}
}

func normalizeFormat(format notify.Format) (string, error) {
	switch format {
	case notify.Plain:
		return "", nil
	case notify.MarkdownV2:
		return string(notify.MarkdownV2), nil
	case notify.HTML:
		return string(notify.HTML), nil
	default:
		return "", Err_0400010001.Sprintf(fmt.Sprintf("unsupported format %q", format))
	}
}

func validateAttachment(attachment notify.Attachment) error {
	switch attachment.Kind {
	case notify.Photo, notify.Video, notify.Audio, notify.Voice, notify.Document, notify.Animation:
	default:
		return Err_0400010001.Sprintf(fmt.Sprintf("unsupported attachment kind %q", attachment.Kind))
	}

	sourceCount := 0
	if attachment.FileID != "" {
		sourceCount++
	}
	if attachment.URL != "" {
		sourceCount++
	}
	if attachment.Path != "" {
		sourceCount++
	}
	if attachment.Content != nil {
		sourceCount++
	}
	if sourceCount != 1 {
		return Err_0400010001.Sprintf(fmt.Sprintf("attachment %q must provide exactly one source", attachment.Kind))
	}
	if attachment.Content != nil && attachment.Name == "" {
		return Err_0400010001.Sprintf(fmt.Sprintf("attachment %q requires name when content is provided", attachment.Kind))
	}
	return nil
}

func validateMediaGroup(attachments []notify.Attachment) error {
	kindSet := make(map[notify.AttachmentKind]struct{})
	for _, attachment := range attachments {
		switch attachment.Kind {
		case notify.Voice, notify.Animation:
			return Err_0400010008.Sprintf(fmt.Sprintf("%s does not support media groups", attachment.Kind))
		}
		kindSet[attachment.Kind] = struct{}{}
	}
	if len(kindSet) == 1 {
		if _, ok := kindSet[notify.Document]; ok {
			return nil
		}
		if _, ok := kindSet[notify.Audio]; ok {
			return nil
		}
		if _, ok := kindSet[notify.Photo]; ok {
			return nil
		}
		if _, ok := kindSet[notify.Video]; ok {
			return nil
		}
	}
	if len(kindSet) == 2 {
		_, hasPhoto := kindSet[notify.Photo]
		_, hasVideo := kindSet[notify.Video]
		if hasPhoto && hasVideo {
			return nil
		}
	}
	return Err_0400010008.Sprintf(strings.Join(mediaGroupKinds(kindSet), ","))
}

func mediaGroupKinds(kindSet map[notify.AttachmentKind]struct{}) []string {
	kinds := make([]string, 0, len(kindSet))
	for kind := range kindSet {
		kinds = append(kinds, string(kind))
	}
	sort.Strings(kinds)
	return kinds
}

func methodAndFieldForAttachment(kind notify.AttachmentKind) (string, string, error) {
	switch kind {
	case notify.Photo:
		return "sendPhoto", "photo", nil
	case notify.Video:
		return "sendVideo", "video", nil
	case notify.Audio:
		return "sendAudio", "audio", nil
	case notify.Voice:
		return "sendVoice", "voice", nil
	case notify.Document:
		return "sendDocument", "document", nil
	case notify.Animation:
		return "sendAnimation", "animation", nil
	default:
		return "", "", Err_0400010001.Sprintf(fmt.Sprintf("unsupported attachment kind %q", kind))
	}
}

func mediaGroupCaption(index int, text string, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if index == 0 {
		return text
	}
	return ""
}

func usesLocalUpload(attachment notify.Attachment) bool {
	return attachment.Path != "" || attachment.Content != nil
}

func remoteAttachmentSource(attachment notify.Attachment) string {
	if attachment.FileID != "" {
		return attachment.FileID
	}
	return attachment.URL
}

func mediaFieldName(index int) string {
	return fmt.Sprintf("file%d", index)
}

func writeUpload(writer *multipart.Writer, field string, attachment notify.Attachment) error {
	fileName := attachmentFileName(attachment)
	part, err := createUploadPart(writer, field, fileName, attachment.MIME)
	if err != nil {
		return Err_0400010002.Sprintf(err)
	}
	if attachment.Path != "" {
		file, err := os.Open(attachment.Path)
		if err != nil {
			return Err_0400010003.Sprintf(attachment.Path, err)
		}
		defer file.Close()
		if _, err := io.Copy(part, file); err != nil {
			return Err_0400010003.Sprintf(attachment.Path, err)
		}
		return nil
	}
	if _, err := part.Write(attachment.Content); err != nil {
		return Err_0400010002.Sprintf(err)
	}
	return nil
}

func attachmentFileName(attachment notify.Attachment) string {
	if attachment.Name != "" {
		return attachment.Name
	}
	return filepath.Base(attachment.Path)
}

func createUploadPart(writer *multipart.Writer, field string, fileName string, mimeType string) (io.Writer, error) {
	if mimeType == "" {
		return writer.CreateFormFile(field, fileName)
	}
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q; filename=%q`, field, fileName))
	header.Set("Content-Type", mimeType)
	return writer.CreatePart(header)
}
