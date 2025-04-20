package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"context"

	openai "github.com/sashabaranov/go-openai"
	"github.com/BurntSushi/toml"
)

func main() {
	var quickPull bool
	var baseCli string
	var repoOwnerCli string
	var repoCli string
	var baseBrachCli string
	flag.BoolVar(&quickPull, "quick-pull", false, "Activate quick-pull mode")
	flag.StringVar(&baseCli,"base", "https://github.com/", "The base url for the tool, default is 'https://github.com/'")
	flag.StringVar(&repoOwnerCli, "repo-owner", "", "The PR reposistory owner")
	flag.StringVar(&repoCli, "repo", "", "The PR repository")
	flag.StringVar(&baseBrachCli, "base-branch", "", "The base branch of the PR")
	flag.Parse()

	currentBranch, err := getCurrentBranch()
	if err != nil {
		fmt.Println("Could not get git current branch")
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("'OPENAI_API_KEY' no set in the enviroment")
		os.Exit(1)
	}

	var base string 
	var repoOwner string
	var repo string
	var baseBranch string

	if ( repoOwnerCli == "" && repoCli == "" &&  baseBrachCli == "" ) {
		fmt.Println("Loading configs from .propener.toml")
		base, repoOwner, repo, baseBranch, err = loadTomlConfigs()
		if err != nil {
			fmt.Println("Cannot load the config from the .propner.toml, exiting...")
			os.Exit(1)
		}
	} else {
		base, repoOwner, repo, baseBranch = baseCli, repoOwnerCli, repoCli, baseBrachCli
		fmt.Println("Loading configs from cli")
	}

	baseUrl := getBasePrUrl(base, repoOwner, repo, baseBranch, currentBranch)
	prePromp := `A partir da lista de commits apresentado abaixo, gere dois itens: 

	- Um título das mudanças no repositório (esse será o título do pull request)
	- O segundo item uma descrição mais longa das mudanças que estão contidas naquele pull request

	A saída deve seguir o seguinte formato:

	<Título>

	<Descrição>
`
	commitBody, _ := generateCommitBody(baseBranch)
	inference, _ := getInference(prePromp, commitBody, apiKey)

	fmt.Println(baseUrl)
	fmt.Println(inference)
}

func loadTomlConfigs() (string, string, string, string, error) {
	var repoCfg struct {
		Repo struct {
			Base       string `toml:"base"`
			RepoOwner  string `toml:"repo_owner"`
			Repo       string `toml:"repo"`
			BaseBranch string `toml:"base_branch"`
		} `toml:"repo"`
	}

	if _, err := os.Stat(".propener.toml"); err != nil {
		fmt.Println("Could not open .propener.toml")
		return "", "", "", "", err
	}

	if _, err := toml.DecodeFile(".propener.toml", &repoCfg); err != nil {
		fmt.Println("Coul not decode .propener.tml")
		return "", "", "", "", err
	}

	return repoCfg.Repo.Base, repoCfg.Repo.RepoOwner, repoCfg.Repo.Repo, repoCfg.Repo.BaseBranch, nil
}

func getBasePrUrl(b string, ro string, r string, bb string, cb string) string {
	return b + ro + "/" + r + "/" + "compare" + "/" + bb + "..." + cb
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getInference(prePrompt string, commitBody string, apiKey string) (string, error) {
	prompt := prePrompt + commitBody

	client := openai.NewClient(apiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func generateCommitBody(bb string) (string, error) {
	commitMessages, err := getCommitMessages(bb)
	if err != nil {
		return "Error retrieving commit messages", err
	}
	summary := "# Summary of changes"
	return summary + "\n\nCommits:\n" + commitMessages, nil
}

func getCommitMessages(bb string) (string, error) {
	cmd := exec.Command("git", "log", `--pretty=format:(%h): %s`, bb+".."+"HEAD", "^origin", "-p")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func urlEncode(s string) string {
	return strings.ReplaceAll(s, " ", "+")
}
