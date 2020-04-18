package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootFlags = struct {
}{}

var rootCmd = &cobra.Command{
	Use:   "patrix",
	Short: "PATRiX - transmit and receive digital modes through the Pulse Audio framework",
}

// Execute is called by main.main() as the entry point to the Cobra framework.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
}
