package store

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"

	"memos/internal/memo"
)

// FileStore は、メモを JSON のファイルに保存する保存先。プログラムを終えても残る。
type FileStore struct {
	path string
}

var _ Store = (*FileStore)(nil)

// NewFileStore は、path のファイルに保存する FileStore を作る。ファイルはまだ無くてよい。
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

// load は、ファイルからすべてのメモを読み込む。ファイルが無ければ空のスライスを返す。
func (s *FileStore) load() ([]*memo.Memo, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var memos []*memo.Memo
	if err := json.Unmarshal(data, &memos); err != nil {
		return nil, fmt.Errorf("%s を読めません: %w", s.path, err)
	}
	return memos, nil
}

// write は、すべてのメモを JSON にしてファイルに書く。
func (s *FileStore) write(memos []*memo.Memo) error {
	data, err := json.Marshal(memos, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Save は、メモを保存する。同じ ID のメモがあれば置き換える。
func (s *FileStore) Save(m *memo.Memo) error {
	memos, err := s.load()
	if err != nil {
		return err
	}
	i := slices.IndexFunc(memos, func(x *memo.Memo) bool { return x.ID == m.ID })
	if i >= 0 {
		memos[i] = m
	} else {
		memos = append(memos, m)
	}
	return s.write(memos)
}

// Get は、ID のメモを返す。無ければ ErrNotFound をラップして返す。
func (s *FileStore) Get(id string) (*memo.Memo, error) {
	memos, err := s.load()
	if err != nil {
		return nil, err
	}
	for _, m := range memos {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
}

// List は、すべてのメモを並べ替えて返す。
func (s *FileStore) List() ([]*memo.Memo, error) {
	memos, err := s.load()
	if err != nil {
		return nil, err
	}
	memo.Sort(memos)
	return memos, nil
}

// Delete は、ID のメモを消す。無ければ ErrNotFound をラップして返す。
func (s *FileStore) Delete(id string) error {
	memos, err := s.load()
	if err != nil {
		return err
	}
	i := slices.IndexFunc(memos, func(x *memo.Memo) bool { return x.ID == id })
	if i < 0 {
		return fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return s.write(slices.Delete(memos, i, i+1))
}
