package application

import (
	"embed"
	"log"
	"path"

	"github.com/IrwantoCia/gomakase/internal/shared/config"
	"github.com/IrwantoCia/gomakase/internal/shared/file"
	"github.com/IrwantoCia/gomakase/internal/shared/parser"
)

type AddService interface {
	Generate(contextName string) error
}

type addService struct {
	EmbedFS      embed.FS
	RootConfig   config.RootSchematic
	PluginConfig config.PluginSchematic
	File         file.File
}

func NewAddService(
	rootConfig config.RootSchematic,
	pluginConfig config.PluginSchematic,
	embedFS embed.FS,
	file file.File,
) AddService {
	return &addService{
		EmbedFS:      embedFS,
		RootConfig:   rootConfig,
		PluginConfig: pluginConfig,
		File:         file,
	}
}

type ImportAction struct {
	OutputPath string
	ImportPath string
	Alias      string
}

type CreateFileAction struct {
	OutputPath string
	Content    []byte
}

type DependencyAction struct {
	OutputPath string
	Dependency string
}

type RouteAction struct {
	OutputPath string
	Route      string
}

type Job struct {
	Type             string
	CreateFileAction *CreateFileAction
	ImportAction     *ImportAction
	DependencyAction *DependencyAction
	RouteAction      *RouteAction
}

func (s *addService) Generate(contextName string) error {
	log.Printf("Adding %s", contextName)

	if s.File.IsPathExists(contextName) {
		log.Printf("Plugin is already exists, skipping...\n")
		return nil
	}

	variables := s.PluginConfig.Variables
	module := s.RootConfig.Module

	// populate data for templates from the variables
	templateData := make(map[string]string)
	for _, variable := range variables {
		switch variable.Name {
		case "Module":
			templateData[variable.Name] = module
		default:
			templateData[variable.Name] = ""
		}
	}

	jobs := []Job{}
	jobError := false
	for _, action := range s.PluginConfig.Actions {
		content, _ := s.EmbedFS.ReadFile(
			path.Join(
				"schematics",
				"plugins",
				contextName,
				"templates",
				action.Template,
			),
		)

		outputPath, _ := s.File.ParseFilePath(action.Output, templateData)
		importPath, _ := s.File.ParseFilePath(action.Import, templateData)

		createFileAction := &CreateFileAction{
			OutputPath: outputPath,
			Content:    content,
		}
		importAction := &ImportAction{
			OutputPath: outputPath,
			ImportPath: importPath,
			Alias:      action.Alias,
		}
		dependencyAction := &DependencyAction{
			OutputPath: outputPath,
			Dependency: action.Dependency,
		}
		routeAction := &RouteAction{
			OutputPath: outputPath,
			Route:      action.Route,
		}
		jobs = append(jobs, Job{
			Type:             action.Type,
			CreateFileAction: createFileAction,
			ImportAction:     importAction,
			DependencyAction: dependencyAction,
			RouteAction:      routeAction,
		})
	}

	if jobError {
		log.Printf("Aborting due to errors...\n")
		return nil
	}

Loop:
	for _, job := range jobs {
		switch job.Type {
		case "create_file":
			created := s.actionCreateFile(job.CreateFileAction, templateData)
			if !created {
				jobError = true
				break Loop
			}
		case "add_import":
			created := s.actionAddImport(job.ImportAction)
			if !created {
				jobError = true
				break Loop
			}
		case "add_dependency":
			created := s.actionAddDependency(job.DependencyAction)
			if !created {
				jobError = true
				break Loop
			}
		case "add_route":
			created := s.actionAddRoute(job.RouteAction)
			if !created {
				jobError = true
				break Loop
			}
		default:
			log.Printf("Unknown action type: %s\n", job.Type)
			jobError = true
		}
	}

	if jobError {
		log.Printf("Completed with errors...\n")
		return nil
	}

	log.Printf("All done!\n")

	return nil
}

func (s *addService) actionCreateFile(createFileAction *CreateFileAction, templateData map[string]string) bool {
	parsedContent, err := s.File.ParseTemplate(createFileAction.Content, templateData)
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		return false
	}
	err = s.File.CreateFile(createFileAction.OutputPath, parsedContent)
	if err != nil {
		log.Printf("Failed to write file: %v\n", err)
		return false
	}
	log.Printf("Created: %s\n", createFileAction.OutputPath)
	return true
}

func (s *addService) actionAddImport(importAction *ImportAction) bool {
	parser := parser.NewASTParser(importAction.OutputPath)
	parser.AddImport(importAction.ImportPath, importAction.Alias)
	parser.WriteFile()
	return true
}

func (s *addService) actionAddDependency(dependencyAction *DependencyAction) bool {
	parser := parser.NewASTParser(dependencyAction.OutputPath)
	parser.AddDependencies([]string{dependencyAction.Dependency})
	parser.WriteFile()
	return true
}

func (s *addService) actionAddRoute(routeAction *RouteAction) bool {
	parser := parser.NewASTParser(routeAction.OutputPath)
	parser.AddRoute(routeAction.Route)
	parser.WriteFile()
	return true
}
