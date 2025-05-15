package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func main() {
	//サブコマンド判定
	if len(os.Args) < 2 || os.Args[1] != "init" {
		color.Yellow("Usage: pushAuto init -r <repo_url>")
		os.Exit(1)
	}

	//initコマンド用のフラグ定義
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	repo := initCmd.String("r", "", "GitHub repository URL (e.g., https://github.com/user/repo.git)")
	initCmd.Parse(os.Args[2:])

	if *repo == "" {
		color.Yellow("リポジトリURLは -r フラグで指定してください。")
		os.Exit(1)
	}

	//既にリポジトリである場合は警告
	if _, err := os.Stat(".git"); err == nil {
		color.Yellow("既にこのディレクトリはGitリポジトリです。")
		os.Exit(1)
	}

	//コマンド実行
	run("git", "init")
	run("git", "remote", "add", "origin", *repo)
	run("git", "add", ".")
	run("git", "commit", "-m", "init")
	run("git", "branch", "-M", "main")
	run("git", "push", "-u", "origin", "main")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)

	null, _ := os.Open(os.DevNull)
	defer null.Close()

	cmd.Stdout = null
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "コマンド実行エラー: %v\n", err)
		os.Exit(1)
	}
}
