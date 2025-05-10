package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goofer",
	Short: "A Type-safe ORM for Go",
	Long:  "A powerful, type-safe ORM for Go that provides an Amazing developer experience",
}

func Execute() error {
	return rootCmd.Execute()
}
