package memo

import (
	"errors"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{name: "ふつうの題", title: "買い物リスト", wantErr: nil},
		{name: "ちょうど20文字", title: strings.Repeat("あ", 20), wantErr: nil},
		{name: "空の題", title: "", wantErr: ErrEmptyTitle},
		{name: "21文字", title: strings.Repeat("あ", 21), wantErr: ErrTitleTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := New(tt.title, "本文")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("New(%q) のエラー = %v、期待したのは %v", tt.title, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if m.Title != tt.title {
				t.Errorf("Title = %q、期待したのは %q", m.Title, tt.title)
			}
			if len(m.ID) != 36 {
				t.Errorf("ID の長さ = %d、期待したのは 36", len(m.ID))
			}
		})
	}
}

func TestEdit(t *testing.T) {
	m := mustNew(t, "買い物リスト", "牛乳と卵を買う")
	before := m.UpdatedAt

	if err := m.Edit("買い物リスト（土曜）", "牛乳と卵とパンを買う"); err != nil {
		t.Fatalf("Edit のエラー = %v", err)
	}
	if m.Title != "買い物リスト（土曜）" || m.Body != "牛乳と卵とパンを買う" {
		t.Errorf("書き換わっていない: %q %q", m.Title, m.Body)
	}
	if m.UpdatedAt.Before(before) {
		t.Errorf("UpdatedAt が前の日時より前: %v", m.UpdatedAt)
	}

	if err := m.Edit("", "本文"); !errors.Is(err, ErrEmptyTitle) {
		t.Errorf("空の題で Edit したエラー = %v、期待したのは ErrEmptyTitle", err)
	}
	if m.Title != "買い物リスト（土曜）" {
		t.Errorf("エラーなのに題が変わった: %q", m.Title)
	}
}

func TestContains(t *testing.T) {
	m := mustNew(t, "読みたい本", "Go の本を2冊")
	tests := []struct {
		word string
		want bool
	}{
		{"本", true},
		{"Go", true},
		{"カレー", false},
	}
	for _, tt := range tests {
		if got := m.Contains(tt.word); got != tt.want {
			t.Errorf("Contains(%q) = %v、期待したのは %v", tt.word, got, tt.want)
		}
	}
}

// mustNew は、テスト用のメモを作る。作れなければテストをそこで止める。
func mustNew(t *testing.T, title, body string) *Memo {
	t.Helper()
	m, err := New(title, body)
	if err != nil {
		t.Fatalf("New(%q) のエラー = %v", title, err)
	}
	return m
}
