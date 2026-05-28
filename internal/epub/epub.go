package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"
)

type Book struct {
	Title    string
	Author   string
	Chapters []Chapter
	reader   *zip.ReadCloser
}

type Chapter struct {
	ID    string
	Title string
	Path  string
}

type container struct {
	Rootfiles []rootfile `xml:"rootfiles>rootfile"`
}

type rootfile struct {
	FullPath string `xml:"full-path,attr"`
}

func Open(filepath string) (*Book, error) {
	r, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("opening epub: %w", err)
	}

	book := &Book{reader: r}

	containerData, err := book.readFile("META-INF/container.xml")
	if err != nil {
		r.Close()
		return nil, fmt.Errorf("reading container.xml: %w", err)
	}

	var c container
	if err := xml.Unmarshal(containerData, &c); err != nil {
		r.Close()
		return nil, fmt.Errorf("parsing container.xml: %w", err)
	}

	if len(c.Rootfiles) == 0 {
		r.Close()
		return nil, fmt.Errorf("no rootfile found in container.xml")
	}

	opfPath := c.Rootfiles[0].FullPath
	opfData, err := book.readFile(opfPath)
	if err != nil {
		r.Close()
		return nil, fmt.Errorf("reading OPF: %w", err)
	}

	opfDir := path.Dir(opfPath)
	if err := book.parseOPF(opfData, opfDir); err != nil {
		r.Close()
		return nil, fmt.Errorf("parsing OPF: %w", err)
	}

	return book, nil
}

func (b *Book) Close() error {
	if b.reader != nil {
		return b.reader.Close()
	}
	return nil
}

func (b *Book) LoadChapter(index int) ([]byte, error) {
	if index < 0 || index >= len(b.Chapters) {
		return nil, fmt.Errorf("chapter index %d out of range", index)
	}
	return b.readFile(b.Chapters[index].Path)
}

func (b *Book) readFile(name string) ([]byte, error) {
	name = strings.TrimPrefix(name, "/")
	for _, f := range b.reader.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("file not found in epub: %s", name)
}

func (b *Book) parseOPF(data []byte, baseDir string) error {
	var pkg opfPackage
	if err := xml.Unmarshal(data, &pkg); err != nil {
		return err
	}

	b.Title = pkg.Metadata.Title
	if len(pkg.Metadata.Creator) > 0 {
		b.Author = pkg.Metadata.Creator[0]
	}

	manifest := make(map[string]manifestItem)
	for _, item := range pkg.Manifest.Items {
		manifest[item.ID] = item
	}

	for _, ref := range pkg.Spine.ItemRefs {
		item, ok := manifest[ref.IDRef]
		if !ok {
			continue
		}
		mediaType := strings.ToLower(item.MediaType)
		if !strings.Contains(mediaType, "html") && !strings.Contains(mediaType, "xml") {
			continue
		}
		chapterPath := item.Href
		if baseDir != "" && baseDir != "." {
			chapterPath = baseDir + "/" + item.Href
		}
		b.Chapters = append(b.Chapters, Chapter{
			ID:    item.ID,
			Title: item.ID,
			Path:  chapterPath,
		})
	}

	return nil
}
