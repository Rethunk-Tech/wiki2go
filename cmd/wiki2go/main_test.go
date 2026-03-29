package main

import (
	"testing"

	"go.abhg.dev/goldmark/wikilink"
)

func TestMakeNameCanonical(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"Hello World", "hello_world"},
		{"already_lower", "already_lower"},
		{"Foo_Bar", "foo_bar"},
		{"  a  b  ", "__a__b__"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			if got := makeNameCanonical(tt.in); got != tt.want {
				t.Errorf("makeNameCanonical(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestMakeNamePretty(t *testing.T) {
	t.Parallel()
	got := makeNamePretty("hello_world")
	want := "Hello World"
	if got != want {
		t.Errorf("makeNamePretty(%q) = %q, want %q", "hello_world", got, want)
	}
}

func TestResolveWikilink(t *testing.T) {
	t.Parallel()
	var res r
	tests := []struct {
		name     string
		target   string
		fragment string
		want     string
	}{
		{"target only", "Foo Bar", "", "foo_bar"},
		{"fragment only", "", "Section", "#Section"},
		{"target and fragment", "My Page", "intro", "my_page#intro"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := &wikilink.Node{
				Target:   []byte(tt.target),
				Fragment: []byte(tt.fragment),
			}
			got, err := res.ResolveWikilink(n)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("ResolveWikilink(...) = %q, want %q", got, tt.want)
			}
		})
	}
}
