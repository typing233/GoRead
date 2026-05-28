package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/goread/goread/internal/config"
	"github.com/goread/goread/internal/epub"
	"github.com/goread/goread/internal/htmlconv"
	"github.com/goread/goread/internal/paginator"
	"github.com/goread/goread/internal/progress"
	"github.com/goread/goread/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: goread <file.epub>\n")
		os.Exit(1)
	}

	epubPath := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
		cfg = config.Default()
	}

	book, err := epub.Open(epubPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening epub: %v\n", err)
		os.Exit(1)
	}
	defer book.Close()

	if len(book.Chapters) == 0 {
		fmt.Fprintf(os.Stderr, "Error: epub has no chapters\n")
		os.Exit(1)
	}

	prog, err := progress.Load()
	if err != nil {
		prog = &progress.Progress{Entries: make(map[string]progress.Entry)}
	}

	chapterIdx := 0
	lineOffset := 0
	if entry, ok := prog.Get(epubPath); ok {
		chapterIdx = entry.Chapter
		lineOffset = entry.LineOffset
		if chapterIdx >= len(book.Chapters) {
			chapterIdx = 0
			lineOffset = 0
		}
	}

	screen, err := ui.InitScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		prog.Set(epubPath, chapterIdx, lineOffset)
		prog.Save()
		screen.Fini()
		os.Exit(0)
	}()

	renderer := ui.NewRenderer(screen, cfg)
	inputHandler := ui.NewInputHandler()

	pag, err := loadChapter(book, chapterIdx, screen, cfg)
	if err != nil {
		screen.Fini()
		fmt.Fprintf(os.Stderr, "Error loading chapter: %v\n", err)
		os.Exit(1)
	}

	if lineOffset >= pag.TotalLines() {
		lineOffset = 0
	}

	render := func() {
		_, height := screen.Size()
		pageHeight := height - 1
		lines := pag.GetLines(lineOffset, pageHeight)
		renderer.RenderPage(lines, lineOffset, pag.TotalLines(),
			book.Chapters[chapterIdx].Title, chapterIdx, len(book.Chapters))
	}

	render()

	for {
		ev := screen.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			w, h := e.Size()
			screen.Sync()
			pag.Reflow(w-2*cfg.SideMargin, h-1)
			if lineOffset >= pag.TotalLines() {
				lineOffset = pag.TotalLines() - 1
				if lineOffset < 0 {
					lineOffset = 0
				}
			}
			render()

		case *tcell.EventKey:
			action := inputHandler.Handle(e)
			_, height := screen.Size()
			pageHeight := height - 1
			halfPage := pageHeight / 2

			switch action {
			case ui.ActionQuit:
				prog.Set(epubPath, chapterIdx, lineOffset)
				prog.Save()
				return

			case ui.ActionScrollDown:
				if lineOffset < pag.TotalLines()-1 {
					lineOffset++
					render()
				}

			case ui.ActionScrollUp:
				if lineOffset > 0 {
					lineOffset--
					render()
				}

			case ui.ActionHalfPageDown:
				lineOffset += halfPage
				if lineOffset >= pag.TotalLines() {
					lineOffset = pag.TotalLines() - 1
				}
				if lineOffset < 0 {
					lineOffset = 0
				}
				render()

			case ui.ActionHalfPageUp:
				lineOffset -= halfPage
				if lineOffset < 0 {
					lineOffset = 0
				}
				render()

			case ui.ActionPageDown:
				lineOffset += pageHeight
				if lineOffset >= pag.TotalLines() {
					lineOffset = pag.TotalLines() - 1
				}
				if lineOffset < 0 {
					lineOffset = 0
				}
				render()

			case ui.ActionPageUp:
				lineOffset -= pageHeight
				if lineOffset < 0 {
					lineOffset = 0
				}
				render()

			case ui.ActionTop:
				lineOffset = 0
				render()

			case ui.ActionBottom:
				lineOffset = pag.TotalLines() - pageHeight
				if lineOffset < 0 {
					lineOffset = 0
				}
				render()

			case ui.ActionNextChapter:
				if chapterIdx < len(book.Chapters)-1 {
					chapterIdx++
					lineOffset = 0
					newPag, err := loadChapter(book, chapterIdx, screen, cfg)
					if err == nil {
						pag = newPag
					}
					render()
				}

			case ui.ActionPrevChapter:
				if chapterIdx > 0 {
					chapterIdx--
					lineOffset = 0
					newPag, err := loadChapter(book, chapterIdx, screen, cfg)
					if err == nil {
						pag = newPag
					}
					render()
				}
			}
		}
	}
}

func loadChapter(book *epub.Book, idx int, screen tcell.Screen, cfg *config.Config) (*paginator.Paginator, error) {
	data, err := book.LoadChapter(idx)
	if err != nil {
		return nil, err
	}

	lines, err := htmlconv.ConvertHTML(data)
	if err != nil {
		return nil, err
	}

	w, h := screen.Size()
	contentWidth := w - 2*cfg.SideMargin
	if contentWidth < 10 {
		contentWidth = w
	}
	pageHeight := h - 1

	return paginator.New(lines, contentWidth, pageHeight), nil
}
