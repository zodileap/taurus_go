package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zodileap/taurus_go/entity/cmd/internal"
)

func main() {
	log.SetFlags(0)
	cmd := &cobra.Command{Use: "github.com/zodileap/taurus_go/entity/cmd"}
	cmd.AddCommand(
		internal.GenerateCmd(),
		internal.NewCmd(),
	)
	_ = cmd.Execute()
}
