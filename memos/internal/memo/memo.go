// Package memo は、memos で扱うメモ1件を表す型と、その確かめをまとめる。
package memo

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"
)

// MaxTitleLen は、題の長さの上限（文字数）。
const MaxTitleLen = 20

var (
	ErrEmptyTitle   = errors.New("題が空です")
	ErrTitleTooLong = errors.New("題が長すぎます")
)

// Memo は、メモ1件を表す。
type Memo struct {
	ID        string
	Title     string
	Body      string
	Pinned    bool
	CreatedAt time.Time
	UpdatedAt time.Time
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
