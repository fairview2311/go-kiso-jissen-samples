// memos は、メモの一覧を表示するコマンド。引数に言葉を渡すと、その言葉を含むメモだけを表示する。
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mattn/go-runewidth"

	"memos/internal/memo"
)

func main() {
	titles := []string{"買い物リスト", "会議のメモ", "読みたい本", "週末の予定"}
	bodies := map[string]string{
		"買い物リスト": "牛乳と卵を買う",
		"会議のメモ":  "来週の会議で予算を決める",
		"読みたい本":  "Go の本を2冊",
		"週末の予定":  "本屋に行く",
	}

	var memos []*memo.Memo
	for _, title := range titles {
		m, err := memo.New(title, bodies[title])
		if err != nil {
			fmt.Println("エラー:", err)
			os.Exit(1)
		}
		memos = append(memos, m)
	}
	memos[1].Pinned = true
	if err := memos[0].Edit("買い物リスト（土曜）", "牛乳と卵とパンを買う"); err != nil {
		fmt.Println("エラー:", err)
		os.Exit(1)
	}

	slices.SortFunc(memos, func(a, b *memo.Memo) int {
		if a.Pinned != b.Pinned {
			if a.Pinned {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Title, b.Title)
	})

	word := ""
	if len(os.Args) > 1 {
		word = os.Args[1]
	}
	found := 0
	for _, m := range memos {
		if word != "" && !m.Contains(word) {
			continue
		}
		printRow(m)
		found++
	}
	fmt.Printf("%d 件のメモ（全部で %d 件）\n", found, len(memos))
}

// printRow は、メモ1件を1行で表示する。題は表示の幅で 22 にそろえる。
func printRow(m *memo.Memo) {
	mark := " "
	if m.Pinned {
		mark = "*"
	}
	fmt.Printf("%s %s %s  %s\n", mark, m.ID[:8], runewidth.FillRight(m.Title, 22), m.Body)
}
