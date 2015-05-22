package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	projectPath = kingpin.Flag("project", "Path to project root").String()
)

func main() {
	kingpin.Parse()
	if *projectPath == "" {
		cur, err := os.Getwd()
		if err != nil {
			kingpin.Fatalf("Failed to retrieve current directory: %s", err)
		}
		*projectPath = cur
	}
	kingpin.FatalIfError(filepath.Walk(*projectPath, walkProject), "walk")
}

func walkProject(path string, info os.FileInfo, err error) error {
	if info.IsDir() && info.Name() == ".git" {
		repo := path[:len(path)-4]
		cmd := exec.Command("git", "status")
		cmd.CombinedOutput
	}

}
