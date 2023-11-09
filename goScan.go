package main

// run command in directory with files and goScan.go:
// go build -o goScan/goScan goScan/goScan.go && goScan/goScan (arguments) | xsv table

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main(){
	f, err := os.Open(".")
    if err != nil {
        fmt.Println(err)
        return
    }
	files, err := f.Readdir(0)
    if err != nil {
        fmt.Println(err)
        return
    }
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", "No files found.")
	}
	for _, v := range files {
		if !v.IsDir() {
			continue
		}
		if isGoProject(v.Name()) {
			f, err := os.Open(v.Name() + "/go.mod")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", "Could not open file.")
			}
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", "Could not read file.")
			}
			var goVersion string
			if hasGoVersion(string(b)) {
				goVersion = getGoVersion(v.Name())
				os.Args = append(os.Args, "1.21.4")
				isArg := false
				for _, args := range os.Args { 
					if goVersion == string(args) {
						isArg = true
						break
					}
				}
				if isArg {
					continue
				}
			}
			dockerFrom := getDockerFrom(v.Name())
			fmt.Printf("%s,%s,%s\n", v.Name(), goVersion, dockerFrom)
		}
	}
}

func isGoProject(pt string) bool {
	_, err := os.Stat(pt + "/go.mod") 
	if err != nil {
		return false
	}

	return true
}

func getGoVersion(pt string) string{
	cmd := exec.Command("grep", "^go", "go.mod")
	cmd.Dir = pt
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(string(output)) != 0{
			fmt.Fprintf(os.Stderr, "error checking go version in %s: %s\noutput was %s\n", pt, err, string(output))
		} else {
			fmt.Fprintf(os.Stderr, "error checking go version in %s: %s\n", pt, err)
		}
	}
	trimmed := strings.TrimSpace(string(output))
	trimmed = strings.TrimPrefix(trimmed, "go ")
	return trimmed
}
func hasGoVersion(content string) bool {
	lines := strings.Split(content, "\n")
	for _, v := range lines{
		if strings.HasPrefix(v, "go "){
			return true
		}
	}
	return false
}
func getDockerFrom(pt string) string{
	cmd := exec.Command("grep", "-m", "1", "^FROM", "Dockerfile")
	cmd.Dir = pt
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	trimmed := strings.TrimSpace(string(output))
	trimmed = strings.TrimPrefix(trimmed, "FROM ")
	index := strings.Index(trimmed, " as ")
	if index != -1 {
		trimmed = trimmed[:index]
	}
	trimmedVersion := strings.TrimLeft(trimmed, "golan:")
	if trimmedVersion == "1.21.4" {
		return "correct version"
	}
	return trimmed
}

