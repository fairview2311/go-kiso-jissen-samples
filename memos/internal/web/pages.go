package web

import (
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"memos/internal/memo"
	"memos/internal/store"
)

// templates と static の中のファイルを、実行ファイルに埋め込む。
//
//go:embed templates/*.html static/*
var files embed.FS

var pages = template.Must(template.New("").Funcs(template.FuncMap{
	"date": func(m *memo.Memo) string { return m.UpdatedAt.Format("2006-01-02 15:04") },
}).ParseFS(files, "templates/*.html"))

// listPage は、一覧の画面に渡すデータ。
type listPage struct {
	Query string
	Memos []*memo.Memo
}

// formPage は、作成と編集の画面に渡すデータ。Memo が nil なら作成。
type formPage struct {
	Memo   *memo.Memo
	Title  string
	Body   string
	Pinned bool
	Error  string
}

// registerPages は、画面のハンドラを mux に登録する。
func (srv *Server) registerPages(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", srv.listPage)
	mux.HandleFunc("GET /memos/new", srv.newPage)
	mux.HandleFunc("POST /memos", srv.createPage)
	mux.HandleFunc("GET /memos/{id}/edit", srv.editPage)
	mux.HandleFunc("POST /memos/{id}", srv.updatePage)
	mux.HandleFunc("POST /memos/{id}/delete", srv.deletePage)
	mux.Handle("GET /static/", http.FileServerFS(files))
}

func (srv *Server) listPage(w http.ResponseWriter, r *http.Request) {
	memos, err := srv.store.List()
	if err != nil {
		http.Error(w, "メモを読めません", http.StatusInternalServerError)
		return
	}
	q := r.URL.Query().Get("q")
	var found []*memo.Memo
	for _, m := range memos {
		if q == "" || m.Contains(q) {
			found = append(found, m)
		}
	}
	render(w, http.StatusOK, "list.html", listPage{Query: q, Memos: found})
}

func (srv *Server) newPage(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "form.html", formPage{})
}

func (srv *Server) createPage(w http.ResponseWriter, r *http.Request) {
	title, body, pinned := r.FormValue("title"), r.FormValue("body"), r.FormValue("pinned") == "on"
	m, err := memo.New(title, body)
	if err != nil {
		render(w, http.StatusBadRequest, "form.html", formPage{Title: title, Body: body, Pinned: pinned, Error: err.Error()})
		return
	}
	m.Pinned = pinned
	if err := srv.store.Save(m); err != nil {
		http.Error(w, "保存できません", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (srv *Server) editPage(w http.ResponseWriter, r *http.Request) {
	m, ok := srv.findMemo(w, r)
	if !ok {
		return
	}
	render(w, http.StatusOK, "form.html", formPage{Memo: m, Title: m.Title, Body: m.Body, Pinned: m.Pinned})
}

func (srv *Server) updatePage(w http.ResponseWriter, r *http.Request) {
	m, ok := srv.findMemo(w, r)
	if !ok {
		return
	}
	title, body, pinned := r.FormValue("title"), r.FormValue("body"), r.FormValue("pinned") == "on"
	if err := m.Edit(title, body); err != nil {
		render(w, http.StatusBadRequest, "form.html", formPage{Memo: m, Title: title, Body: body, Pinned: pinned, Error: err.Error()})
		return
	}
	m.Pinned = pinned
	if err := srv.store.Save(m); err != nil {
		http.Error(w, "保存できません", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (srv *Server) deletePage(w http.ResponseWriter, r *http.Request) {
	if err := srv.store.Delete(r.PathValue("id")); err != nil && !errors.Is(err, store.ErrNotFound) {
		http.Error(w, "消せません", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// findMemo は、パスの {id} のメモを探す。無ければ 404 の画面を返して false を返す。
func (srv *Server) findMemo(w http.ResponseWriter, r *http.Request) (*memo.Memo, bool) {
	m, err := srv.store.Get(r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "メモが見つかりません", http.StatusNotFound)
		return nil, false
	}
	if err != nil {
		http.Error(w, "メモを読めません", http.StatusInternalServerError)
		return nil, false
	}
	return m, true
}

// render は、テンプレート name に data をはめ込んで、status で返す。
func render(w http.ResponseWriter, status int, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := pages.ExecuteTemplate(w, name, data); err != nil {
		slog.Error("画面を作れません", "err", err)
	}
}
