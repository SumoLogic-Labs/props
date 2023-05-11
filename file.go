package props

import (
	"context"
	"fmt"
	"os"

	"github.com/magiconair/properties"
)

type File struct {
	Path string
}

func (f File) Poll(ctx context.Context) (*properties.Properties, error) {
	b, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", f.Path, err)
	}
	l := &properties.Loader{
		DisableExpansion: true,
		Encoding:         properties.UTF8,
	}
	return l.LoadBytes(b)
}

func NewFileSource(filename string) *File {
	return &File{
		Path: filename,
	}
}
