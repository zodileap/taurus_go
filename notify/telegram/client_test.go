package telegram

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zodileap/taurus_go/notify"
)

func TestSendMessageFunctionJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/bottest-token/sendMessage")
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			t.Fatalf("expected application/json, got %s", got)
		}
		payload := decodeJSONMap(t, r.Body)
		assertStringValue(t, payload, "chat_id", "12345")
		assertStringValue(t, payload, "text", "Build Ready\n\nall checks passed")
		assertStringValue(t, payload, "parse_mode", "HTML")
		assertBoolValue(t, payload, "disable_notification", true)
		assertBoolValue(t, payload, "protect_content", true)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":101}}`))
	}))
	defer server.Close()

	result, err := Send(
		context.Background(),
		"test-token",
		"12345",
		notify.Notification{
			Title:          "Build Ready",
			Body:           "all checks passed",
			Format:         notify.HTML,
			Silent:         true,
			ProtectContent: true,
		},
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	assertResult(t, result, []string{"101"})
}

func TestClientSendSingleAttachments(t *testing.T) {
	tempDir := t.TempDir()
	videoPath := filepath.Join(tempDir, "sample.mp4")
	if err := os.WriteFile(videoPath, []byte("video-bytes"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	tests := []struct {
		name         string
		methodPath   string
		notification notify.Notification
		assert       func(t *testing.T, r *http.Request)
	}{
		{
			name:       "photo content multipart",
			methodPath: "/bottest-token/sendPhoto",
			notification: notify.Notification{
				Body:   "photo body",
				Format: notify.MarkdownV2,
				Attachments: []notify.Attachment{{
					Kind:    notify.Photo,
					Name:    "photo.jpg",
					MIME:    "image/jpeg",
					Content: []byte("photo-bytes"),
				}},
			},
			assert: func(t *testing.T, r *http.Request) {
				form := parseMultipartForm(t, r)
				assertString(t, r.FormValue("chat_id"), "12345")
				assertString(t, r.FormValue("caption"), "photo body")
				assertString(t, r.FormValue("parse_mode"), "MarkdownV2")
				assertUploadedFile(t, form, "photo", "photo-bytes")
			},
		},
		{
			name:       "video path multipart",
			methodPath: "/bottest-token/sendVideo",
			notification: notify.Notification{
				Title: "video title",
				Attachments: []notify.Attachment{{
					Kind: notify.Video,
					Path: videoPath,
				}},
			},
			assert: func(t *testing.T, r *http.Request) {
				form := parseMultipartForm(t, r)
				assertString(t, r.FormValue("chat_id"), "12345")
				assertString(t, r.FormValue("caption"), "video title")
				assertUploadedFile(t, form, "video", "video-bytes")
			},
		},
		{
			name:       "audio url json",
			methodPath: "/bottest-token/sendAudio",
			notification: notify.Notification{
				Body: "audio body",
				Attachments: []notify.Attachment{{
					Kind: notify.Audio,
					URL:  "https://example.com/audio.mp3",
				}},
			},
			assert: func(t *testing.T, r *http.Request) {
				payload := decodeJSONMap(t, r.Body)
				assertStringValue(t, payload, "chat_id", "12345")
				assertStringValue(t, payload, "audio", "https://example.com/audio.mp3")
				assertStringValue(t, payload, "caption", "audio body")
			},
		},
		{
			name:       "voice file id json",
			methodPath: "/bottest-token/sendVoice",
			notification: notify.Notification{
				Body: "voice body",
				Attachments: []notify.Attachment{{
					Kind:   notify.Voice,
					FileID: "voice-file-id",
				}},
			},
			assert: func(t *testing.T, r *http.Request) {
				payload := decodeJSONMap(t, r.Body)
				assertStringValue(t, payload, "chat_id", "12345")
				assertStringValue(t, payload, "voice", "voice-file-id")
				assertStringValue(t, payload, "caption", "voice body")
			},
		},
		{
			name:       "document content explicit caption",
			methodPath: "/bottest-token/sendDocument",
			notification: notify.Notification{
				Body: "document body",
				Attachments: []notify.Attachment{{
					Kind:    notify.Document,
					Name:    "report.pdf",
					Content: []byte("pdf-bytes"),
					Caption: "custom caption",
				}},
			},
			assert: func(t *testing.T, r *http.Request) {
				form := parseMultipartForm(t, r)
				assertString(t, r.FormValue("caption"), "custom caption")
				assertUploadedFile(t, form, "document", "pdf-bytes")
			},
		},
		{
			name:       "animation url json",
			methodPath: "/bottest-token/sendAnimation",
			notification: notify.Notification{
				Body:           "<b>alert</b>",
				Format:         notify.HTML,
				Silent:         true,
				ProtectContent: true,
				Attachments: []notify.Attachment{{
					Kind: notify.Animation,
					URL:  "https://example.com/anim.gif",
				}},
			},
			assert: func(t *testing.T, r *http.Request) {
				payload := decodeJSONMap(t, r.Body)
				assertStringValue(t, payload, "animation", "https://example.com/anim.gif")
				assertStringValue(t, payload, "caption", "<b>alert</b>")
				assertStringValue(t, payload, "parse_mode", "HTML")
				assertBoolValue(t, payload, "disable_notification", true)
				assertBoolValue(t, payload, "protect_content", true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertPath(t, r, tt.methodPath)
				tt.assert(t, r)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":202}}`))
			}))
			defer server.Close()

			client := New("test-token", "12345", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
			result, err := client.Send(context.Background(), tt.notification)
			if err != nil {
				t.Fatalf("Client.Send() error = %v", err)
			}
			assertResult(t, result, []string{"202"})
		})
	}
}

func TestClientSendMediaGroups(t *testing.T) {
	tests := []struct {
		name         string
		notification notify.Notification
		assert       func(t *testing.T, r *http.Request)
	}{
		{
			name: "photo video multipart",
			notification: notify.Notification{
				Title:  "Deploy",
				Body:   "photo and video",
				Format: notify.MarkdownV2,
				Attachments: []notify.Attachment{
					{
						Kind:    notify.Photo,
						Name:    "cover.jpg",
						Content: []byte("cover-bytes"),
					},
					{
						Kind:    notify.Video,
						URL:     "https://example.com/demo.mp4",
						Caption: "clip",
					},
				},
			},
			assert: func(t *testing.T, r *http.Request) {
				form := parseMultipartForm(t, r)
				assertString(t, r.FormValue("chat_id"), "12345")
				media := decodeMediaField(t, r.FormValue("media"))
				if len(media) != 2 {
					t.Fatalf("expected 2 media items, got %d", len(media))
				}
				assertString(t, media[0]["type"].(string), "photo")
				assertString(t, media[0]["media"].(string), "attach://file0")
				assertString(t, media[0]["caption"].(string), "Deploy\n\nphoto and video")
				assertString(t, media[0]["parse_mode"].(string), "MarkdownV2")
				assertString(t, media[1]["type"].(string), "video")
				assertString(t, media[1]["media"].(string), "https://example.com/demo.mp4")
				assertString(t, media[1]["caption"].(string), "clip")
				assertString(t, media[1]["parse_mode"].(string), "MarkdownV2")
				assertUploadedFile(t, form, "file0", "cover-bytes")
			},
		},
		{
			name: "document json",
			notification: notify.Notification{
				Body: "docs",
				Attachments: []notify.Attachment{
					{Kind: notify.Document, FileID: "doc-1"},
					{Kind: notify.Document, URL: "https://example.com/doc-2.pdf"},
				},
			},
			assert: func(t *testing.T, r *http.Request) {
				payload := decodeJSONMap(t, r.Body)
				assertStringValue(t, payload, "chat_id", "12345")
				media := payload["media"].([]any)
				if len(media) != 2 {
					t.Fatalf("expected 2 media items, got %d", len(media))
				}
				first := media[0].(map[string]any)
				second := media[1].(map[string]any)
				assertStringValue(t, first, "type", "document")
				assertStringValue(t, first, "media", "doc-1")
				assertStringValue(t, first, "caption", "docs")
				assertStringValue(t, second, "type", "document")
				assertStringValue(t, second, "media", "https://example.com/doc-2.pdf")
			},
		},
		{
			name: "audio json",
			notification: notify.Notification{
				Body: "audios",
				Attachments: []notify.Attachment{
					{Kind: notify.Audio, FileID: "audio-1"},
					{Kind: notify.Audio, FileID: "audio-2"},
				},
			},
			assert: func(t *testing.T, r *http.Request) {
				payload := decodeJSONMap(t, r.Body)
				media := payload["media"].([]any)
				if len(media) != 2 {
					t.Fatalf("expected 2 media items, got %d", len(media))
				}
				first := media[0].(map[string]any)
				second := media[1].(map[string]any)
				assertStringValue(t, first, "type", "audio")
				assertStringValue(t, first, "media", "audio-1")
				assertStringValue(t, first, "caption", "audios")
				assertStringValue(t, second, "type", "audio")
				assertStringValue(t, second, "media", "audio-2")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertPath(t, r, "/bottest-token/sendMediaGroup")
				tt.assert(t, r)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":[{"message_id":301},{"message_id":302}]}`))
			}))
			defer server.Close()

			client := New("test-token", "12345", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
			result, err := client.Send(context.Background(), tt.notification)
			if err != nil {
				t.Fatalf("Client.Send() error = %v", err)
			}
			assertResult(t, result, []string{"301", "302"})
		})
	}
}

