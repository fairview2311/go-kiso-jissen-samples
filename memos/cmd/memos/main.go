// memos は、メモを SQLite のデータベース memos.db に保存するコマンド。
//
//	memos add <題> <本文>    メモを足す
//	memos list [言葉]        一覧を表示する。言葉を渡すと、その言葉を含むメモだけ
//	memos delete <ID の先頭>  メモを消す
//	memos import <ディレクトリ> ディレクトリの .txt のファイルを、メモとしてまとめて取り込む
//	memos serve [アドレス]     Web API と画面を動かす（アドレスの既定は localhost:8080）
//	memos version             バージョンを表示する
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"memos/internal/importer"
	"memos/internal/memo"
	"memos/internal/store"
	"memos/internal/web"
)

const dataFile = "memos.db"

// version は、memos のバージョン。go build の -ldflags で書き換える。
var version = "dev"

func main() {
	s, err := store.NewSqlStore(dataFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
	err = run(os.Stdout, s, os.Args[1:])
	s.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
}

// run は、args の最初の言葉で処理を選び、結果を w に書き出す。
func run(w io.Writer, s store.Store, args []string) error {
	if len(args) == 0 {
		return list(w, s, "")
	}
	switch args[0] {
	case "add":
		if len(args) != 3 {
			return errors.New("使い方: memos add <題> <本文>")
		}
		return add(w, s, args[1], args[2])
	case "list":
		word := ""
		if len(args) > 1 {
			word = args[1]
		}
		return list(w, s, word)
	case "delete":
		if len(args) != 2 {
			return errors.New("使い方: memos delete <ID の先頭>")
		}
		return remove(w, s, args[1])
	case "import":
		if len(args) != 2 {
			return errors.New("使い方: memos import <ディレクトリ>")
		}
		return importDir(w, s, args[1])
	case "serve":
		addr := "localhost:8080"
		if len(args) > 1 {
			addr = args[1]
		}
		return serve(s, addr)
	case "version":
		fmt.Fprintln(w, "memos", version)
		return nil
	default:
		return fmt.Errorf("知らないコマンドです: %s（add・list・delete・import・serve・version のどれか）", args[0])
	}
}

// add は、新しいメモを作って保存する。
func add(w io.Writer, s store.Store, title, body string) error {
	m, err := memo.New(title, body)
	if err != nil {
		return err
	}
	if err := s.Save(m); err != nil {
		return err
	}
	fmt.Fprintf(w, "追加しました: %s %s\n", m.ID[:8], m.Title)
	return nil
}

// list は、word を含むメモを1行ずつ書き出す。word が空なら全部。
func list(w io.Writer, s store.Store, word string) error {
	memos, err := s.List()
	if err != nil {
		return err
	}
	found := 0
	for _, m := range memos {
		if word != "" && !m.Contains(word) {
			continue
		}
		fmt.Fprintln(w, m)
		found++
	}
	fmt.Fprintf(w, "%d 件のメモ（全部で %d 件）\n", found, len(memos))
	return nil
}

// remove は、ID が prefix で始まるメモを1件消す。
func remove(w io.Writer, s store.Store, prefix string) error {
	memos, err := s.List()
	if err != nil {
		return err
	}
	var hit []*memo.Memo
	for _, m := range memos {
		if strings.HasPrefix(m.ID, prefix) {
			hit = append(hit, m)
		}
	}
	switch len(hit) {
	case 0:
		return fmt.Errorf("%w: %s", store.ErrNotFound, prefix)
	case 1:
		if err := s.Delete(hit[0].ID); err != nil {
			return err
		}
		fmt.Fprintf(w, "消しました: %s %s\n", hit[0].ID[:8], hit[0].Title)
		return nil
	default:
		return fmt.Errorf("%s で始まるメモが %d 件あります。ID をもっと長く書いてください", prefix, len(hit))
	}
}

// importDir は、dir の .txt のファイルを同時に読んでメモにし、保存する。
func importDir(w io.Writer, s store.Store, dir string) error {
	start := time.Now()
	memos, err := importer.ReadDir(dir, 8)
	if err != nil {
		return err
	}
	for _, m := range memos {
		if err := s.Save(m); err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "%d 件を取り込みました（%v）\n", len(memos), time.Since(start).Round(time.Millisecond))
	return nil
}

// serve は、addr で Web API と画面を動かす。Ctrl+C を押すと、処理中のリクエストを終えてから止まる。
func serve(s store.Store, addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: addr, Handler: web.NewServer(s).Handler()}
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()
	slog.Info("待ち受けを始めました", "url", "http://"+addr, "version", version)

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	slog.Info("止めています")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	slog.Info("止めました")
	return nil
}
