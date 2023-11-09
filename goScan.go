package main

// run command in directory with files and goScan.go:
// go build -o goScan/goScan goScan/goScan.go && goScan/goScan (arguments) | xsv table
// go repos with go version >1.11 will be shown as >1.11 (version is not in go.mod)

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var bvar bool
	flag.BoolVar(&bvar, "l", false, "latest")
	flag.Parse()
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
	latestVersion := getLatestGoVersion()
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
				if bvar {
					os.Args = append(os.Args, latestVersion)
				}
				isArg := false
				for _, args := range os.Args[1:] {
					if goVersion == string(args) {
						isArg = true
						break
					}
				}
				if isArg {
					continue
				}
			}
			if goVersion == "" {
				goVersion = ">1.11"
			}
			dockerFrom := getDockerFrom(v.Name())
			fmt.Printf("%s,%s,%s\n", v.Name(), goVersion, dockerFrom)
		}
	}
}

func isGoProject(pt string) bool {
	_, err := os.Stat(pt + "/go.mod")
	return err == nil
}

func getGoVersion(pt string) string {
	cmd := exec.Command("grep", "^go", "go.mod")
	cmd.Dir = pt
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(string(output)) != 0 {
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
	for _, v := range lines {
		if strings.HasPrefix(v, "go ") {
			return true
		}
	}
	return false
}
func getDockerFrom(pt string) string {
	cmd := exec.Command("grep", "-m", "1", "^FROM", "Dockerfile")
	cmd.Dir = pt
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	trimmed := strings.TrimSpace(string(output))
	trimmed = strings.TrimPrefix(trimmed, "FROM ")
	index := strings.Index(strings.ToLower(trimmed), " as ")
	if index != -1 {
		trimmed = trimmed[:index]
	}
	trimmedVersion := strings.TrimLeft(trimmed, "golan:")
	if trimmedVersion == "1.21.4" {
		return "correct version"
	}
	return trimmed
}

func getLatestGoVersion() string {
	const filteredContentStart = 20
	const filteredContentLength = 128
	const goRemover = 6
	res, err := http.Get("https://go.dev/dl/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	content := string(body)
	indexStable := strings.Index(strings.ToLower(content), "id=\"stable\"")
	filteredContent := content[indexStable+filteredContentStart:indexStable+filteredContentLength]
	indexVersion := strings.Index(filteredContent, "id=\"")
	indexVersionEnd := strings.Index(filteredContent, "\">")
	version := filteredContent[indexVersion+goRemover:indexVersionEnd]
	return version
}