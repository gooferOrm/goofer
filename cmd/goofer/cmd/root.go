package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	version = "0.1.0"
	// verbose output flag
	verbose bool
	// configFile path
	configFile string
)

var rootCmd = &cobra.Command{
	Use:     "goofer",
	Short:   "A Type-safe ORM for Go",
	Long:    "A powerful, type-safe ORM for Go that provides an Amazing developer experience",
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file (default is ./goofer.yaml)")

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number of Goofer ORM CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Goofer ORM v%s\n", version)
		},
	})
}

// printVerbose prints a message if verbose mode is enabled
func printVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
