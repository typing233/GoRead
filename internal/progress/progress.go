package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Chapter    int       `json:"chapter"`
	LineOffset int       `json:"line_offset"`
	LastRead   time.Time `json:"last_read"`
}

type Progress struct {
	Entries map[string]Entry `json:"entries"`
}

func dataDir() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, "goread")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "goread")
}

func progressPath() string {
	return filepath.Join(dataDir(), "progress.json")
}

func Load() (*Progress, error) {
	data, err := os.ReadFile(progressPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &Progress{Entries: make(map[string]Entry)}, nil
		}
		return nil, err
	}

	p := &Progress{Entries: make(map[string]Entry)}
	if err := json.Unmarshal(data, p); err != nil {
		return &Progress{Entries: make(map[string]Entry)}, nil
	}
	return p, nil
}

func (p *Progress) Save() error {
	dir := dataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(progressPath(), data, 0644)
}

func (p *Progress) Get(epubPath string) (Entry, bool) {
	abs, err := filepath.Abs(epubPath)
	if err != nil {
		abs = epubPath
	}
	e, ok := p.Entries[abs]
	return e, ok
}

func (p *Progress) Set(epubPath string, chapter, lineOffset int) {
	abs, err := filepath.Abs(epubPath)
	if err != nil {
		abs = epubPath
	}
	p.Entries[abs] = Entry{
		Chapter:    chapter,
		LineOffset: lineOffset,
		LastRead:   time.Now(),
	}
}
