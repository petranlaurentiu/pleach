package action

import (
	"strings"
	"testing"
	"time"
)

func TestWelcomeGreeting(t *testing.T) {
	cases := map[int]string{
		8:  "Good morning",
		14: "Good afternoon",
		21: "Good evening",
	}
	for hour, want := range cases {
		if got := welcomeGreeting(hour); got != want {
			t.Fatalf("hour %d: got %q, want %q", hour, got, want)
		}
	}
}

func TestWelcomeMessageIncludesName(t *testing.T) {
	now := time.Date(2026, 8, 17, 21, 0, 0, 0, time.UTC)
	got := welcomeMessage(now, "Laurentiu", 0)
	if !strings.Contains(got, "Laurentiu") && !strings.Contains(got, "Good evening") {
		t.Fatalf("unexpected welcome: %q", got)
	}
}

func TestWelcomeMessageRotatesSlowly(t *testing.T) {
	now := time.Date(2026, 8, 17, 21, 0, 0, 0, time.UTC)
	first := welcomeMessage(now, "Laurentiu", time.Second)
	same := welcomeMessage(now, "Laurentiu", 5*time.Second)
	next := welcomeMessage(now, "Laurentiu", 6*time.Second)
	if first != same {
		t.Fatalf("message changed too soon: %q vs %q", first, same)
	}
	if next == first {
		t.Fatalf("message did not change after %s", welcomeRotateEvery)
	}
}

func TestFirstName(t *testing.T) {
	if got := firstName("Laurentiu Petran"); got != "Laurentiu" {
		t.Fatalf("got %q", got)
	}
}
