/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path"

	"github.com/IrwantoCia/gomakase/embed"
	"github.com/IrwantoCia/gomakase/internal/add_context/application"
	"github.com/IrwantoCia/gomakase/internal/shared/command"
	"github.com/IrwantoCia/gomakase/internal/shared/config"
	"github.com/IrwantoCia/gomakase/internal/shared/file"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add <plugin_name>",
	Short:   "Add a new plugin to the project",
	Long:    `Add a new plugin to the project. For available plugins, see the gomakase list command.`,
	Args:    cobra.ExactArgs(1),
	Example: `gomakase add <plugin_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		pluginList, err := embed.SchematicsFS.ReadDir(path.Join("schematics", "plugins"))
		if err != nil {
			log.Fatalf("Error reading schematics directory: %v", err)
		}
		if len(pluginList) == 0 {
			log.Fatalf("No plugins found.")
		}
		selectedPlugin := ""
		for _, plugin := range pluginList {
			if plugin.Name() == pluginName {
				selectedPlugin = plugin.Name()
				break
			}
		}

		if selectedPlugin == "" {
			log.Fatalf("Plugin not found.")
		}

		log.Printf("Adding plugin: %s", selectedPlugin)

		// read the root config file
		rootConfigFile, err := os.ReadFile("gen.yaml")
		if err != nil {
			log.Fatalf("Error reading root config file: %v", err)
		}
		rootConfig, err := config.LoadSchematic[config.RootSchematic](rootConfigFile)
		if err != nil {
			log.Fatalf("Error loading root config: %v", err)
		}

		// read the plugin config file
		pluginConfigFile, err := embed.SchematicsFS.ReadFile(
			path.Join(
				"schematics",
				"plugins",
				selectedPlugin,
				"schematic.yaml",
			),
		)
		if err != nil {
			log.Fatalf("Error reading plugin config file: %v", err)
		}
		pluginConfig, err := config.LoadSchematic[config.PluginSchematic](pluginConfigFile)
		if err != nil {
			log.Fatalf("Error loading plugin config: %v", err)
		}

		file := file.NewFile()

		addService := application.NewAddService(
			rootConfig,
			pluginConfig,
			embed.SchematicsFS,
			file,
		)
		addService.Generate(pluginName)

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
		err = command.NPMInstall()
		if err != nil {
			log.Fatalf("Error running npm install: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
