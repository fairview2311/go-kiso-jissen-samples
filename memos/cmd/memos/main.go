// memos は、メモの一覧を表示するコマンド。引数に言葉を渡すと、その言葉を含むメモだけを表示する。
package main

import (
	"fmt"
	"io"
	"os"

	"memos/internal/memo"
	"memos/internal/store"
)

func main() {
	s := store.NewMemoryStore()
	if err := seed(s); err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
	word := ""
	if len(os.Args) > 1 {
		word = os.Args[1]
	}
	if err := run(os.Stdout, s, word); err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
}

// seed は、練習用のメモを4件作って s に保存する。
func seed(s store.Store) error {
	titles := []string{"買い物リスト", "会議のメモ", "読みたい本", "週末の予定"}
	bodies := map[string]string{
		"買い物リスト": "牛乳と卵を買う",
		"会議のメモ":  "来週の会議で予算を決める",
		"読みたい本":  "Go の本を2冊",
		"週末の予定":  "本屋に行く",
	}
	for i, title := range titles {
		m, err := memo.New(title, bodies[title])
		if err != nil {
			return err
		}
		if i == 0 {
			if err := m.Edit("買い物リスト（土曜）", "牛乳と卵とパンを買う"); err != nil {
				return err
			}
		}
		m.Pinned = i == 1
		if err := s.Save(m); err != nil {
			return err
		}
	}
	return nil
}

// run は、s のメモのうち word を含むものを w に1行ずつ書き出す。word が空なら全部。
func run(w io.Writer, s store.Store, word string) error {
	list, err := s.List()
	if err != nil {
		return err
	}
	found := 0
	for _, m := range list {
		if word != "" && !m.Contains(word) {
			continue
		}
		fmt.Fprintln(w, m)
		found++
	}
	fmt.Fprintf(w, "%d 件のメモ（全部で %d 件）\n", found, len(list))
	return nil
}
