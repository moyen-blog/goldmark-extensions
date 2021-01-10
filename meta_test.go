package meta

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

func TestMeta(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Meta,
		),
	)
	source := `---
title: goldmark-meta
ignored:
  a: 2
  b: [3, 4]
---
# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct {
		Title string
	}{}
	if err := Unmarshal(context, &out); err != nil {
		t.Error("YAML unmarshal failed")
	}
	if out.Title != "goldmark-meta" {
		t.Errorf("Title must be %s, but got %v", "goldmark-meta", out.Title)
	}
	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("should render '<h1>Hello goldmark-meta</h1>', but '%s'", buf.String())
	}
}

func TestNoMeta(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Meta,
		),
	)
	source := `# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct {
		Title string
	}{}
	if err := Unmarshal(context, &out); err != nil {
		t.Error("YAML unmarshal failed")
	}
	if out.Title != "" {
		t.Errorf("Title must be empty, but got %v", out.Title)
	}
	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("should render '<h1>Hello goldmark-meta</h1>', but '%s'", buf.String())
	}
}

func TestEmptyMeta(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Meta,
		),
	)
	source := `---
---
# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct {
		Title string
	}{}
	if err := Unmarshal(context, &out); err != nil {
		t.Error("YAML unmarshal failed")
	}
	if out.Title != "" {
		t.Errorf("Title must be empty, but got %v", out.Title)
	}
	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("should render '<h1>Hello goldmark-meta</h1>', but '%s'", buf.String())
	}
}

func TestMetaError(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			New(),
		),
	)
	source := `---
bad:
  - : {
  }
    - markdown
    - goldmark
---
# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct{}{}
	if err := Unmarshal(context, &out); err == nil {
		t.Error("Should throw unmarshal error")
	}
}
