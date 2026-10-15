# Go言語 実践入門 練習用ファイル

書籍『Go言語 実践入門 基礎文法からWeb APIとWebアプリを1本作るまで』（早瀬 透）の練習用ファイルです。

本では、メモ帳のアプリ「memos」を1冊を通して育てます。このリポジトリには、**各章を終えた時点の memos** が入っています。
本文のコードの枠と、その章の練習用ファイルは同じ内容です。

- 書いたコードが動かないとき、本のとおりの完成形と見比べる
- 途中の章から読み始めるとき、前の章の終わりの状態をそろえる

といった使い方をしてください。

## ファイルの取り方

取り方は3つあります。やりたいことに合わせて選んでください。

| やりたいこと | 取り方 |
|---|---|
| 1つのファイルだけ見比べたい・コピーしたい | 下の表の［ファイルを見る］から、ブラウザでファイルを開く。右上のコピーのボタンで中身をコピーできます |
| ある章の終わりの memos を丸ごと手元に置きたい | 下の表の［ZIP］を押してダウンロードし、展開する |
| 全部の章を手元に置いて、章を切り替えたい | `git clone` でクローンし、`git checkout ch05` のように切り替える（下の「git でクローンする」） |

ZIP を展開したり git でクローンしたりしたら、中の `memos` のフォルダーを、VS Code の［ファイル］メニューの［フォルダーを開く...］（Mac は［開く...］）で開きます。

| 章 | この章の終わりの memos | ファイルを見る | ダウンロード |
|---|---|---|---|
| 第1章 | 「こんにちは、memos です」と表示する `main.go` | [ch01](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch01/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch01.zip) |
| 第2章 | VS Code で整形した `main.go` | [ch02](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch02/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch02.zip) |
| 第3章 | メモ1件を表示する | [ch03](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch03/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch03.zip) |
| 第4章 | 題を確かめる関数とエラー | [ch04](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch04/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch04.zip) |
| 第5章 | メモの一覧と、キーワードでの検索 | [ch05](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch05/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch05.zip) |
| 第6章 | Memo 型とメソッド | [ch06](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch06/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch06.zip) |
| 第7章 | `cmd/memos` と `internal/memo` にファイルを分けた形 | [ch07](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch07/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch07.zip) |
| 第8章 | 保存先の Store インターフェース（MemoryStore・FileStore） | [ch08](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch08/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch08.zip) |
| 第9章 | memos.json に保存する add・list・delete | [ch09](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch09/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch09.zip) |
| 第10章 | テスト付きのコマンド | [ch10](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch10/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch10.zip) |
| 第11章 | ファイルをまとめて取り込む import | [ch11](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch11/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch11.zip) |
| 第12章 | Web API（serve） | [ch12](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch12/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch12.zip) |
| 第13章 | SQLite に保存する Web API | [ch13](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch13/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch13.zip) |
| 第14章 | ブラウザで使える画面 | [ch14](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch14/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch14.zip) |
| 第15章 | テスト・ログ付きの完成形（このページの `memos` と同じ） | [ch15](https://github.com/fairview2311/go-kiso-jissen-samples/tree/ch15/memos) | [ZIP](https://github.com/fairview2311/go-kiso-jissen-samples/archive/refs/tags/ch15.zip) |

### git でクローンする

Git をインストールしている人は、ターミナルで次のように打ちます。

```text
git clone https://github.com/fairview2311/go-kiso-jissen-samples.git
cd go-kiso-jissen-samples
git checkout ch05
```

`ch01`〜`ch15` は Git のタグ（ある時点のファイルに付けた名前）です。`git checkout ch05` で第5章の終わりの状態に、`git checkout main` で完成形に戻ります。自分で書き換えたファイルは、切り替える前に別の場所へ写しておいてください。

## 動かし方

Go 1.27 をインストールしてから、`memos` のフォルダーを VS Code で開き、［ターミナル］メニューの［新しいターミナル］で開いたターミナルで打ちます。

| 章 | 動かすコマンド |
|---|---|
| 第1〜6章 | `go run .` |
| 第7〜11章 | `go run ./cmd/memos`（第9章からは `go run ./cmd/memos list` のようにサブコマンドを付けます） |
| 第12〜15章 | `go run ./cmd/memos serve` を打ち、ブラウザで http://localhost:8080/ を開きます。止めるときは Ctrl＋C |

テストは、第10章からどの章でも `go test ./...` で動きます。

第13章から使う SQLite のドライバ（modernc.org/sqlite）は、初めて動かしたときに自動でダウンロードされます。インターネットにつながった状態で動かしてください。

## 完成形のフォルダーの中身（第15章）

```text
memos/
├── cmd/memos/          コマンドの入り口（main.go）
├── internal/memo/      メモの型（Memo）
├── internal/store/     保存先（メモリ・JSON のファイル・SQLite）
├── internal/importer/  ファイルをまとめて取り込む処理
├── internal/web/       Web API と画面（テンプレート・CSS）
└── testdata/           curl で送る JSON と、取り込みの見本のファイル
```

## うまくいかないとき

- `go version` が 1.27 より古い … 本の第1章の手順で Go 1.27 をインストールし、VS Code を開き直します
- Windows で `curl` の結果がおかしい … Windows の PowerShell では `curl` ではなく `curl.exe` と打ちます（第12章）
- Windows で `go test -race` が動かない … `-race` には C のコンパイラ（gcc）が必要です。付けずに `go test ./...` で確かめます（第11章）

## 確かめたバージョン

- Go 1.27.2
- modernc.org/sqlite v1.60.1・github.com/mattn/go-runewidth v0.0.31
- Windows 11・macOS（2026 年 10 月）

## ライセンス

MIT License（[LICENSE](LICENSE) を見てください）。本のコードは、自分のアプリに自由に使ってかまいません。
