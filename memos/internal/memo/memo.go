// Package memo は、memos で扱うメモ1件を表す型と、その確かめをまとめる。
package memo

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"

	"github.com/mattn/go-runewidth"
)

// MaxTitleLen は、題の長さの上限（文字数）。
const MaxTitleLen = 20

var (
	ErrEmptyTitle   = errors.New("題が空です")
	ErrTitleTooLong = errors.New("題が長すぎます")
)

// Memo は、メモ1件を表す。
type Memo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Pinned    bool      `json:"pinned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New は、題を確かめてから新しいメモを作る。ID は UUID で付ける。
func New(title, body string) (*Memo, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	now := time.Now()
	return &Memo{
		ID:        uuid.New().String(),
		Title:     title,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
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

// String は、メモを1行の文字列にする。fmt.Println などは、この形で表示する。
func (m *Memo) String() string {
	mark := " "
	if m.Pinned {
		mark = "*"
	}
	return fmt.Sprintf("%s %s %s  %s", mark, m.ID[:8], runewidth.FillRight(m.Title, 22), m.Body)
}

// Sort は、ピン留めしたメモを先頭に、残りを題の順に並べ替える。
func Sort(memos []*Memo) {
	slices.SortFunc(memos, func(a, b *Memo) int {
		if a.Pinned != b.Pinned {
			if a.Pinned {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Title, b.Title)
	})
}

// validateTitle は、題が空でなく、MaxTitleLen 文字以内かを確かめる。
func validateTitle(title string) error {
	n := utf8.RuneCountInString(title)
	if n == 0 {
		return ErrEmptyTitle
	}
	if n > MaxTitleLen {
		return fmt.Errorf("%w（%d 文字。%d 文字まで）", ErrTitleTooLong, n, MaxTitleLen)
	}
	return nil
}
