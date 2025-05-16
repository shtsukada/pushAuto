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

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".envの読み込みに失敗しました。")
	}

	//サブコマンド判定
	if len(os.Args) < 2 || os.Args[1] != "init" {
		color.Yellow("Usage: pushAuto init -r <repo_url>")
		os.Exit(1)
	}

	//initコマンド用のフラグ定義
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

	//コマンド実行
	run("git", "init")
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
	token := os.Getenv("GITHUB_TOKEN")
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

	req, err := http.NewRequest("POST", "https://github.com/shtsukada/repos", bytes.NewBuffer(jsonBody))
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
		return "", fmt.Errorf("Github API Error:%s", resp.Status)
	}

	var result RepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.CloneURL, nil
}
