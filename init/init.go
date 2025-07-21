package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "init" {
		color.Yellow("Usage: pushAuto init -r <repo_url>")
		os.Exit(1)
	}

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	repo := initCmd.String("r", "", "GitHub repository URL (e.g., https://github.com/user/repo.git)")
	private := initCmd.Bool("private", false, "private")
	initCmd.Parse(os.Args[2:])

	if *repo == "" {
		color.Yellow("リポジトリURLは -r フラグで指定してください。")
		os.Exit(1)
	}

	cloneURL, err := createGitHubRepo(*repo, *private)
	if err != nil {
		color.Red("Githubリポジトリ作成失敗:%v", err)
		os.Exit(1)
	}
	color.Green("Githubリポジトリ作成成功:%s", cloneURL)

	run("git", "init")
	run("git", "secrets", "--install")
	run("git", "secrets", "--register-aws")
	color.Green("git-secrets を自動導入しました")
	run("git", "remote", "add", "origin", cloneURL)
	run("git", "add", ".")
	run("git", "commit", "-m", "init")
	run("git", "branch", "-M", "main")
	run("git", "push", "-u", "origin", "main")

	color.Green("Githubリポジトリ初期化、push完了")

}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "コマンド実行エラー: %v\n", err)
		os.Exit(1)
	}
}

type RepoResponse struct {
	CloneURL string `json:"clone_url"`
	HTMLURL  string `json:"html_url"`
	Name     string `json:"name"`
}

func createGitHubRepo(repoName string, private bool) (string, error) {
	token := loadToken()
	if token == "" {
		return "", errors.New("環境変数 GITHUB_TOKENが設定されていません。")
	}

	body := map[string]interface{}{
		"name":    repoName,
		"private": private,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/user/repos", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("github API Error:%s", resp.Status)
	}

	var result RepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.CloneURL, nil
}

func loadToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}

	home, err := os.UserHomeDir()
	if err == nil {
		err = godotenv.Load(filepath.Join(home, ".env"))
		if err == nil {
			token = os.Getenv("GITHUB_TOKEN")
			if token != "" {
				return token
			}
		}
	}

	log.Fatal("GITHUB_TOKENが設定されていません。")
	return ""
}
