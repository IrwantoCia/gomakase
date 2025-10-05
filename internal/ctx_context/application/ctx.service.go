package application

import (
	"embed"
	"errors"
	"log"
	"path"
	"strings"

	SEmbed "github.com/IrwantoCia/gomakase/embed"
	"github.com/IrwantoCia/gomakase/internal/shared/config"
	"github.com/IrwantoCia/gomakase/internal/shared/file"
)

type CtxService interface {
	Generate(contextName string) error
}

type ctxService struct {
	RootConfig    config.RootSchematic
	ContextConfig config.ContextSchematic
	EmbedFS       embed.FS
	File          file.File
}

func NewCtxService(
	file file.File,
	rootConfig config.RootSchematic,
	contextConfig config.ContextSchematic,
) CtxService {
	return &ctxService{
		EmbedFS:       SEmbed.SchematicsFS,
		File:          file,
		RootConfig:    rootConfig,
		ContextConfig: contextConfig,
	}
}

type Job struct {
	OutputPath string
	Content    []byte
}

func (s *ctxService) Generate(
	contextName string,
) error {
	log.Printf("Generating a new context: %s\n", contextName)

	if s.File.IsPathExists(path.Join("internal", strings.ToLower(contextName))) {
		log.Printf("Context already exists, skipping...\n")
		return nil
	}

	variables := s.ContextConfig.Variables
	// populate data for templates from the variables
	data := make(map[string]string)
	for _, variable := range variables {
		switch variable.Name {
		case "Module":
			data[variable.Name] = s.RootConfig.Module
		case "ContextName":
			data[variable.Name] = contextName
		default:
			data[variable.Name] = ""
		}
	}

	jobs := []Job{}
	jobError := false
	for _, action := range s.ContextConfig.Actions {
		content, _ := s.EmbedFS.ReadFile(
			path.Join(
				"schematics",
				"context",
				"templates",
				action.Template,
			),
		)

		outputPath, _ := s.File.ParseFilePath(action.Output, data)
		jobs = append(jobs, Job{
			OutputPath: outputPath,
			Content:    content,
		})
	}

	if jobError {
		log.Printf("Job error, skipping...\n")
		return errors.New("job error")
	}

	for _, job := range jobs {
		log.Printf("Creating file: %s\n", job.OutputPath)
		parsedContent, err := s.File.ParseTemplate(job.Content, data)
		if err != nil {
			log.Printf("Error parsing template: %v\n", err)
			continue
		}
		err = s.File.CreateFile(job.OutputPath, parsedContent)
		if err != nil {
			log.Printf("Error creating file %s: %v\n\n", job.OutputPath, err)
			continue
		}
	}

	return nil
}
