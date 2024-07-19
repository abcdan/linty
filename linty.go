package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	ignore "github.com/sabhiram/go-gitignore"
)

type LintyConfig struct {
	Workers   int      `json:"workers"`
	Gitignore bool     `json:"gitignore"`
	Ignore    []string `json:"ignore"`
	Lint      []struct {
		Type  string `json:"type"`
		Regex string `json:"regex"`
	} `json:"lint"`
}

type LintResult struct {
	File   string
	Result bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: linty <path-to-js-files>")
		os.Exit(1)
	}
	jsPath := os.Args[1]

	config := readConfig(filepath.Join(jsPath, "linty.json"))

	var gitIgnore *ignore.GitIgnore
	if config.Gitignore {
		gitIgnore = loadGitignore()
	}

	files := getFiles(".", config, gitIgnore)
	results := runLintChecks(files, config.Workers, config, jsPath)

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

func loadGitignore() *ignore.GitIgnore {
	data, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		fmt.Printf("Failed to read .gitignore: %v\n", err)
		return nil
	}

	return ignore.CompileIgnoreLines(strings.Split(string(data), "\n")...)
}

func getFiles(root string, config LintyConfig, gitIgnore *ignore.GitIgnore) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and files based on .gitignore
		if gitIgnore != nil && gitIgnore.MatchesPath(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directories and files based on ignore array
		for _, pattern := range config.Ignore {
			if strings.HasSuffix(pattern, "/") {
				if info.IsDir() && strings.HasPrefix(path, pattern) {
					return filepath.SkipDir
				}
			} else {
				match, _ := filepath.Match(pattern, filepath.Base(path))
				if match {
					return nil
				}
			}
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to get files: %v\n", err)
		os.Exit(1)
	}
	return files
}

func runLintChecks(files []string, workers int, config LintyConfig, jsPath string) []LintResult {
	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))
	resultChan := make(chan LintResult, len(files))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				result := runLintCheck(file, config, jsPath)
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

func runLintCheck(file string, config LintyConfig, jsPath string) LintResult {
	for _, lintConfig := range config.Lint {
		match, _ := regexp.MatchString(lintConfig.Regex, file)
		if match {
			cmd := exec.Command("node", filepath.Join(jsPath, fmt.Sprintf("%s.js", lintConfig.Type)), file)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Failed to run lint check on %s: %v\n", file, err)
				return LintResult{File: file, Result: false}
			}

			result := strings.TrimSpace(string(output)) == "true"
			return LintResult{File: file, Result: result}
		}
	}

	// If no specific linting script matches, return true by default
	return LintResult{File: file, Result: true}
}
