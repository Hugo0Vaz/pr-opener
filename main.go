package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	var urlFlag string
	var quickPull bool
	flag.StringVar(&urlFlag, "url", "", "URL for PR generation (optional)")
	flag.BoolVar(&quickPull, "quick-pull", false, "Activate quick-pull mode")
	flag.Parse()

	var baseUrl string
	var apiKey string

	if urlFlag == "" {
		fmt.Println("Loading configs from .propener.toml ")
		var err error
		baseUrl, apiKey, err = loadTomlConfigs()
		if err != nil {
			fmt.Println("Cannot load the config from the .propner, exiting...")
		}
	} else {
		fmt.Println("URL argument passed, ignoring any .propener.toml")
		baseUrl = urlFlag

		envKey := os.Getenv("OPENAI_API_KEY")

		if envKey == "" {
			fmt.Println("No OPENAI_API_KEY in the enviroment, exiting...")
		}
		apiKey = envKey
	}

	fmt.Println(baseUrl, apiKey)
}

func loadTomlConfigs() (string, string, error) {
	var repoCfg struct {
		Repo struct {
			Base       string `toml:"base"`
			RepoOwner  string `toml:"repo_owner"`
			Repo       string `toml:"repo"`
			BaseBranch string `toml:"base_branch"`
		} `toml:"repo"`

		Ai struct {
			OpenAiApiKey string `toml:"openai_api_key"`
		} `toml:"ai"`
	}

	if _, err := os.Stat(".propener.toml"); err != nil {
		fmt.Println("Could not open .propener.toml")
	}

	if _, err := toml.DecodeFile(".propener.toml", &repoCfg); err != nil {
		fmt.Println("Coul not decode .propener.tml")
	}

	cb, err := getCurrentBranch()
	if err != nil {
		fmt.Println(err)
	}

	base := getBasePrUrl(repoCfg.Repo.Base, repoCfg.Repo.RepoOwner, repoCfg.Repo.Repo, repoCfg.Repo.BaseBranch, cb)
	key := repoCfg.Ai.OpenAiApiKey

	return base, key, nil
}

func getBasePrUrl(b string, ro string, r string, bb string, cb string) string {
	return b + ro + "/" + r + "/" + "compare" + "/" + bb + "..." + cb
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
