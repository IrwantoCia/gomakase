/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path"

	"github.com/IrwantoCia/gomakase/embed"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available plugins",
	Long: `List all available plugins that can be used with the 'add' command.
Each plugin represents a schematic that can be added to your project.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available plugins:")

		pluginsDir := path.Join("schematics", "plugins")
		entries, err := embed.SchematicsFS.ReadDir(pluginsDir)
		if err != nil {
			log.Fatalf("Error reading schematics directory: %v", err)
		}

		if len(entries) == 0 {
			fmt.Println("No plugins found.")
			return
		}

		for _, entry := range entries {
			pluginName := entry.Name()
			fmt.Printf("  • %s\n", pluginName)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
