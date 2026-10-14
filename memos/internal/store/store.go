// Package store は、メモを保存する場所を表すインターフェースと、その実装をまとめる。
package store

import (
	"errors"

	"memos/internal/memo"
)

// ErrNotFound は、指定した ID のメモが無いときのエラー。
var ErrNotFound = errors.New("メモが見つかりません")

// Store は、メモを保存する場所。保存先ごとに、この4つのメソッドを持つ型を作る。
type Store interface {
	Save(m *memo.Memo) error
	Get(id string) (*memo.Memo, error)
	List() ([]*memo.Memo, error)
	Delete(id string) error
}
