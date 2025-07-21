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


## 使用技術
- 言語：Go(標準ライブラリ中心)
- API連携：GitHub REST API v3
- ログ可視化：faith/colorによるCLIカラー出力
- セキュリティ強化：git-secrets による AWSシークレット流出対策


## 使用例
```bash
# GitHubリポジトリ新規作成 + 初期化 + push
pushAutoInit init -r my-new-repo --private

# 差分があればコミット＆プッシュ
pushAutoPush push -m "fix: update documentation"
```

## 設計方針
- 手動リポジトリ初期化作業の反復削減が目的
- git-secrets 導入、.env によるAPIキー管理により安全性向上

## 今後の改善予定
- pushAuto への統合（単一バイナリ化）
- cobraによるサブコマンドCLI化
- --help オプション自動生成
- .gitignore 自動生成対応