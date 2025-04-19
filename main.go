package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var urlFlag string
	var quickPull bool
	flag.StringVar(&urlFlag, "url", "", "URL for PR generation (optional)")
	flag.BoolVar(&quickPull, "quick-pull", false, "Activate quick-pull mode")
	flag.Parse()

	// Check if .propener.toml exists
	if _, err := os.Stat(".propener.toml"); err == nil {
		// Load from file
		content, err := ioutil.ReadFile(".propener.toml")
		if err != nil {
			log.Fatalf("Error reading .propener.toml: %v", err)
		}
		urlFlag = strings.TrimSpace(string(content))
	}

	mainBranch, err := getMainBranch()
	if err != nil {
		log.Fatalf("Error getting main branch: %v", err)
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		log.Fatalf("Error getting current branch: %v", err)
	}

	if currentBranch == "main" || currentBranch == "develop" || currentBranch == "master" {
		log.Fatalf("Current branch (%s) cannot be used to open a PR", currentBranch)
	}

	title, body := "", ""
	if quickPull {
		title = "PR from " + currentBranch
		body = "Quick PR from branch " + currentBranch
	} else {
		title = generateTitle()
		body = generateBody()
	}

	prURL := fmt.Sprintf("https://github.com/your-org/your-repo/compare/%s...%s?expand=1&title=%s&body=%s", mainBranch, currentBranch, urlEncode(title), urlEncode(body))
	fmt.Println("PR URL:", prURL)
}

func getMainBranch() (string, error) {
	candidates := []string{"main", "develop", "master"}
	for _, branch := range candidates {
		cmd := exec.Command("git", "rev-parse", "--verify", branch)
		if err := cmd.Run(); err == nil {
			return branch, nil
		}
	}
	return "", fmt.Errorf("no main branch found")
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func generateTitle() string {
	// Call ChatGPT API to generate title; stub implementation
	return "Generated PR Title"
}

func generateBody() string {
	commitMessages, err := getCommitMessages()
	if err != nil {
		return "Error retrieving commit messages"
	}
	// Call ChatGPT API to generate summary; stub implementation
	summary := "Summary of changes"
	return summary + "\n\nCommits:\n" + commitMessages
}

func getCommitMessages() (string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:(%h): %s", "HEAD", "^origin")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func urlEncode(s string) string {
	// Simple URL encoding, more robust implementation might use url.QueryEscape
	return strings.ReplaceAll(s, " ", "+")
}
