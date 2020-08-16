package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/webview/webview"
)

var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

func view(path string) {
	debug := false
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("PostgreSQL Explain Flame")
	w.SetSize(1024, 768, webview.HintNone)
	p, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	url := "file://" + p
	// fmt.Println(url)
	w.Navigate(url)
	w.Run()
}

// Open calls the OS default program for uri
func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmd := exec.Command(run, uri)
	return cmd.Start()
}
