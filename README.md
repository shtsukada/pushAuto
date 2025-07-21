# pushAuto

## 概要
pushAutoはGitHubリポジトリの初期化及びコミット、プッシュ作業を効率化するためのCLIツールです。GitHub APIを活用し、ローカル作業からリポジトリ作成、初回プッシュ、通常のコミット、プッシュを自動化できます。<br>このツールは日常的なGit作業の繰り返しを削減し、API連携、CLI設計、認証管理のスキル習得を目的に作成しました。

## 機能一覧
| 機能 | 説明 |
|------|------|
| `init`| GitHub APIを使用してリポジトリを作成し、初回Pushまで自動化 |
| `push`| 変更差分がある場合のみcommit→pushを実行 |
| `--private`| `init`実行時にプライベートリポジトリの作成が可能 |
| `.env`| `GITHUB_TOKEN` は `.env` または環境変数で安全に管理 |
| 差分チェック | 空コミット、空pushを防止(`git diff --cached --quiet` 使用) |
| カラー出力| `faith/color`を用いた視認性の高いログ出力 |





## 前提条件

以下のツールがインストールされている必要があります：

- git
- git-secrets（https://github.com/awslabs/git-secrets）

macOSの場合は以下でインストール可能です：
```bash
brew install git-secrets
```
