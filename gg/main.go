package main

import (
	"github.com/AlecAivazis/survey/v2"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	workspaceData, err := os.ReadFile("go.work")
	if err != nil {
		log.Println("Failed to read workspace data:", err)
		log.Println("Please run this command in the root of your workspace")
		log.Println("- go work init")
		return
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
		"command",
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

	workspaceData = append(workspaceData, []byte("\nuse "+moduleName)...)
	if err := os.WriteFile("go.work", workspaceData, 0644); err != nil {
		log.Println("Failed to write workspace data:", err)
		log.Println("Please write this line to your go.work file")
		log.Println("- use " + moduleName)
		return
	}

	log.Println("Done")
}
