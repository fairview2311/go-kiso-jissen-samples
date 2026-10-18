package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // database/sql に "sqlite" のドライバを登録する

	"memos/internal/memo"
)

// SqlStore は、メモを SQLite のデータベースに保存する保存先。
type SqlStore struct {
	db *sql.DB
}

var _ Store = (*SqlStore)(nil)

const schema = `
CREATE TABLE IF NOT EXISTS memos (
	id         TEXT PRIMARY KEY,
	title      TEXT NOT NULL,
	body       TEXT NOT NULL,
	pinned     INTEGER NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
)`

// NewSqlStore は、path の SQLite のファイルを開き、表が無ければ作る。
// busy_timeout と WAL は、同時に書き込んだときに待つための設定。接続ごとに効くよう、DSN に書く。
func NewSqlStore(path string) (*SqlStore, error) {
	dsn := path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("表を作れません: %w", err)
	}
	return &SqlStore{db: db}, nil
}

// Close は、データベースを閉じる。
func (s *SqlStore) Close() error {
	return s.db.Close()
}

// Save は、メモを保存する。同じ ID のメモがあれば書き換える。
func (s *SqlStore) Save(m *memo.Memo) error {
	_, err := s.db.Exec(`
		INSERT INTO memos (id, title, body, pinned, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			title = excluded.title, body = excluded.body, pinned = excluded.pinned,
			updated_at = excluded.updated_at`,
		m.ID, m.Title, m.Body, m.Pinned,
		m.CreatedAt.Format(time.RFC3339Nano), m.UpdatedAt.Format(time.RFC3339Nano))
	return err
}

// Get は、ID のメモを返す。無ければ ErrNotFound をラップして返す。
func (s *SqlStore) Get(id string) (*memo.Memo, error) {
	row := s.db.QueryRow(`SELECT id, title, body, pinned, created_at, updated_at FROM memos WHERE id = ?`, id)
	m, err := scanMemo(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return m, err
}

// List は、すべてのメモを並べ替えて返す。
func (s *SqlStore) List() ([]*memo.Memo, error) {
	rows, err := s.db.Query(`SELECT id, title, body, pinned, created_at, updated_at FROM memos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var memos []*memo.Memo
	for rows.Next() {
		m, err := scanMemo(rows)
		if err != nil {
			return nil, err
		}
		memos = append(memos, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	memo.Sort(memos)
	return memos, nil
}

// Delete は、ID のメモを消す。無ければ ErrNotFound をラップして返す。
func (s *SqlStore) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM memos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return nil
}

// scanner は、*sql.Row と *sql.Rows の両方が持つ Scan のメソッドだけのインターフェース。
type scanner interface {
	Scan(dest ...any) error
}

// scanMemo は、1行を読んでメモにする。
func scanMemo(r scanner) (*memo.Memo, error) {
	var m memo.Memo
	var created, updated string
	if err := r.Scan(&m.ID, &m.Title, &m.Body, &m.Pinned, &created, &updated); err != nil {
		return nil, err
	}
	var err error
	if m.CreatedAt, err = time.Parse(time.RFC3339Nano, created); err != nil {
		return nil, err
	}
	if m.UpdatedAt, err = time.Parse(time.RFC3339Nano, updated); err != nil {
		return nil, err
	}
	return &m, nil
}
