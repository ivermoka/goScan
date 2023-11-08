package main

// run command in directory with files and goScan.go:
// go build -o goScan/goScan goScan/goScan.go && goScan/goScan | xsv table
import (
	"fmt"
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
			goVersion := getGoVersion(v.Name())
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
		fmt.Fprintf(os.Stderr, "error checking go version in %s: %s\noutput was %s", pt, err, string(output))
	}
	trimmed := strings.TrimSpace(string(output))
	trimmed = strings.TrimPrefix(trimmed, "go ")
	return trimmed
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
	return trimmed
}