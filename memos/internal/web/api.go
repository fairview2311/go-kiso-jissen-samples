// Package web は、memos の Web API と画面を作る。
package web

import (
	"encoding/json/v2"
	"errors"
	"log"
	"net/http"
	"time"

	"memos/internal/memo"
	"memos/internal/store"
)

// Server は、メモの保存先を持ち、HTTP のリクエストに答える。
type Server struct {
	store store.Store
}

// NewServer は、s を保存先にする Server を作る。
func NewServer(s store.Store) *Server {
	return &Server{store: s}
}

// Handler は、URL とハンドラの対応を登録した http.Handler を返す。
func (srv *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/memos", srv.listMemos)
	mux.HandleFunc("POST /api/memos", srv.createMemo)
	mux.HandleFunc("GET /api/memos/{id}", srv.getMemo)
	mux.HandleFunc("PUT /api/memos/{id}", srv.updateMemo)
	mux.HandleFunc("DELETE /api/memos/{id}", srv.deleteMemo)
	srv.registerPages(mux)
	return logRequests(mux)
}

// memoInput は、追加と更新のリクエストの JSON の形。
type memoInput struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Pinned bool   `json:"pinned"`
}

func (srv *Server) listMemos(w http.ResponseWriter, r *http.Request) {
	memos, err := srv.store.List()
	if err != nil {
		writeError(w, err)
		return
	}
	word := r.URL.Query().Get("q")
	found := []*memo.Memo{}
	for _, m := range memos {
		if word == "" || m.Contains(word) {
			found = append(found, m)
		}
	}
	writeJSON(w, http.StatusOK, found)
}

func (srv *Server) createMemo(w http.ResponseWriter, r *http.Request) {
	var in memoInput
	if err := json.UnmarshalRead(r.Body, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "JSON を読めません: " + err.Error()})
		return
	}
	m, err := memo.New(in.Title, in.Body)
	if err != nil {
		writeError(w, err)
		return
	}
	m.Pinned = in.Pinned
	if err := srv.store.Save(m); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (srv *Server) getMemo(w http.ResponseWriter, r *http.Request) {
	m, err := srv.store.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (srv *Server) updateMemo(w http.ResponseWriter, r *http.Request) {
	m, err := srv.store.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	var in memoInput
	if err := json.UnmarshalRead(r.Body, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "JSON を読めません: " + err.Error()})
		return
	}
	if err := m.Edit(in.Title, in.Body); err != nil {
		writeError(w, err)
		return
	}
	m.Pinned = in.Pinned
	if err := srv.store.Save(m); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (srv *Server) deleteMemo(w http.ResponseWriter, r *http.Request) {
	if err := srv.store.Delete(r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// errorBody は、エラーのときに返す JSON の形。
type errorBody struct {
	Error string `json:"error"`
}

// writeJSON は、v を JSON にして、ステータスコード status で返す。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.MarshalWrite(w, v); err != nil {
		log.Printf("JSON を書けません: %v", err)
	}
}

// writeError は、エラーの種類に合ったステータスコードで、エラーの JSON を返す。
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, store.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, memo.ErrEmptyTitle), errors.Is(err, memo.ErrTitleTooLong):
		status = http.StatusBadRequest
	}
	writeJSON(w, status, errorBody{Error: err.Error()})
}

// logRequests は、リクエストごとに、メソッド・パス・かかった時間を記録するミドルウェア。
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start).Round(time.Microsecond))
	})
}
