package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

type LintyConfig struct {
	Abcdan    string   `json:"abcdan"`
	Gitignore bool     `json:"gitignore"`
	Ignore    []string `json:"ignore"`
	Verbose   bool     `json:"verbose"`
	Secure    bool     `json:"secure"`
	Lint      []struct {
		Type   string `json:"type"`
		Regex  string `json:"regex"`
		Linter string `json:"linter"`
	} `json:"lint"`
}

type LintResult struct {
	File   string `json:"file"`
	Result bool   `json:"result"`
	Error  string `json:"error"`
}

type GitHubCommit struct {
	Commit struct {
		Author struct {
			Name string `json:"name"`
		} `json:"author"`
	} `json:"commit"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: linty <path-to-js-files>")
		os.Exit(1)
	}
	jsPath := os.Args[1]

	config := readConfig(filepath.Join(jsPath, "linty.json"))

	if config.Secure {
		if !checkIntegrity() {
			fmt.Println("Integrity check failed. Aborting.")
			os.Exit(1)
		}

		if !checkLintycheckAuthor() {
			fmt.Println("LINTYCHECK file was not updated by the authorized user. Aborting.")
			os.Exit(1)
		}
	}

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
		for _, lintConfig := range config.Lint {
			if match, _ := regexp.MatchString(lintConfig.Regex, file); match {
				result := runLintCheck(file, lintConfig, jsPath, config)
				results = append(results, result)
			}
		}
	}

	return results
}

func runLintCheck(file string, lintConfig struct {
	Type   string `json:"type"`
	Regex  string `json:"regex"`
	Linter string `json:"linter"`
}, jsPath string, config LintyConfig) LintResult {
	logVerbose(config, "Running lint check on file: %s with linter: %s", file, lintConfig.Linter)
	cmd := exec.Command("node", "linty.js", lintConfig.Linter, file)
	logVerbose(config, "Executing command: %s", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		logVerbose(config, "Command output: %s", string(output))
		return LintResult{File: file, Result: false, Error: fmt.Sprintf("Failed to run lint check: %v", err)}
	}

	var lintResults []LintResult
	if err := json.Unmarshal(output, &lintResults); err != nil {
		logVerbose(config, "Command output: %s", string(output))
		return LintResult{File: file, Result: false, Error: fmt.Sprintf("Failed to parse lint results: %v", err)}
	}

	for _, result := range lintResults {
		logVerbose(config, "Test: %s, File: %s, Result: %t", lintConfig.Type, result.File, result.Result)
		if !result.Result {
			return result
		}
	}

	return LintResult{File: file, Result: true}
}

func checkIntegrity() bool {
	data, err := ioutil.ReadFile("LINTYCHECK")
	if err != nil {
		logError("Failed to read LINTYCHECK file: %v", err)
		return false
	}

	lines := strings.Split(string(data), "\n")
	checksums := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			checksums[parts[0]] = parts[1]
		}
	}

	if !verifyChecksum("linty.go", checksums["linty.go"]) {
		return false
	}

	if !verifyChecksum("linty.js", checksums["linty.js"]) {
		return false
	}

	return true
}

func verifyChecksum(file, expectedChecksum string) bool {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logError("Failed to read file %s: %v", file, err)
		return false
	}

	checksum := md5.Sum(data)
	actualChecksum := hex.EncodeToString(checksum[:])

	if actualChecksum != expectedChecksum {
		logError("Checksum mismatch for file %s. Expected: %s, Actual: %s", file, expectedChecksum, actualChecksum)
		return false
	}

	return true
}

func checkLintycheckAuthor() bool {
	url := "https://api.github.com/repos/abcdan/linty/commits?path=LINTYCHECK&page=1&per_page=1"
	resp, err := http.Get(url)
	if err != nil {
		logError("Failed to fetch LINTYCHECK commit information: %v", err)
		return false
	}
	defer resp.Body.Close()

	var commits []GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		logError("Failed to parse LINTYCHECK commit information: %v", err)
		return false
	}

	if len(commits) == 0 {
		logError("No commits found for LINTYCHECK file")
		return false
	}

	author := commits[0].Commit.Author.Name
	if author != "abcdan" {
		logError("LINTYCHECK file was last updated by %s, expected abcdan", author)
		return false
	}

	return true
}

func logError(format string, args ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}

func logVerbose(config LintyConfig, format string, args ...interface{}) {
	if config.Verbose {
		fmt.Printf("VERBOSE: "+format+"\n", args...)
	}
}
