package buffer

import "testing"

func TestFindAllLiteral(t *testing.T) {
	b := NewBufferFromString("foo bar foo\nfoo", "", BTDefault)
	b.Settings["ignorecase"] = true
	matches, err := b.FindAll("foo", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 3 {
		t.Fatalf("got %d matches, want 3", len(matches))
	}
}

func TestFindAllEmpty(t *testing.T) {
	b := NewBufferFromString("hello", "", BTDefault)
	b.Settings["ignorecase"] = true
	matches, err := b.FindAll("", false)
	if err != nil {
		t.Fatal(err)
	}
	if matches != nil {
		t.Fatalf("empty query should have no matches, got %v", matches)
	}
}