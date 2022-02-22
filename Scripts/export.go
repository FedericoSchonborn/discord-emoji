package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultWidth = 64
)

var (
	sourceFlag = flag.String("s", "./Sources", "Source directory")
	exportFlag = flag.String("e", "./Export", "Export directory")
	widthFlag  = flag.Int("w", defaultWidth, "Image width")
)

func main() {
	flag.Parse()

	width := strconv.Itoa(*widthFlag)
	var suffix string
	if *widthFlag != defaultWidth {
		suffix = "_" + width
	}

	exportDir, _ := filepath.Abs(*exportFlag)

	if err := os.RemoveAll(exportDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	readmeFile, _ := os.Create(filepath.Join(exportDir, "README.md"))
	if _, err := io.WriteString(readmeFile, "<!-- markdownlint-disable MD041 MD045 -->\n"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	svgDir, _ := filepath.Abs(*sourceFlag)
	filepath.WalkDir(svgDir, func(path string, entry fs.DirEntry, err error) error {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".svg" {
			return nil
		}

		emoji := strings.TrimSuffix(name, ".svg")
		exportName := emoji + suffix + ".png"
		exportPath := filepath.Join(exportDir, exportName)

		cmd := exec.Command("inkscape", "--export-filename="+exportPath, "--export-width="+width, path)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			return nil
		}

		relPath, _ := filepath.Rel(exportDir, exportPath)
		if _, err := io.WriteString(readmeFile, "![](./"+relPath+" \""+emoji+"\")\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			return nil
		}

		return nil
	})
}
