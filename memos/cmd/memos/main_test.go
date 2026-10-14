package main

import (
	"bytes"
	"strings"
	"testing"

	"memos/internal/store"
)

func TestRun(t *testing.T) {
	s := store.NewMemoryStore()
	var out bytes.Buffer

	if err := run(&out, s, []string{"add", "読みたい本", "Go の本を2冊"}); err != nil {
		t.Fatalf("add のエラー = %v", err)
	}
	if !strings.Contains(out.String(), "追加しました:") {
		t.Errorf("add の出力 = %q", out.String())
	}

	out.Reset()
	if err := run(&out, s, []string{"list", "本"}); err != nil {
		t.Fatalf("list のエラー = %v", err)
	}
	if !strings.Contains(out.String(), "1 件のメモ（全部で 1 件）") {
		t.Errorf("list の出力 = %q", out.String())
	}

	if err := run(&out, s, []string{"remove"}); err == nil {
		t.Error("知らないサブコマンドなのにエラーにならない")
	}
}
