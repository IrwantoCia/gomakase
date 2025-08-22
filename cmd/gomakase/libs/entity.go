package libs

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Task struct {
	Template string
	Data     map[string]interface{}
	Outpath  string
	Type     string
}

func NewTask(template string, outpath string, data map[string]interface{}, taskType string) *Task {
	if taskType != "file" && taskType != "folder" {
		errMsg := fmt.Sprintf("invalid task type: %s, must be 'file' or 'folder'\n", taskType)
		panic(errMsg)
	}

	return &Task{
		Template: template,
		Data:     data,
		Outpath:  outpath,
		Type:     taskType,
	}
}

func (t *Task) ToFile(templatesFS embed.FS) error {
	if strings.HasSuffix(t.Outpath, ".html") {
		// copy the file instead
		content, err := templatesFS.ReadFile(t.Template)
		if err != nil {
			return err
		}
		err = os.WriteFile(t.Outpath, content, 0666)
		if err != nil {
			return err
		}
		return nil
	}

	tmpl, err := template.ParseFS(templatesFS, t.Template)
	if err != nil {
		return err
	}

	file, err := os.Create(t.Outpath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, t.Data)
}

type TaskList []Task

func (tl TaskList) Run(templatesFS embed.FS) error {
	for _, task := range tl {
		switch task.Type {
		case "file":
			err := task.ToFile(templatesFS)
			if err != nil {
				return err
			}
		case "folder":
			err := os.MkdirAll(task.Outpath, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
