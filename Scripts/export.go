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
	defaultSize  = 64
	readmeHeader = `# Discord Emoji

Discord emoji made by me (mostly for the Rust Programming Language Community Server).

## Preview

`
)

var (
	sourceFlag = flag.String("s", "./Sources", "Source directory")
	exportFlag = flag.String("e", "./Export", "Export directory")
	sizeFlag   = flag.Int("z", defaultSize, "Image size")
)

func main() {
	flag.Parse()

	size := strconv.Itoa(*sizeFlag)
	var suffix string
	if *sizeFlag != defaultSize {
		suffix = "_" + size
	}

	root, _ := os.Getwd()

	exportDir, _ := filepath.Abs(*exportFlag)
	if err := os.RemoveAll(exportDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	readmeFile, _ := os.Create("README.md")
	if _, err := io.WriteString(readmeFile, readmeHeader); err != nil {
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

		cmd := exec.Command("inkscape", "--export-filename="+exportPath, "--export-width="+size, path)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			return nil
		}

		relPath, _ := filepath.Rel(root, exportPath)
		if _, err := io.WriteString(readmeFile, "!["+emoji+"](./"+relPath+" \""+emoji+"\")\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			return nil
		}

		return nil
	})
}
