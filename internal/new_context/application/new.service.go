package application

import (
	"embed"
	"errors"
	"log"
	"path"

	SEmbed "github.com/IrwantoCia/gomakase/embed"
	"github.com/IrwantoCia/gomakase/internal/shared/config"
	"github.com/IrwantoCia/gomakase/internal/shared/file"
)

type NewService interface {
	Generate(name string, schematic config.ProjectSchematic) error
}

func NewNewService(file file.File) NewService {
	return &newService{
		EmbedFS: SEmbed.SchematicsFS,
		File:    file,
	}
}

type newService struct {
	EmbedFS embed.FS
	File    file.File
}

type Job struct {
	OutputPath string
	Content    []byte
}

func (s newService) Generate(name string, schematic config.ProjectSchematic) error {
	log.Printf("Generating a new project: %s\n", name)

	if s.File.IsPathExists(name) {
		log.Printf("Project already exists, skipping...\n")
		return nil
	}

	variables := schematic.Variables
	// populate data for templates from the variables
	data := make(map[string]string)
	for _, variable := range variables {
		switch variable.Name {
		case "Module":
			data[variable.Name] = name
		default:
			data[variable.Name] = ""
		}
	}

	jobs := []Job{}
	jobError := false
	for _, action := range schematic.Actions {
		if s.File.IsPathExists(action.Output) {
			log.Printf("Path already exists, skipping...\n")
			jobError = true
			continue
		}

		content, err := s.EmbedFS.ReadFile(
			path.Join(
				"schematics",
				"project",
				"templates",
				action.Template,
			))
		if err != nil {
			log.Printf("Error reading template file...\n%v", err)
			jobError = true
			continue
		}

		// parse the output path
		outputPath, err := s.File.ParseFilePath(action.Output, data)
		if err != nil {
			log.Printf("Error parsing output path...\n%v", err)
			jobError = true
			continue
		}

		parsedContent, err := s.File.ParseTemplate(content, data)
		if err != nil {
			log.Printf("Error parsing template...\n%v", err)
			jobError = true
			continue
		}

		jobs = append(jobs, Job{
			OutputPath: outputPath,
			Content:    parsedContent,
		})

	}

	if jobError {
		log.Printf("Job error, skipping...\n")
		return errors.New("job error")
	}

	for _, job := range jobs {
		log.Printf("Creating file: %s\n", job.OutputPath)
		err := s.File.CreateFile(job.OutputPath, job.Content)
		if err != nil {
			log.Printf("Error creating file %s: %v\n\n", job.OutputPath, err)
			continue
		}
	}

	return nil
}
