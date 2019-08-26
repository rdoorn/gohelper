package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func MustGetString(cobra *cobra.Command, name string) string {
	s, err := cobra.Flags().GetString(name)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	return s
}

func MustGetInt(cobra *cobra.Command, name string) int {
	s, err := cobra.Flags().GetInt(name)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	return s
}
