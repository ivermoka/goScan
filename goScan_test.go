package main

// go test | xsv table
// (while in goScan folder)
// not finished, and not relevant, since go versions before 1.11 didn't have go version in go.mod file

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGoScan(t *testing.T){
    dirs := []string{}    
    goVersions := []string{
		"1.0",
		"1.1",
		"1.2",
		"1.3",
		"1.4",
		"1.5",
		"1.6",
		"1.7",
		"1.8",
		"1.9",
		"1.10",
		"1.11",
		"1.12",
		"1.13",
		"1.14",
		"1.15",
		"1.16",
		"1.17",
		"1.18",
		"1.19",
		"1.20",
		"1.21.1",
		"1.21.2",
	}

    for i := 0; i <= len(goVersions) -1; i++ {
        versionStr := strings.Replace(goVersions[i], ".", "_", -1)
        tempDir, err := os.MkdirTemp("", fmt.Sprintf("go%s_tempDir", versionStr))        
        if err != nil{
            log.Fatal(err)
        }
        defer os.RemoveAll(tempDir)
        // go mod content and path creation
        goModContent := "module test1\n\ngo " + goVersions[i] + "\n"
		goModPath := tempDir + "/go.mod"
        // write to temp file
		if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
			t.Fatal(err)
		}  
        dirs = append(dirs, tempDir)
    }

    // list all temporary directories

	for i, dir := range dirs {
		goVersion := getGoVersion(dir)
		// dockerFrom := getDockerFrom(dir.Path)
		fmt.Printf("%s: %s,%s\n", fmt.Sprint(i), dir, goVersion)
	}
}
