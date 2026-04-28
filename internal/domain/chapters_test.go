package domain_test

import (
	"testing"

	"github.com/rubenmarques/youtube-api-integration/internal/domain"
)

func TestParseChapters_EmptyDescription(t *testing.T) {
	got := domain.ParseChapters("")
	if got == nil {
		t.Fatal("expected non-nil slice, got nil")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestParseChapters_NoTimestamps(t *testing.T) {
	desc := "This is a video description.\nNo timestamps here.\nJust plain text."
	got := domain.ParseChapters(desc)
	if got == nil {
		t.Fatal("expected non-nil slice, got nil")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestParseChapters_MMSS(t *testing.T) {
	desc := "0:00 Intro\n1:30 Main Content\n5:45 Outro"
	got := domain.ParseChapters(desc)
	if len(got) != 3 {
		t.Fatalf("expected 3 chapters, got %d: %v", len(got), got)
	}
	want := []domain.Chapter{
		{Title: "Intro", Timestamp: "0:00", Seconds: 0},
		{Title: "Main Content", Timestamp: "1:30", Seconds: 90},
		{Title: "Outro", Timestamp: "5:45", Seconds: 345},
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("chapter[%d]: got %+v, want %+v", i, got[i], w)
		}
	}
}

func TestParseChapters_HMMSS(t *testing.T) {
	desc := "0:00 Start\n1:00:00 Part 2\n2:30:15 Part 3"
	got := domain.ParseChapters(desc)
	if len(got) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(got))
	}
	if got[1].Seconds != 3600 {
		t.Errorf("expected 3600 seconds for 1:00:00, got %d", got[1].Seconds)
	}
	if got[2].Seconds != 9015 {
		t.Errorf("expected 9015 seconds for 2:30:15, got %d", got[2].Seconds)
	}
}

func TestParseChapters_MixedFormats(t *testing.T) {
	desc := "0:00 Intro\n45:30 Middle\n1:10:00 End"
	got := domain.ParseChapters(desc)
	if len(got) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(got))
	}
	if got[0].Seconds != 0 {
		t.Errorf("expected 0 seconds, got %d", got[0].Seconds)
	}
	if got[1].Seconds != 2730 {
		t.Errorf("expected 2730 seconds for 45:30, got %d", got[1].Seconds)
	}
	if got[2].Seconds != 4200 {
		t.Errorf("expected 4200 seconds for 1:10:00, got %d", got[2].Seconds)
	}
}

func TestParseChapters_TimestampWithNoTitle_Skipped(t *testing.T) {
	desc := "0:00\n1:30 Valid Chapter\n2:00"
	got := domain.ParseChapters(desc)
	if len(got) != 1 {
		t.Fatalf("expected 1 chapter, got %d: %v", len(got), got)
	}
	if got[0].Title != "Valid Chapter" {
		t.Errorf("expected title 'Valid Chapter', got %q", got[0].Title)
	}
}

func TestParseChapters_TimestampMidLine_NotMatched(t *testing.T) {
	// Timestamps embedded mid-line should NOT match (timestamp must be at line start)
	desc := "Watch at 1:30 for the best part\nCheck 5:00 out\n0:00 Actual Intro"
	got := domain.ParseChapters(desc)
	// Only "0:00 Actual Intro" should match since the others have text before the timestamp
	if len(got) != 1 {
		t.Fatalf("expected 1 chapter, got %d: %v", len(got), got)
	}
	if got[0].Title != "Actual Intro" {
		t.Errorf("expected title 'Actual Intro', got %q", got[0].Title)
	}
}

func TestParseChapters_CorrectSecondsCalculation(t *testing.T) {
	tests := []struct {
		ts      string
		seconds int
	}{
		{"0:00 Test", 0},
		{"0:30 Test", 30},
		{"1:00 Test", 60},
		{"1:30 Test", 90},
		{"10:00 Test", 600},
		{"1:00:00 Test", 3600},
		{"1:01:01 Test", 3661},
		{"2:30:00 Test", 9000},
	}
	for _, tt := range tests {
		got := domain.ParseChapters(tt.ts)
		if len(got) != 1 {
			t.Errorf("input %q: expected 1 chapter, got %d", tt.ts, len(got))
			continue
		}
		if got[0].Seconds != tt.seconds {
			t.Errorf("input %q: expected %d seconds, got %d", tt.ts, tt.seconds, got[0].Seconds)
		}
	}
}
