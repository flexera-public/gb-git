package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

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
	w := Walker{ProjectPath: *projectPath, Emit: emitGitRepos, Process: runGit}
	(&w).Walk()
}

// Walker encapsulates the state needed to recursively traverse a project path and apply a
// filepath.WalkFunc asynchronously.
// The PathGen function should emit path to the walker in channel.
// The PathProcessor function consumes these path and produce output and error output
// asynchronously to the channels it returns.
type Walker struct {
	ProjectPath string
	Emit        PathEmitFunc
	Process     PathProcessFunc
	emitOut     chan string
}

type PathEmitFunc func(path string, info os.FileInfo, out chan string) error
type PathProcessFunc func(path string, out chan string, errOut chan string)

// NOTE: this is not concurrent safe atm
func (w *Walker) Walk() error {
	emitOut := make(chan string)
	out := make(chan string)
	errOut := make(chan string)
	var wg sync.WaitGroup
	w.emitOut = emitOut
	go func() {
		filepath.Walk(w.ProjectPath, w.fwalk)
		close(emitOut)
	}()
	for o := range emitOut {
		wg.Add(1)
		go func(p string) {
			w.Process(p, out, errOut)
			wg.Done()
		}(o)
	}
	go func() {
		wg.Wait()
		close(out)
		close(errOut)
	}()
	select {
	case o := <-out:
		fmt.Printf("OUT: %s\n", o)
	case e := <-errOut:
		fmt.Printf("ERR: %s\n", e)

	}
	return nil
}

func (w *Walker) fwalk(path string, info os.FileInfo, err error) error {
	//fmt.Printf("WALKING %s\n", path)
	return w.Emit(path, info, w.emitOut)
}

// emitGitRepos
func emitGitRepos(path string, info os.FileInfo, emitOut chan string) error {
	if info.IsDir() && info.Name() == ".git" {
		emitOut <- path[:len(path)-4]
		return filepath.SkipDir
	}
	return nil
}

func runGit(path string, out chan string, errOut chan string) {
	cmd := exec.Command("git", "status")
	cmd.Stdout = &channelWriter{out}
	cmd.Stderr = &channelWriter{errOut}
	cmd.Dir = path
	cmd.Start()
}

type channelWriter struct {
	Channel chan string
}

func (c *channelWriter) Write(b []byte) (n int, err error) {
	c.Channel <- string(b)
	return len(b), nil
}