func TestSendValidationErrors(t *testing.T) {
	tests := []struct {
		name         string
		client       *Client
		notification notify.Notification
		code         string
	}{
		{
			name:         "empty token",
			client:       New("", "12345"),
			notification: notify.Notification{Body: "hello"},
			code:         "0400010001",
		},
		{
			name:         "empty chat id",
			client:       New("token", ""),
			notification: notify.Notification{Body: "hello"},
			code:         "0400010001",
		},
		{
			name:         "empty notification",
			client:       New("token", "12345"),
			notification: notify.Notification{},
			code:         "0400010001",
		},
		{
			name:   "attachment source conflict",
			client: New("token", "12345"),
			notification: notify.Notification{Attachments: []notify.Attachment{{
				Kind:   notify.Photo,
				URL:    "https://example.com/a.jpg",
				FileID: "file-id",
			}}},
			code: "0400010001",
		},
		{
			name:   "content missing name",
			client: New("token", "12345"),
			notification: notify.Notification{Attachments: []notify.Attachment{{
				Kind:    notify.Document,
				Content: []byte("doc"),
			}}},
			code: "0400010001",
		},
		{
			name:   "voice in media group",
			client: New("token", "12345"),
			notification: notify.Notification{Attachments: []notify.Attachment{
				{Kind: notify.Voice, FileID: "voice-1"},
				{Kind: notify.Voice, FileID: "voice-2"},
			}},
			code: "0400010008",
		},
		{
			name:   "document mixed with photo",
			client: New("token", "12345"),
			notification: notify.Notification{Attachments: []notify.Attachment{
				{Kind: notify.Document, FileID: "doc-1"},
				{Kind: notify.Photo, FileID: "photo-1"},
			}},
			code: "0400010008",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.client.Send(context.Background(), tt.notification)
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), "code:"+tt.code) {
				t.Fatalf("expected error code %s, got %v", tt.code, err)
			}
		})
	}
}

