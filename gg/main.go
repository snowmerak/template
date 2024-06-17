package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
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

	log.Println("Done")
}
