/*
Copyright 2026 The llm-d Authors

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

// Package mock provides mock implementations for testing.
package mock

import (
	"context"
	"io"
	"time"

	"github.com/llm-d-incubation/batch-gateway/internal/files_store/api"
)

// MockBatchFilesClient is a mock implementation of the BatchFilesClient interface.
type MockBatchFilesClient struct{}

// NewMockBatchFilesClient creates a new mock client.
func NewMockBatchFilesClient() *MockBatchFilesClient {
	return &MockBatchFilesClient{}
}

// Store stores a file in the files storage.
func (m *MockBatchFilesClient) Store(ctx context.Context, location string, fileSizeLimit int64, reader io.Reader) (*api.BatchFileMetadata, error) {
	return &api.BatchFileMetadata{
		Location: location,
		Size:     0,
		ModTime:  time.Now(),
	}, nil
}

// Retrieve retrieves a file from the files storage.
func (m *MockBatchFilesClient) Retrieve(ctx context.Context, location string) (io.Reader, *api.BatchFileMetadata, error) {
	return nil, &api.BatchFileMetadata{
		Location: location,
		Size:     0,
		ModTime:  time.Now(),
	}, nil
}

// List lists the files in the specified location.
func (m *MockBatchFilesClient) List(ctx context.Context, location string) ([]api.BatchFileMetadata, error) {
	return []api.BatchFileMetadata{}, nil
}

// Delete deletes the file in the specified location.
func (m *MockBatchFilesClient) Delete(ctx context.Context, location string) error {
	return nil
}

// GetContext returns a derived context for a call.
func (m *MockBatchFilesClient) GetContext(parentCtx context.Context, timeLimit time.Duration) (context.Context, context.CancelFunc) {
	if timeLimit > 0 {
		return context.WithTimeout(parentCtx, timeLimit)
	}
	return context.WithCancel(parentCtx)
}

// Close closes the client.
func (m *MockBatchFilesClient) Close() error {
	return nil
}
