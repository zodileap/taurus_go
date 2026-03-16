package notify

import (
	"context"
	"testing"
)

type fakeSender struct {
	got Notification
}

func (f *fakeSender) Send(_ context.Context, n Notification) (*Result, error) {
	f.got = n
	return &Result{
		Provider:   "fake",
		MessageIDs: []string{"1"},
		Raw:        []byte("ok"),
	}, nil
}

func TestNotificationSenderContract(t *testing.T) {
	var sender Sender = &fakeSender{}

	n := Notification{
		Title:          "title",
		Body:           "body",
		Silent:         true,
		ProtectContent: true,
		Attachments: []Attachment{
			{
				Kind:    Document,
				Name:    "file.txt",
				MIME:    "text/plain",
				Content: []byte("payload"),
			},
		},
	}

	res, err := sender.Send(context.Background(), n)
	if err != nil {
		t.Fatalf("发送通知失败: %v", err)
	}
	if res.Provider != "fake" || len(res.MessageIDs) != 1 || string(res.Raw) != "ok" {
		t.Fatalf("发送结果不正确: %+v", res)
	}

	if n.Format != Plain {
		t.Fatalf("零值格式应等于 Plain，实际为 %q", n.Format)
	}
}
