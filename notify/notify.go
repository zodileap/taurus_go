package notify

import "context"

// Format 定义消息的文本格式。
type Format string

const (
	// Plain 纯文本。
	Plain Format = ""
	// MarkdownV2 Telegram MarkdownV2 格式。
	MarkdownV2 Format = "MarkdownV2"
	// HTML Telegram HTML 格式。
	HTML Format = "HTML"
)

// AttachmentKind 定义通知附件类型。
type AttachmentKind string

const (
	Photo     AttachmentKind = "photo"
	Video     AttachmentKind = "video"
	Audio     AttachmentKind = "audio"
	Voice     AttachmentKind = "voice"
	Document  AttachmentKind = "document"
	Animation AttachmentKind = "animation"
)

// Notification 定义统一通知模型。
type Notification struct {
	Title          string
	Body           string
	Format         Format
	Silent         bool
	ProtectContent bool
	Attachments    []Attachment
}

// Attachment 定义通知附件。
type Attachment struct {
	Kind    AttachmentKind
	Name    string
	MIME    string
	Caption string
	FileID  string
	URL     string
	Path    string
	Content []byte
}

// Result 定义发送结果。
type Result struct {
	Provider   string
	MessageIDs []string
	Raw        []byte
}

// Sender 定义统一发送接口。
type Sender interface {
	Send(ctx context.Context, n Notification) (*Result, error)
}
