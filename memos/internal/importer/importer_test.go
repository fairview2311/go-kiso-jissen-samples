package importer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"memos/internal/memo"
)

func TestParallelMapKeepsOrder(t *testing.T) {
	got, err := ParallelMap([]int{1, 2, 3, 4, 5}, 2, func(n int) (string, error) {
		return strings.Repeat("*", n), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"*", "**", "***", "****", "*****"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q、期待したのは %q", i, got[i], want[i])
		}
	}
}

func TestReadDir(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.txt", "買い物リスト\n牛乳と卵を買う\n")
	writeFile(t, dir, "b.txt", "会議のメモ\n来週の会議で予算を決める\n")
	writeFile(t, dir, "memo.md", "拡張子が違うので読まない\n")

	memos, err := ReadDir(dir, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(memos) != 2 || memos[0].Title != "買い物リスト" || memos[1].Body != "来週の会議で予算を決める" {
		t.Errorf("読んだメモ = %v", memos)
	}

	writeFile(t, dir, "c.txt", "\n題の無いファイル\n")
	if _, err := ReadDir(dir, 4); !errors.Is(err, memo.ErrEmptyTitle) {
		t.Errorf("題の無いファイルがあるときのエラー = %v、期待したのは ErrEmptyTitle", err)
	}
}

// BenchmarkReadDir は、1000 個のファイルを、同時に読む数を変えて読む時間を測る。
func BenchmarkReadDir(b *testing.B) {
	dir := b.TempDir()
	for i := range 1000 {
		path := filepath.Join(dir, fmt.Sprintf("%04d.txt", i))
		if err := os.WriteFile(path, []byte(fmt.Sprintf("メモ %d\n本文 %d\n", i, i)), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	for _, limit := range []int{1, 4, 8} {
		b.Run(fmt.Sprintf("limit=%d", limit), func(b *testing.B) {
			for b.Loop() {
				if _, err := ReadDir(dir, limit); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
