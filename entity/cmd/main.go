package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yohobala/taurus_go/entity/cmd/internal"
)

func main() {
	log.SetFlags(0)
	cmd := &cobra.Command{Use: "github.com/yohobala/taurus_go/entity/cmd"}
	cmd.AddCommand(
		internal.GenerateCmd(),
	)
	_ = cmd.Execute()
}
