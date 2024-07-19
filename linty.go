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

	ignore "github.com/sabhiram/go-gitignore"
)

type LintyConfig struct {
	Workers   int      `json:"workers"`
	Gitignore bool     `json:"gitignore"`
	Ignore    []string `json:"ignore"`
	Verbose   bool     `json:"verbose"`
	Lint      []struct {
		Type  string `json:"type"`
		Regex string `json:"regex"`
	} `json:"lint"`
}

type LintResult struct {
	File   string
	Result bool
	Error  string
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
	results := runLintChecks(files, config, jsPath)

	for _, result := range results {
		if !result.Result {
			fmt.Printf("Lint failed for file: %s\nError: %s\n", result.File, result.Error)
			os.Exit(1)
		}
	}

	fmt.Println("All lint checks passed!")
}

func readConfig(path string) LintyConfig {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logError("Failed to read config: %v", err)
		os.Exit(1)
	}

	var config LintyConfig
	if err := json.Unmarshal(data, &config); err != nil {
		logError("Failed to parse config: %v", err)
		os.Exit(1)
	}

	logVerbose(config, "Config loaded: %+v", config)
	return config
}

func loadGitignore() *ignore.GitIgnore {
	data, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		logError("Failed to read .gitignore: %v", err)
		return nil
	}

	logVerbose(LintyConfig{Verbose: true}, ".gitignore loaded")
	return ignore.CompileIgnoreLines(strings.Split(string(data), "\n")...)
}

func getFiles(root string, config LintyConfig, gitIgnore *ignore.GitIgnore) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldSkipFile(path, info, config, gitIgnore) {
			logVerbose(config, "Skipping file: %s", path)
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
			logVerbose(config, "Found file: %s", path)
		}
		return nil
	})
	if err != nil {
		logError("Failed to get files: %v", err)
		os.Exit(1)
	}
	return files
}

func shouldSkipFile(path string, info os.FileInfo, config LintyConfig, gitIgnore *ignore.GitIgnore) bool {
	if gitIgnore != nil && gitIgnore.MatchesPath(path) {
		if info.IsDir() {
			return true
		}
		return false
	}

	for _, pattern := range config.Ignore {
		if strings.HasSuffix(pattern, "/") {
			if info.IsDir() && strings.HasPrefix(path, pattern) {
				return true
			}
		} else {
			match, _ := filepath.Match(pattern, filepath.Base(path))
			if match {
				return true
			}
		}
	}

	if strings.Contains(path, ".github") {
		if info.IsDir() {
			return true
		}
		return false
	}

	return false
}

func runLintChecks(files []string, config LintyConfig, jsPath string) []LintResult {
	var results []LintResult

	for _, file := range files {
		result := runLintCheck(file, config, jsPath)
		results = append(results, result)
	}

	return results
}

func runLintCheck(file string, config LintyConfig, jsPath string) LintResult {
	logVerbose(config, "Running lint check on file: %s", file)
	for _, lintConfig := range config.Lint {
		logVerbose(config, "Checking with regex: %s", lintConfig.Regex)
		match, _ := regexp.MatchString(lintConfig.Regex, file)
		if match {
			logVerbose(config, "Matched regex: %s", lintConfig.Regex)
			cmd := exec.Command("node", filepath.Join(jsPath, fmt.Sprintf("%s.js", lintConfig.Type)), file)
			logVerbose(config, "Executing command: %s", cmd.String())
			output, err := cmd.CombinedOutput()
			if err != nil {
				return LintResult{File: file, Result: false, Error: fmt.Sprintf("Failed to run lint check: %v", err)}
			}

			var lintResults []LintResult
			if err := json.Unmarshal(output, &lintResults); err != nil {
				return LintResult{File: file, Result: false, Error: fmt.Sprintf("Failed to parse lint results: %v", err)}
			}

			for _, result := range lintResults {
				if !result.Result {
					return result
				}
			}
		}
	}

	return LintResult{File: file, Result: true}
}

func logError(format string, args ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}

func logVerbose(config LintyConfig, format string, args ...interface{}) {
	if config.Verbose {
		fmt.Printf("VERBOSE: "+format+"\n", args...)
	}
}
