/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nuclio

import (
	"io"
	"sync"
)

// Response can be returned from functions, allowing the user to specify various fields
type Response struct {
	StatusCode  int
	ContentType string
	Headers     map[string]interface{}
	Body        []byte
}

func (r *Response) IsStream() bool {
	return false
}

type ResponseStream struct {
	StatusCode  *int
	ContentType string
	Headers     map[string]interface{}
	Body        io.ReadCloser

	writer io.Writer
	mu     sync.Mutex
}

// NewResponseStream creates a new ResponseStream backed by io.Pipe.
func NewResponseStream() *ResponseStream {
	reader, writer := io.Pipe()
	return &ResponseStream{
		Body:   reader,
		writer: writer,
		mu:     sync.Mutex{},
	}
}

// NewCustomResponseStream allows creating a ResponseStream with custom reader and writer.
func NewCustomResponseStream(reader io.ReadCloser, writer io.Writer) *ResponseStream {
	return &ResponseStream{
		Body:   reader,
		writer: writer,
		mu:     sync.Mutex{},
	}
}

func (s *ResponseStream) GetWriter() io.Writer {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.writer
}

// StreamFrom asynchronously copies data from the provided reader.
func (s *ResponseStream) StreamFrom(reader io.Reader) (int64, error) {
	writer := s.GetWriter()

	if writer == nil {
		return 0, io.ErrClosedPipe
	}

	return io.Copy(writer, reader)
}

// SendChunk writes a chunk of data to the response stream.
func (s *ResponseStream) SendChunk(chunk []byte) (int, error) {
	writer := s.GetWriter()

	if writer == nil {
		return 0, io.ErrClosedPipe
	}

	return s.writer.Write(chunk)
}

// StopStreaming finalizes the response by closing the writer and setting the status code.
func (s *ResponseStream) StopStreaming(statusCode int) {
	// set status code before closing the writer
	// this is important because code will be set right after the end of streaming
	s.SetStatusCode(statusCode)

	s.CloseWriter()
}

// SetStatusCode sets the status code for the response stream
func (s *ResponseStream) SetStatusCode(statusCode int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.StatusCode = &statusCode
}

// CloseWriter closes the writer
func (s *ResponseStream) CloseWriter() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pipeWriter, ok := s.writer.(io.Closer); ok {
		_ = pipeWriter.Close()
	}
	s.writer = nil
}

func (s *ResponseStream) IsStream() bool {
	return true
}

type ProcessingResult interface {
	// IsStream checks if the result is a stream
	IsStream() bool
}
