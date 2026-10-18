package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode/utf8"
)

const maxTitleLen = 20

var (
	ErrEmptyTitle   = errors.New("題が空です")
	ErrTitleTooLong = errors.New("題が長すぎます")
)

func main() {
	titles := []string{"買い物リスト", "会議のメモ", "読みたい本"}
	bodies := map[string]string{
		"買い物リスト": "牛乳と卵を買う",
		"会議のメモ":  "来週の会議で予算を決める",
		"読みたい本":  "Go の本を2冊",
	}
	pinned := map[string]bool{
		"会議のメモ": true,
	}

	titles = append(titles, "週末の予定")
	bodies["週末の予定"] = "本屋に行く"

	for _, title := range titles {
		if err := validateTitle(title); err != nil {
			fmt.Println("エラー:", err)
			os.Exit(1)
		}
	}

	slices.SortFunc(titles, func(a, b string) int {
		if pinned[a] != pinned[b] {
			if pinned[a] {
				return -1
			}
			return 1
		}
		return strings.Compare(a, b)
	})

	word := ""
	if len(os.Args) > 1 {
		word = os.Args[1]
	}
	found := 0
	for _, title := range titles {
		body := bodies[title]
		if word != "" && !strings.Contains(title, word) && !strings.Contains(body, word) {
			continue
		}
		printMemo(title, body, pinned[title])
		found++
	}
	fmt.Printf("%d 件のメモ（全部で %d 件）\n", found, len(titles))
}

// validateTitle は、題が空でなく、maxTitleLen 文字以内かを確かめる。
func validateTitle(title string) error {
	n, _ := titleLength(title)
	if n == 0 {
		return ErrEmptyTitle
	}
	if n > maxTitleLen {
		return fmt.Errorf("%w（%d 文字。%d 文字まで）", ErrTitleTooLong, n, maxTitleLen)
	}
	return nil
}

// titleLength は、題の文字数とバイト数を返す。
func titleLength(title string) (int, int) {
	return utf8.RuneCountInString(title), len(title)
}

// printMemo は、メモ1件を区切りの線の下に表示する。
func printMemo(title, body string, pinned bool) {
	for i := 0; i < 24; i++ {
		fmt.Print("-")
	}
	fmt.Println()
	fmt.Printf("題: %s\n", title)
	fmt.Printf("本文: %s\n", body)
	if pinned {
		fmt.Println("ピン留め: はい")
	} else {
		fmt.Println("ピン留め: いいえ")
	}
}
