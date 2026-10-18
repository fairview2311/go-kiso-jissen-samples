package web

import (
	"encoding/json/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"memos/internal/memo"
	"memos/internal/store"
)

func TestAPI(t *testing.T) {
	ts := httptest.NewServer(NewServer(store.NewMemoryStore()).Handler())
	defer ts.Close()

	// 足す
	res := do(t, http.MethodPost, ts.URL+"/api/memos", `{"title": "買い物リスト", "body": "牛乳と卵を買う"}`)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("POST のステータス = %d、期待したのは 201", res.StatusCode)
	}
	var created memo.Memo
	if err := json.UnmarshalRead(res.Body, &created); err != nil {
		t.Fatal(err)
	}

	// 1件を取る
	res = do(t, http.MethodGet, ts.URL+"/api/memos/"+created.ID, "")
	if res.StatusCode != http.StatusOK {
		t.Errorf("GET のステータス = %d、期待したのは 200", res.StatusCode)
	}

	// 題が空なら 400
	res = do(t, http.MethodPost, ts.URL+"/api/memos", `{"title": ""}`)
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("空の題の POST のステータス = %d、期待したのは 400", res.StatusCode)
	}

	// 消すと 204、もう一度取ると 404
	res = do(t, http.MethodDelete, ts.URL+"/api/memos/"+created.ID, "")
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE のステータス = %d、期待したのは 204", res.StatusCode)
	}
	res = do(t, http.MethodGet, ts.URL+"/api/memos/"+created.ID, "")
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("消した後の GET のステータス = %d、期待したのは 404", res.StatusCode)
	}
}

func TestListPage(t *testing.T) {
	s := store.NewMemoryStore()
	m, err := memo.New("会議のメモ", "来週の会議で予算を決める")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Save(m); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(NewServer(s).Handler())
	defer ts.Close()

	res := do(t, http.MethodGet, ts.URL+"/", "")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), "会議のメモ") {
		t.Errorf("一覧の画面 = %d %q", res.StatusCode, body)
	}
}

// do は、method と body でリクエストを送り、レスポンスを返す。本文はテストの終わりに閉じる。
func do(t *testing.T, method, url, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { res.Body.Close() })
	return res
}
