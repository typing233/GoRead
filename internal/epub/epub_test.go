package epub

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", "test.epub")
}

func TestOpenEpub(t *testing.T) {
	book, err := Open(testdataPath())
	if err != nil {
		t.Fatalf("failed to open epub: %v", err)
	}
	defer book.Close()

	if book.Title != "Test Book" {
		t.Errorf("expected title 'Test Book', got '%s'", book.Title)
	}

	if book.Author != "Test Author" {
		t.Errorf("expected author 'Test Author', got '%s'", book.Author)
	}

	if len(book.Chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(book.Chapters))
	}

	data, err := book.LoadChapter(0)
	if err != nil {
		t.Fatalf("failed to load chapter 0: %v", err)
	}
	if len(data) == 0 {
		t.Error("chapter 0 data is empty")
	}
}
