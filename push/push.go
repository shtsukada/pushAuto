package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "push" {
		color.Yellow("Usage: pushAuto push -m \"commit message\"")
		os.Exit(1)
	}

	pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
	message := pushCmd.String("m", "", "commit message")
	pushCmd.Parse(os.Args[2:])

	if *message == "" {
		color.Yellow("コミットメッセージは -m フラグで指定してください。")
		os.Exit(1)
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		color.Yellow("このディレクトリはGitリポジトリではありません。")
		os.Exit(1)
	}

	run("git", "add", ".")

	if isStagingAreaEmpty() {
		color.Yellow("ステージされた変更がありません。コミットをスキップします。")
		return
	}

	run("git", "commit", "-m", *message)
	run("git", "push")
}

func isStagingAreaEmpty() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	return err == nil
}

func run(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "コマンド実行エラー: %v\n", err)
		os.Exit(1)
	}
}
