package libs

import (
	"embed"
	"path/filepath"
	"strings"
)

type Service struct {
	FS           embed.FS
	ProjectName  string
	TemplateRoot string
}

func NewService(projectName string, templateRoot string, fs embed.FS) *Service {
	return &Service{
		FS:           fs,
		ProjectName:  projectName,
		TemplateRoot: templateRoot,
	}
}

func (s *Service) GetPaymentContextTasks(contextName string) (TaskList, error) {
	files, err := s.ListFiles(s.TemplateRoot)
	if err != nil {
		return nil, err
	}

	folderTasks := TaskList{}
	fileTasks := TaskList{}
	folderMap := map[string]bool{} // to check unique folders

	baseFolder := filepath.Join("internal", s.ConvertToLowerCamelCase(contextName)+"_context")

	for _, file := range files {
		outpath, err := s.GetOutpath(file)
		if err != nil {
			return nil, err
		}
		outpath = filepath.Join(baseFolder, outpath)

		outpathDir := filepath.Dir(outpath)
		if _, ok := folderMap[outpathDir]; !ok {
			folderMap[outpathDir] = true
			folderTasks = append(folderTasks, *NewTask(
				outpathDir,
				outpathDir,
				map[string]interface{}{},
				"folder",
			))
		}

		fileTasks = append(fileTasks, *NewTask(
			file,
			outpath,
			map[string]interface{}{
				"ProjectName": s.ProjectName,
			},
			"file",
		))
	}

	return append(folderTasks, fileTasks...), nil
}

func (s *Service) GetContextTasks(contextName string) (TaskList, error) {
	files, err := s.ListFiles(s.TemplateRoot)
	if err != nil {
		return nil, err
	}

	folderTasks := TaskList{}
	fileTasks := TaskList{}
	folderMap := map[string]bool{} // to check unique folders

	baseFolder := filepath.Join("internal", s.ConvertToLowerCamelCase(contextName)+"_context")

	for _, file := range files {
		outpath, err := s.GetOutpath(file)
		if err != nil {
			return nil, err
		}
		outpath = filepath.Join(baseFolder, outpath)

		outpathDir := filepath.Dir(outpath)
		if _, ok := folderMap[outpathDir]; !ok {
			folderMap[outpathDir] = true
			folderTasks = append(folderTasks, *NewTask(
				outpathDir,
				outpathDir,
				map[string]interface{}{},
				"folder",
			))
		}

		fileTasks = append(fileTasks, *NewTask(
			file,
			outpath,
			map[string]interface{}{
				"ProjectName":      s.ProjectName,
				"ContextName":      s.ConvertToUpperCamelCase(contextName),
				"ContextNameCamel": s.ConvertToLowerCamelCase(contextName),
			},
			"file",
		))
	}

	return append(folderTasks, fileTasks...), nil
}

func (s *Service) GetInitTasks() (TaskList, error) {
	files, err := s.ListFiles(s.TemplateRoot)
	if err != nil {
		return nil, err
	}

	folderTasks := TaskList{}
	fileTasks := TaskList{}
	folderMap := map[string]bool{} // to check unique folders

	for _, file := range files {
		outpath, err := s.GetOutpath(file)
		if err != nil {
			return nil, err
		}
		outpath = filepath.Join(s.ProjectName, outpath)

		outpathDir := filepath.Dir(outpath)
		if _, ok := folderMap[outpathDir]; !ok {
			folderMap[outpathDir] = true
			folderTasks = append(folderTasks, *NewTask(
				outpathDir,
				outpathDir,
				map[string]interface{}{},
				"folder",
			))
		}

		fileTasks = append(fileTasks, *NewTask(
			file,
			outpath,
			map[string]interface{}{
				"ProjectName": s.ProjectName,
			},
			"file",
		))
	}

	return append(folderTasks, fileTasks...), nil
}

func (s *Service) ListFiles(path string) ([]string, error) {
	files, err := s.FS.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	results := []string{}
	for _, file := range files {
		if file.IsDir() {
			files, err := s.ListFiles(filepath.Join(path, file.Name()))
			if err != nil {
				return []string{}, err
			}
			results = append(results, files...)
		} else {
			// full path
			results = append(results, filepath.Join(path, file.Name()))
		}
	}

	return results, nil
}

func (s *Service) GetOutpath(path string) (string, error) {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, ext)
	fullpath := filepath.Join(dir, filename)
	fullpath = strings.ReplaceAll(fullpath, s.TemplateRoot, "")

	return fullpath, nil
}

func (s *Service) ConvertToUpperCamelCase(text string) string {
	words := strings.Split(text, "_")
	for i, word := range words {
		if i == 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		} else {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}

func (s *Service) ConvertToLowerCamelCase(text string) string {
	words := strings.Split(text, "_")
	for i, word := range words {
		if i == 0 {
			words[i] = strings.ToLower(word[:1]) + word[1:]
		} else {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}
