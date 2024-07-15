package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
)

type LintyConfig struct {
	Workers int `json:"workers"`
	Php     []struct {
		Type  string `json:"type"`
		Regex string `json:"regex"`
	} `json:"php"`
}

type LintResult struct {
	File   string
	Result bool
}

func main() {
	config := readConfig("linty/linty.json")

	files := getGoFiles(".")
	results := runLintChecks(files, config.Workers, config)

	for _, result := range results {
		if !result.Result {
			fmt.Printf("Lint failed for file: %s\n", result.File)
			os.Exit(1)
		}
	}

	fmt.Println("All lint checks passed!")
}

func readConfig(path string) LintyConfig {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read config: %v\n", err)
		os.Exit(1)
	}

	var config LintyConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	return config
}

func getGoFiles(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to get Go files: %v\n", err)
		os.Exit(1)
	}
	return files
}

func runLintChecks(files []string, workers int, config LintyConfig) []LintResult {
	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))
	resultChan := make(chan LintResult, len(files))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				result := runLintCheck(file, config)
				resultChan <- result
			}
		}()
	}

	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	wg.Wait()
	close(resultChan)

	var results []LintResult
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func runLintCheck(file string, config LintyConfig) LintResult {
	for _, phpConfig := range config.Php {
		match, _ := regexp.MatchString(phpConfig.Regex, file)
		if match {
			cmd := exec.Command("node", fmt.Sprintf("linty/%s.js", phpConfig.Type), file)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Failed to run lint check on %s: %v\n", file, err)
				return LintResult{File: file, Result: false}
			}

			result := string(output) == "true"
			return LintResult{File: file, Result: result}
		}
	}

	cmd := exec.Command("node", "linty/generic.js", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to run lint check on %s: %v\n", file, err)
		return LintResult{File: file, Result: false}
	}

	result := string(output) == "true"
	return LintResult{File: file, Result: result}
}
