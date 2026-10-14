package store

import (
	"fmt"
	"maps"
	"slices"

	"memos/internal/memo"
)

// MemoryStore は、メモをメモリーのマップに置く保存先。プログラムを終えると消える。
type MemoryStore struct {
	memos map[string]*memo.Memo
}

// MemoryStore が Store のメソッドをすべて持っているかを、コンパイルのときに確かめる。
var _ Store = (*MemoryStore)(nil)

// NewMemoryStore は、空の MemoryStore を作る。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{memos: map[string]*memo.Memo{}}
}

// Save は、メモを保存する。同じ ID のメモがあれば置き換える。
func (s *MemoryStore) Save(m *memo.Memo) error {
	s.memos[m.ID] = m
	return nil
}

// Get は、ID のメモを返す。無ければ ErrNotFound をラップして返す。
func (s *MemoryStore) Get(id string) (*memo.Memo, error) {
	m, ok := s.memos[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return m, nil
}

// List は、すべてのメモを並べ替えて返す。
func (s *MemoryStore) List() ([]*memo.Memo, error) {
	list := slices.Collect(maps.Values(s.memos))
	memo.Sort(list)
	return list, nil
}

// Delete は、ID のメモを消す。無ければ ErrNotFound をラップして返す。
func (s *MemoryStore) Delete(id string) error {
	if _, ok := s.memos[id]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	delete(s.memos, id)
	return nil
}
