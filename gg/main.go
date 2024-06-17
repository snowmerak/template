package main

import (
	"github.com/AlecAivazis/survey/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get current working directory:", err)
		return
	}
	path := ""

	workspaceData := make([]byte, 0)
	for {
		workspaceData, err = os.ReadFile(filepath.Join(pwd, "go.work"))
		if err == nil {
			break
		}

		path = filepath.Join(filepath.Base(pwd), path)
		pwd = filepath.Dir(pwd)

		if pwd == "/" || pwd == "" || strings.HasSuffix(pwd, ":\\") {
			log.Println("Failed to find workspace file")
			return
		}
	}

	templateOptions := []string{
		"hello",
		"http-client",
		"http-server",
		"mono",
		"nats",
		"postgres",
		"redis",
		"executable",
	}

	selectedTemplate := ""
	if err := survey.AskOne(&survey.Select{
		Message: "Select template",
		Options: templateOptions,
	}, &selectedTemplate); err != nil {
		return
	}

	moduleName := ""
	if err := survey.AskOne(&survey.Input{
		Message: "Enter module name",
	}, &moduleName); err != nil {
		return
	}

	cmd := exec.Command("gonew", "github.com/snowmerak/template/"+strings.ReplaceAll(selectedTemplate, "-", ""), moduleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

	modulePath := filepath.Join(path, filepath.Base(moduleName))

	workspaceData = append(workspaceData, []byte("\nuse "+modulePath)...)
	if err := os.WriteFile(filepath.Join(pwd, "go.work"), workspaceData, 0644); err != nil {
		log.Println("Failed to write workspace data:", err)
		log.Println("Please write this line to your go.work file")
		log.Println("- use " + modulePath)
		return
	}

	log.Println("Done")
}
