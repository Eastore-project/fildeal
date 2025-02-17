package buffer

import (
	"fmt"
	"path/filepath"
)

// LocalBuffer implements Buffer for local storage
type LocalBuffer struct {
	port *int // Optional port for server URL generation
}

func NewLocalBuffer(port ...int) Buffer {
	var p *int
	if len(port) > 0 {
		p = &port[0]
	}
	return &LocalBuffer{
		port: p,
	}
}

func (b *LocalBuffer) Store(filePath string) (*Response, error) {
	fileName := filepath.Base(filePath)
	return &Response{
		URL:  filePath,
		Hash: fileName,
	}, nil
}

// StoreForServer creates a local URL for server access
func (b *LocalBuffer) StoreForServer(filePath string) (*Response, error) {
	if b.port == nil {
		return b.Store(filePath)
	}

	fileName := filepath.Base(filePath)
	localURL := fmt.Sprintf("http://localhost:%d/download/car?file_name=%s", *b.port, fileName)
	return &Response{
		URL:  localURL,
		Hash: fileName,
	}, nil
}
