package libs

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InitProject(projectName string, templatesFS embed.FS) {
	// check if project folder exists
	if _, err := os.Stat(projectName); err == nil {
		fmt.Printf("Project folder %s already exists\n", projectName)
		os.Exit(1)
	}

	// get tasks
	service := NewService(projectName, "templates/init", templatesFS)
	tasks, err := service.GetInitTasks()
	if err != nil {
		fmt.Println("Error getting tasks")
		os.Exit(1)
	}

	// create project folder
	fmt.Println("Creating project folder")
	os.MkdirAll(projectName, 0755)

	tasks.Run(templatesFS)

	// cd to project folder
	fmt.Println("Changing directory to project folder")
	os.Chdir(projectName)

	// run go mod tidy
	fmt.Println("Running go mod tidy")
	exec.Command("go", "mod", "tidy").Run()

	// npm install
	fmt.Println("Running npm install")
	exec.Command("npm", "install").Run()

	// build css
	fmt.Println("Running npm run build:css")
	exec.Command("npm", "run", "build:css").Run()

}

func CreateContext(contextName string, templatesFS embed.FS) {
	// check if on project folder, by checking if go.mod exists
	if _, err := os.Stat("go.mod"); err != nil {
		fmt.Println("Not in project folder")
		os.Exit(1)
	}

	// get project name, go.mod folder name
	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)
	fmt.Println("Project name: ", projectName)

	// get tasks
	service := NewService(projectName, "templates/context", templatesFS)
	tasks, err := service.GetContextTasks(contextName)
	if err != nil {
		fmt.Println("Error getting tasks")
		os.Exit(1)
	}

	tasks.Run(templatesFS)

	// run go mod tidy
	fmt.Println("Running go mod tidy")
	exec.Command("go", "mod", "tidy").Run()
}

func CreatePaymentContext(contextName string, templatesFS embed.FS) {
	// check if on project folder, by checking if go.mod exists
	if _, err := os.Stat("go.mod"); err != nil {
		fmt.Println("Not in project folder")
		os.Exit(1)
	}

	// get project name, go.mod folder name
	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)
	fmt.Println("Project name: ", projectName)

	// get tasks
	service := NewService(projectName, "templates/payment_context", templatesFS)
	tasks, err := service.GetPaymentContextTasks(contextName)
	if err != nil {
		fmt.Println("Error getting tasks")
		os.Exit(1)
	}

	tasks.Run(templatesFS)

	// run go mod tidy
	fmt.Println("Running go mod tidy")
	exec.Command("go", "mod", "tidy").Run()
}
