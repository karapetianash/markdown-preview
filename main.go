package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>Markdown Preview Tool</title>
</head>
<body>
`
	footer = `
</body>
</html>
`
)

func main() {
	filename := flag.String("file", "", "Markdown file to preview.")
	skipPreview := flag.Bool("s", false, "Skip auto-preview.")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run is a coordinating function
func run(filename string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err = saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview()
}

// parseContent function parse Markdown into HTML
func parseContent(input []byte) []byte {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

// saveHTML function saves result into a file
func saveHTML(outName string, htmlData []byte) error {
	return os.WriteFile(outName, htmlData, 0664)
}

// preview function runs executable command to show an HTML file
func preview() error {
	cName := ""
	cParams := make([]string, 0)

	// Define executable depending on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = append(cParams, "/C", "start")
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS is not supported")
	}

	// Appending filename to parameters slice
	cParams = append(cParams, cName)

	// Locating executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Opening the file using default program
	err = exec.Command(cPath, cParams...).Run()

	// Giving the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)

	return err
}
