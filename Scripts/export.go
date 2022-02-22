package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	exportDir, _ := filepath.Abs("./Export")
	if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	readmeFile, _ := os.Create(filepath.Join(exportDir, "README.md"))

	svgDir, _ := filepath.Abs("./Sources")
	filepath.WalkDir(svgDir, func(path string, entry fs.DirEntry, err error) error {
		if filepath.Ext(path) != ".svg" {
			return nil
		}

		svgName := filepath.Base(path)
		baseName := strings.TrimSuffix(svgName, ".svg")
		pngName := baseName + ".png"
		exportPath := filepath.Join(exportDir, pngName)
		cmd := exec.Command("inkscape", "--export-filename="+exportPath, "--export-width=32", path)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			return nil
		}

		if _, err := io.WriteString(readmeFile, "!["+baseName+"](./"+pngName+")"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			return nil
		}

		return nil
	})
}
