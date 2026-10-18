package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const maxTitleLen = 20

var (
	ErrEmptyTitle   = errors.New("題が空です")
	ErrTitleTooLong = errors.New("題が長すぎます")
)

// Memo は、メモ1件を表す。
type Memo struct {
	Title     string
	Body      string
	Pinned    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMemo は、題を確かめてから新しいメモを作る。
func NewMemo(title, body string) (*Memo, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	now := time.Now()
	return &Memo{Title: title, Body: body, CreatedAt: now, UpdatedAt: now}, nil
}

// Edit は、メモの題と本文を書き換え、更新した日時を記録する。
func (m *Memo) Edit(title, body string) error {
	if err := validateTitle(title); err != nil {
		return err
	}
	m.Title = title
	m.Body = body
	m.UpdatedAt = time.Now()
	return nil
}

// Contains は、題か本文に word が含まれるかを返す。
func (m *Memo) Contains(word string) bool {
	return strings.Contains(m.Title, word) || strings.Contains(m.Body, word)
}

// Print は、メモ1件を区切りの線の下に表示する。
func (m *Memo) Print() {
	fmt.Println(strings.Repeat("-", 24))
	fmt.Printf("題: %s\n", m.Title)
	fmt.Printf("本文: %s\n", m.Body)
	if m.Pinned {
		fmt.Println("ピン留め: はい")
	} else {
		fmt.Println("ピン留め: いいえ")
	}
	fmt.Printf("更新: %s\n", m.UpdatedAt.Format("2006-01-02 15:04"))
}

func main() {
	titles := []string{"買い物リスト", "会議のメモ", "読みたい本", "週末の予定"}
	bodies := map[string]string{
		"買い物リスト": "牛乳と卵を買う",
		"会議のメモ":  "来週の会議で予算を決める",
		"読みたい本":  "Go の本を2冊",
		"週末の予定":  "本屋に行く",
	}

	var memos []*Memo
	for _, title := range titles {
		m, err := NewMemo(title, bodies[title])
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

	slices.SortFunc(memos, func(a, b *Memo) int {
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
		m.Print()
		found++
	}
	fmt.Printf("%d 件のメモ（全部で %d 件）\n", found, len(memos))
}

// validateTitle は、題が空でなく、maxTitleLen 文字以内かを確かめる。
func validateTitle(title string) error {
	n := utf8.RuneCountInString(title)
	if n == 0 {
		return ErrEmptyTitle
	}
	if n > maxTitleLen {
		return fmt.Errorf("%w（%d 文字。%d 文字まで）", ErrTitleTooLong, n, maxTitleLen)
	}
	return nil
}
