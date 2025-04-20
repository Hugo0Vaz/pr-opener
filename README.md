# pr-opener
A GitHub PR opener script

## Usage

The pr-opener project allows you to generate a GitHub pull request URL automatically from your current git branch. It gathers commit messages, leverages the OpenAI API to generate a PR title and description, and assembles the final URL.

### Prerequisites

- Go and make, must be installed.
- The `OPENAI_API_KEY` environment variable needs to be set.
- You must be working within a git repository.
- Optionally, create a `.propener.toml` file in the project root with the following format:

```toml
[repo]
base = "https://github.com/"
repo_owner = "Hugo0Vaz"
repo = "pr-opener"
base_branch = "main"
```

# Installation

Run `make` to build the project and produce the binary `pr-opener`. To install, run:

    sudo make install

# Usage

## Using Configuration File (.propener.toml):

Simply run the project without providing CLI flags:

```bash
propener
```
The tool will load the configuration from the `.propener.toml` file if no flags are provided.

## Using CLI Flags:

Override the configuration by passing flags. For example:

```bash
propener --repo-owner "Hugo0Vaz" --repo "pr-opener" --base-branch "main"
```

The tool will:
1. Determine your current git branch.
2. Generate a list of commit messages.
3. Use the OpenAI API to infer a PR title and detailed description.
4. Assemble and print a GitHub PR URL.

### Example

For a branch named `feature/new-feature`, the tool might produce a URL similar to:

```
https://github.com/Hugo0Vaz/pr-opener/compare/main...feature/new-feature?quick_pull=1&title=<generated_title>
```

Enjoy using pr-opener!

