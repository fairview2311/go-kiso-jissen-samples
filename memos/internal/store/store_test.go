package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"memos/internal/memo"
)

// TestStores は、同じテストを MemoryStore と FileStore の両方で動かす。
func TestStores(t *testing.T) {
	stores := []struct {
		name string
		s    Store
	}{
		{"MemoryStore", NewMemoryStore()},
		{"FileStore", NewFileStore(filepath.Join(t.TempDir(), "memos.json"))},
		{"SqlStore", newTestSqlStore(t)},
	}
	for _, tt := range stores {
		t.Run(tt.name, func(t *testing.T) {
			testStore(t, tt.s)
		})
	}
}

// testStore は、Store の4つのメソッドが約束どおりに動くかを確かめる。
func testStore(t *testing.T, s Store) {
	t.Helper()
	list, err := s.List()
	if err != nil || len(list) != 0 {
		t.Fatalf("最初の List = %d 件, %v、期待したのは 0 件", len(list), err)
	}

	m, err := memo.New("買い物リスト", "牛乳と卵を買う")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Save(m); err != nil {
		t.Fatalf("Save のエラー = %v", err)
	}
	got, err := s.Get(m.ID)
	if err != nil {
		t.Fatalf("Get のエラー = %v", err)
	}
	if got.Title != "買い物リスト" {
		t.Errorf("Get した題 = %q、期待したのは %q", got.Title, "買い物リスト")
	}

	if err := s.Delete(m.ID); err != nil {
		t.Fatalf("Delete のエラー = %v", err)
	}
	if _, err := s.Get(m.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("消した後の Get のエラー = %v、期待したのは ErrNotFound", err)
	}
	if err := s.Delete(m.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("2回目の Delete のエラー = %v、期待したのは ErrNotFound", err)
	}
}

// TestFileStoreKeepsMemos は、FileStore に保存したメモが、別の FileStore からも読めることを確かめる。
func TestFileStoreKeepsMemos(t *testing.T) {
	path := filepath.Join(t.TempDir(), "memos.json")
	m, err := memo.New("会議のメモ", "来週の会議で予算を決める")
	if err != nil {
		t.Fatal(err)
	}
	if err := NewFileStore(path).Save(m); err != nil {
		t.Fatal(err)
	}
	got, err := NewFileStore(path).Get(m.ID)
	if err != nil {
		t.Fatalf("別の FileStore からの Get のエラー = %v", err)
	}
	if got.Body != m.Body || !got.CreatedAt.Equal(m.CreatedAt) {
		t.Errorf("読み直したメモ = %q %v、期待したのは %q %v", got.Body, got.CreatedAt, m.Body, m.CreatedAt)
	}
}

// newTestSqlStore は、テスト用の SQLite のファイルを一時ディレクトリに作る。テストが終わると閉じる。
func newTestSqlStore(t *testing.T) *SqlStore {
	t.Helper()
	s, err := NewSqlStore(filepath.Join(t.TempDir(), "memos.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// TestSqlStoreConcurrentSave は、50 個の goroutine から同時に保存しても、50 件とも残ることを確かめる。
func TestSqlStoreConcurrentSave(t *testing.T) {
	s := newTestSqlStore(t)
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			m, err := memo.New(fmt.Sprintf("メモ %d", i), "同時に保存する")
			if err != nil {
				t.Error(err)
				return
			}
			if err := s.Save(m); err != nil {
				t.Errorf("Save のエラー = %v", err)
			}
		})
	}
	wg.Wait()
	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 50 {
		t.Errorf("保存できたメモ = %d 件、期待したのは 50 件", len(list))
	}
}
