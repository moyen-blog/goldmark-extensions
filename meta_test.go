package meta

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

var markdown = goldmark.New(
	goldmark.WithExtensions(
		New(),
	),
)

func TestMeta(t *testing.T) {
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

func TestInvalidMetaBuffer(t *testing.T) {
	source := `# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	context.Set(contextKey, 0) // Not the expected bytes.Buffer
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(parser.NewContext())); err != nil {
		panic(err)
	}
	out := struct {
		Title string
	}{}
	if err := Unmarshal(context, &out); err == nil {
		t.Error("Should throw unavailable YAML buffer error")
	}
}
