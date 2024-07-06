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
	if len(os.Args) < 2 {
		log.Println("Usage: gg <command>")
		log.Println("Commands:")
		log.Println("  init: Initialize workspace")
		log.Println("  clone: Clone a repository from template")
		return
	}

	switch os.Args[1] {
	case "init":
		rootWorkspace := ""
		if err := survey.AskOne(&survey.Input{
			Message: "Enter root workspace repository",
		}, &rootWorkspace); err != nil {
			log.Println("Failed to get root workspace repository:", err)
			return
		}

		cmd := exec.Command("go", "work", "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println("Failed to initialize workspace:", err)
			return
		}

		workspaceData, err := os.ReadFile("go.work")
		if err != nil {
			log.Println("Failed to read workspace data:", err)
			return
		}

		workspaceData = append([]byte("// "+rootWorkspace+"\n"), workspaceData...)
		if err := os.WriteFile("go.work", workspaceData, 0644); err != nil {
			log.Println("Failed to write workspace data:", err)
			return
		}
	case "clone":
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

		rootWorkspace := strings.TrimPrefix(strings.SplitN(string(workspaceData), "\n", 2)[0], "// ")

		templateOptions := []string{
			"hello: a simple hello world",
			"http-client: a simple http client",
			"http-server: a simple http server based on gin",
			"mono: a mono-repo application",
			"nats: a simple nats client",
			"postgres: a simple postgres and postgis client based on sqlc",
			"redis: a simple redis client based on rueidis",
			"executable: a simple executable application",
			"cassandra: a simple cassandra client based on gocql and gocqlx",
			"s3: a simple s3 client based on minio",
			"mongo: a simple mongo client based on mongo-driver",
			"ollama: a simple ollama client",
			"empty: an empty module",
			"google-ai-studio: a simple google ai studio client",
		}

		selectedTemplate := ""
		if err := survey.AskOne(&survey.Select{
			Message: "Select template",
			Options: templateOptions,
		}, &selectedTemplate); err != nil {
			return
		}

		selectedTemplate = strings.SplitN(selectedTemplate, ":", 2)[0]

		moduleName := ""
		if err := survey.AskOne(&survey.Input{
			Message: "Enter module name",
		}, &moduleName); err != nil {
			return
		}

		cmd := exec.Command("gonew", "github.com/snowmerak/template/"+strings.ReplaceAll(selectedTemplate, "-", ""), filepath.ToSlash(filepath.Join(rootWorkspace, moduleName)))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return
		}

		modulePath := filepath.ToSlash(filepath.Join(path, filepath.Base(moduleName)))

		workspaceData = append(workspaceData, []byte("\nuse "+modulePath)...)
		if err := os.WriteFile(filepath.Join(pwd, "go.work"), workspaceData, 0644); err != nil {
			log.Println("Failed to write workspace data:", err)
			log.Println("Please write this line to your go.work file")
			log.Println("- use " + modulePath)
			return
		}
	}

	log.Println("Done")
}
