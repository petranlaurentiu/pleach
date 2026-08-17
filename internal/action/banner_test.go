package action

import "testing"

func TestComposeBannerWidth(t *testing.T) {
	lines := composeBanner("HI", bannerBig, 1)
	if len(lines) != 5 {
		t.Fatalf("got %d rows", len(lines))
	}
	if len(lines[0]) < 11 {
		t.Fatalf("too narrow: %q", lines[0])
	}
}

func TestRenderSubtitleIsReadable(t *testing.T) {
	const msg = "Welcome back, Laurentiu. Shall we begin?"
	lines := renderSubtitle(msg, 80)
	if len(lines) != 1 || lines[0] != msg {
		t.Fatalf("expected plain text, got %#v", lines)
	}
}
