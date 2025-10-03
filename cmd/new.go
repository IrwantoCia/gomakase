/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"path/filepath"

	"github.com/IrwantoCia/gomakase/embed"
	"github.com/IrwantoCia/gomakase/internal/new_context/application"
	"github.com/IrwantoCia/gomakase/internal/shared/command"
	"github.com/IrwantoCia/gomakase/internal/shared/config"
	"github.com/IrwantoCia/gomakase/internal/shared/file"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new <project_name>",
	Short:   "Create a new project",
	Args:    cobra.ExactArgs(1),
	Example: `gomakase new <project_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		projectConfigFile := filepath.Join("schematics", "project", "schematic.yaml")
		projectConfigFileContent, err := embed.SchematicsFS.ReadFile(projectConfigFile)
		if err != nil {
			log.Fatalf("Error reading project config file: %v", err)
		}

		projectSchematic, err := config.LoadSchematic[config.ProjectSchematic](projectConfigFileContent)
		if err != nil {
			log.Fatalf("Error loading project config: %v", err)
		}

		if len(projectSchematic.Actions) == 0 {
			log.Fatalf("No actions found in project schematic")
		}

		file := file.NewFile()
		newService := application.NewNewService(file)
		newService.Generate(projectName, projectSchematic)

		// after commands
		command := command.NewCommand()
		err = command.ChangeFolder(projectName)
		if err != nil {
			log.Fatalf("Error changing folder: %v", err)
		}
		err = command.GoModTidy()
		if err != nil {
			log.Fatalf("Error running go mod tidy: %v", err)
		}
		err = command.NPMInstall()
		if err != nil {
			log.Fatalf("Error running npm install: %v", err)
		}
		err = command.GoFmt()
		if err != nil {
			log.Fatalf("Error running go fmt: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
