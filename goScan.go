package main

// run command in directory with files and goScan.go:
// go build -o go1/goScan go1/goScan.go && go1/goScan | xsv table
import (
	"fmt"
	"log"
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
		log.Fatal(err)
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