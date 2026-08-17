package action

import "testing"

func TestFindPromptLabel(t *testing.T) {
	if got := findPromptLabel(0, 0, false); got != "Find 0: " {
		t.Fatalf("got %q", got)
	}
	if got := findPromptLabel(2, 5, false); got != "Find 2/5: " {
		t.Fatalf("got %q", got)
	}
	if got := findPromptLabel(1, 3, true); got != "Find (regex) 1/3: " {
		t.Fatalf("got %q", got)
	}
}

func TestIsFindPrompt(t *testing.T) {
	if !isFindPrompt("Find") || !isFindPrompt("FindRegex") {
		t.Fatal("expected find prompt types")
	}
	if isFindPrompt("Command") {
		t.Fatal("command is not a find prompt")
	}
}
