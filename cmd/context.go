/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/IrwantoCia/gomakase/embed"
	"github.com/IrwantoCia/gomakase/internal/ctx_context/application"
	"github.com/IrwantoCia/gomakase/internal/shared/command"
	"github.com/IrwantoCia/gomakase/internal/shared/config"
	"github.com/IrwantoCia/gomakase/internal/shared/file"
	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Generate a new context",
	Long: `Generate a new context. For example:

gomakase context <context_name>
`,
	Args:    cobra.ExactArgs(1),
	Example: `gomakase context <context_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]

		// read the root config file
		rootConfigFile, err := os.ReadFile("gen.yaml")
		if err != nil {
			log.Fatalf("Error reading root config file: %v", err)
		}
		rootConfig, err := config.LoadSchematic[config.RootSchematic](rootConfigFile)
		if err != nil {
			log.Fatalf("Error loading root config: %v", err)
		}

		contextConfigFile := filepath.Join("schematics", "context", "schematic.yaml")
		contextConfigFileContent, err := embed.SchematicsFS.ReadFile(contextConfigFile)
		if err != nil {
			log.Fatalf("Error reading context config file: %v", err)
		}
		contextConfig, err := config.LoadSchematic[config.ContextSchematic](contextConfigFileContent)
		if err != nil {
			log.Fatalf("Error loading context config: %v", err)
		}

		file := file.NewFile()
		contextService := application.NewCtxService(file, rootConfig, contextConfig)
		err = contextService.Generate(contextName)
		if err != nil {
			log.Fatalf("Error generating context: %v", err)
		}

		// after commands
		command := command.NewCommand()
		err = command.GoModTidy()
		if err != nil {
			log.Fatalf("Error running go mod tidy: %v", err)
		}
		err = command.GoFmt()
		if err != nil {
			log.Fatalf("Error running go fmt: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
