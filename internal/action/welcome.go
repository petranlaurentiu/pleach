package action

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"
)

func welcomeName() string {
	if u, err := user.Current(); err == nil {
		if n := firstName(u.Name); n != "" {
			return n
		}
		if u.Username != "" {
			return u.Username
		}
	}
	if n := os.Getenv("USER"); n != "" {
		return n
	}
	return "sir"
}

func firstName(full string) string {
	for _, part := range strings.Fields(full) {
		if part != "" {
			return part
		}
	}
	return ""
}

func welcomeGreeting(hour int) string {
	switch {
	case hour < 12:
		return "Good morning"
	case hour < 18:
		return "Good afternoon"
	default:
		return "Good evening"
	}
}

const welcomeRotateEvery = 6 * time.Second

func welcomeLines(now time.Time, name string) []string {
	if name == "" {
		name = "sir"
	}
	greeting := welcomeGreeting(now.Hour())
	return []string{
		fmt.Sprintf("%s, %s. All systems online.", greeting, name),
		fmt.Sprintf("%s. Pleach is at your service.", greeting),
		fmt.Sprintf("Systems check complete. Ready when you are, %s.", name),
		fmt.Sprintf("Welcome back, %s. Shall we begin?", name),
	}
}

func welcomeMessage(now time.Time, name string, elapsed time.Duration) string {
	lines := welcomeLines(now, name)
	idx := int(elapsed/welcomeRotateEvery) % len(lines)
	if idx < 0 {
		idx = 0
	}
	return lines[idx]
}
