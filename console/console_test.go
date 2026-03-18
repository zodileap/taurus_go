package console

import (
	"io"
	"os"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	return string(output)
}

func TestExamples(t *testing.T) {
	got := Examples("go test ./...", "go vet ./...")
	want := "  go test ./...\n  go vet ./..."
	if got != want {
		t.Fatalf("unexpected examples output: %q", got)
	}
}

func TestModule(t *testing.T) {
	got := captureStdout(t, func() {
		Module("build %s", "api")
	})
	if got != "*** build api ***\n" {
		t.Fatalf("unexpected module output: %q", got)
	}
}

func TestStep(t *testing.T) {
	got := captureStdout(t, func() {
		Step("download config")
	})
	if got != "-> download config\n" {
		t.Fatalf("unexpected step output: %q", got)
	}
}

func TestSubStep(t *testing.T) {
	got := captureStdout(t, func() {
		SubStep("parse env")
	})
	if got != "   -> parse env\n" {
		t.Fatalf("unexpected substep output: %q", got)
	}
}

func TestSkip(t *testing.T) {
	got := captureStdout(t, func() {
		Skip("config file missing")
	})
	if got != "   ! config file missing [skip]\n" {
		t.Fatalf("unexpected skip output: %q", got)
	}
}

func TestDone(t *testing.T) {
	got := captureStdout(t, func() {
		Done("deploy completed")
	})
	if got != "[done] deploy completed\n" {
		t.Fatalf("unexpected done output: %q", got)
	}
}
