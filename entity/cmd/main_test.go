package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/zodileap/taurus_go/entity/cmd/internal"
)

func TestRootCommandBuilds(t *testing.T) {
	cmd := &cobra.Command{Use: "github.com/zodileap/taurus_go/entity/cmd"}
	cmd.AddCommand(
		internal.GenerateCmd(),
		internal.NewCmd(),
	)

	if len(cmd.Commands()) != 2 {
		t.Fatalf("根命令子命令数量不正确: %d", len(cmd.Commands()))
	}
}
