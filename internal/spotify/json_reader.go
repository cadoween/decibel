package spotify

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/cadoween/decibel/pkg/ext"
	"github.com/rs/zerolog"
)

type JSONReader struct {
	// bufferPool helps reduce memory allocations when reading files.
	bufferPool sync.Pool
}

func NewJSONReader() *JSONReader {
	return &JSONReader{
		bufferPool: sync.Pool{
			New: func() any {
				return bufio.NewReaderSize(nil, 32*1024) // 32KB buffer
			},
		},
	}
}

func (r *JSONReader) ReadStreamsFromFolder(ctx context.Context, folderPath string) ([]Stream, error) {
	logger := zerolog.Ctx(ctx)
	streamsChan := make(chan []Stream)
	errChan := make(chan error)

	var (
		wg        sync.WaitGroup
		jsonFiles []string
	)

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ext.JSON {
			jsonFiles = append(jsonFiles, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath.WalkDir: %w", err)
	}
	logger.Debug().Int("files_found", len(jsonFiles)).Msg("Found JSON files to process")

	for _, path := range jsonFiles {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:

				streams, err := r.readStreamsFromFile(ctx, filePath)
				if err != nil {
					select {
					case errChan <- fmt.Errorf("r.readStreamsFromFile %s: %w", filePath, err):
					case <-ctx.Done():
					}
					return

				}

				select {
				case streamsChan <- streams:
					logger.Debug().
						Str("file", filepath.Base(filePath)).
						Int("streams_count", len(streams)).
						Msg("Processed file")
				case <-ctx.Done():
				}
			}
		}(path)
	}

	go func() {
		wg.Wait()
		close(streamsChan)
		close(errChan)
	}()

	var allStreams []Stream
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err, ok := <-errChan:
			if ok && err != nil {
				return nil, fmt.Errorf("reading streams: %w", err)
			}
		case streams, ok := <-streamsChan:
			if !ok {
				logger.Info().Int("total_streams", len(allStreams)).Msg("Completed reading all streams")
				return allStreams, nil
			}
			allStreams = append(allStreams, streams...)

		}
	}
}

func (r *JSONReader) readStreamsFromFile(ctx context.Context, path string) ([]Stream, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Str("file", filepath.Base(path)).Msg("Reading file")

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer func() { _ = file.Close() }()

	reader := r.bufferPool.Get().(*bufio.Reader)
	reader.Reset(file)
	defer r.bufferPool.Put(reader)

	streams := make([]Stream, 0, 20000) // max estimation of streams per file

	decoder := json.NewDecoder(reader)
	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("decoder.Token: %w", err)
	}

	for decoder.More() {
		var stream Stream
		if err := decoder.Decode(&stream); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("decoder.Decode: %w", err)
		}
		streams = append(streams, stream)
	}

	logger.Debug().
		Str("file", filepath.Base(path)).
		Int("streams_count", len(streams)).
		Msg("Successfully read streams from file")

	return streams, nil
}
