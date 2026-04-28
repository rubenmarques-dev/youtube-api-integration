package domain

import (
	"regexp"
	"strconv"
	"strings"
)

// timestampRe matches a timestamp at the start of a line (with optional leading whitespace),
// capturing the timestamp and the rest of the line as the title.
var timestampRe = regexp.MustCompile(`^\s*(\d{1,2}:\d{2}(?::\d{2})?)\s*(.*)`)

// ParseChapters extracts YouTube chapters from a video description.
// Chapters are lines that start with a timestamp in the format:
//
//	MM:SS, H:MM:SS, or HH:MM:SS (optionally followed by a space and chapter title)
//
// The first chapter must start at 0:00 for YouTube to recognise them as chapters;
// we parse all matching lines regardless.
func ParseChapters(description string) []Chapter {
	chapters := []Chapter{}
	lines := strings.Split(description, "\n")
	for _, line := range lines {
		matches := timestampRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		ts := matches[1]
		title := strings.TrimSpace(matches[2])
		if title == "" {
			continue
		}
		seconds, err := parseTimestampToSeconds(ts)
		if err != nil {
			continue
		}
		chapters = append(chapters, Chapter{
			Title:     title,
			Timestamp: ts,
			Seconds:   seconds,
		})
	}
	return chapters
}

// parseTimestampToSeconds converts a timestamp string (MM:SS or H:MM:SS or HH:MM:SS)
// to total seconds.
func parseTimestampToSeconds(ts string) (int, error) {
	parts := strings.Split(ts, ":")
	var h, m, s int
	var err error
	switch len(parts) {
	case 2:
		m, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		s, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
	case 3:
		h, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		m, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		s, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}
	}
	return h*3600 + m*60 + s, nil
}