func TestSendTransportAndResponseErrors(t *testing.T) {
	t.Run("http status error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "bad gateway", http.StatusBadGateway)
		}))
		defer server.Close()

		client := New("token", "12345", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.Send(context.Background(), notify.Notification{Body: "hello"})
		if err == nil || !strings.Contains(err.Error(), "code:0400010005") {
			t.Fatalf("expected http status error, got %v", err)
		}
	})

	t.Run("telegram api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()

		client := New("token", "12345", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.Send(context.Background(), notify.Notification{Body: "hello"})
		if err == nil || !strings.Contains(err.Error(), "code:0400010006") {
			t.Fatalf("expected api error, got %v", err)
		}
	})

	t.Run("invalid json response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		client := New("token", "12345", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.Send(context.Background(), notify.Notification{Body: "hello"})
		if err == nil || !strings.Contains(err.Error(), "code:0400010007") {
			t.Fatalf("expected decode error, got %v", err)
		}
	})

	t.Run("transport error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		client := New("token", "12345", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		server.Close()

		_, err := client.Send(context.Background(), notify.Notification{Body: "hello"})
		if err == nil || !strings.Contains(err.Error(), "code:0400010004") {
			t.Fatalf("expected transport error, got %v", err)
		}
	})
}

func assertPath(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.URL.Path; got != want {
		t.Fatalf("expected path %s, got %s", want, got)
	}
}

func decodeJSONMap(t *testing.T, body io.Reader) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	return payload
}

func parseMultipartForm(t *testing.T, r *http.Request) *multipart.Form {
	t.Helper()
	if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
		t.Fatalf("expected multipart/form-data, got %s", got)
	}
	if err := r.ParseMultipartForm(16 << 20); err != nil {
		t.Fatalf("parse multipart form: %v", err)
	}
	return r.MultipartForm
}

func assertUploadedFile(t *testing.T, form *multipart.Form, field string, want string) {
	t.Helper()
	files := form.File[field]
	if len(files) != 1 {
		t.Fatalf("expected 1 uploaded file for %s, got %d", field, len(files))
	}
	file, err := files[0].Open()
	if err != nil {
		t.Fatalf("open uploaded file: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("read uploaded file: %v", err)
	}
	assertString(t, string(data), want)
}

func decodeMediaField(t *testing.T, value string) []map[string]any {
	t.Helper()
	var media []map[string]any
	if err := json.Unmarshal([]byte(value), &media); err != nil {
		t.Fatalf("decode media: %v", err)
	}
	return media
}

func assertResult(t *testing.T, result *notify.Result, wantIDs []string) {
	t.Helper()
	if result == nil {
		t.Fatalf("result is nil")
	}
	assertString(t, result.Provider, "telegram")
	if len(result.MessageIDs) != len(wantIDs) {
		t.Fatalf("expected %d ids, got %d", len(wantIDs), len(result.MessageIDs))
	}
	for idx, want := range wantIDs {
		assertString(t, result.MessageIDs[idx], want)
	}
	if len(result.Raw) == 0 {
		t.Fatalf("expected raw response")
	}
}

func assertStringValue(t *testing.T, payload map[string]any, key string, want string) {
	t.Helper()
	value, ok := payload[key]
	if !ok {
		t.Fatalf("expected key %s", key)
	}
	got, ok := value.(string)
	if !ok {
		t.Fatalf("expected string value for %s, got %T", key, value)
	}
	assertString(t, got, want)
}

func assertBoolValue(t *testing.T, payload map[string]any, key string, want bool) {
	t.Helper()
	value, ok := payload[key]
	if !ok {
		t.Fatalf("expected key %s", key)
	}
	got, ok := value.(bool)
	if !ok {
		t.Fatalf("expected bool value for %s, got %T", key, value)
	}
	if got != want {
		t.Fatalf("expected %s=%v, got %v", key, want, got)
	}
}

func assertString(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
