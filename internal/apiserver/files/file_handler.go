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

// The file provides HTTP handlers for file-related API endpoints.
// It implements the OpenAI compatible Files API endpoints for uploading, downloading, listing, and deleting files.
package files

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/llm-d-incubation/batch-gateway/internal/apiserver/common"
	"github.com/llm-d-incubation/batch-gateway/internal/files_store/api"
	"github.com/llm-d-incubation/batch-gateway/internal/shared/openai"
	"github.com/llm-d-incubation/batch-gateway/internal/util/logging"
)

type FilesApiHandler struct {
	config      *common.ServerConfig
	filesClient api.BatchFilesClient
}

func NewFilesApiHandler(config *common.ServerConfig, filesClient api.BatchFilesClient) *FilesApiHandler {
	return &FilesApiHandler{
		config:      config,
		filesClient: filesClient,
	}
}

func (c *FilesApiHandler) GetRoutes() []common.Route {
	return []common.Route{
		{
			Method:      http.MethodPost,
			Pattern:     "/v1/files",
			HandlerFunc: c.CreateFile,
		},
		{
			Method:      http.MethodDelete,
			Pattern:     "/v1/files/{file_id}",
			HandlerFunc: c.DeleteFile,
		},
		{
			Method:      http.MethodGet,
			Pattern:     "/v1/files/{file_id}/content",
			HandlerFunc: c.DownloadFile,
		},
		{
			Method:      http.MethodGet,
			Pattern:     "/v1/files",
			HandlerFunc: c.ListFiles,
		},
		{
			Method:      http.MethodGet,
			Pattern:     "/v1/files/{file_id}",
			HandlerFunc: c.RetrieveFile,
		},
	}
}

func (c *FilesApiHandler) CreateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetRequestLogger(r)

	// Check Content-Length header before reading the body
	maxFileSize := c.config.GetMaxFileSizeBytes()
	if r.ContentLength > maxFileSize {
		logger.Info("file size exceeds limit", "contentLength", r.ContentLength, "limit", maxFileSize)
		apiErr := openai.NewAPIError(
			http.StatusBadRequest,
			"",
			fmt.Sprintf("File size exceeds the maximum allowed size of %d bytes", maxFileSize),
			nil,
		)
		common.WriteAPIError(ctx, w, apiErr)
		return
	}

	// Read form file from request
	reader, filename, err := common.ReadFormFile(r, "file")
	if err != nil {
		logger.Error(err, "failed to read form file from request")
		common.WriteInternalServerError(ctx, w)
		return
	}

	purpose := r.FormValue("purpose")

	// Parse expires_after parameters if provided, otherwise use default TTL from config
	var expiresAt int64
	expiresAfterAnchor := r.FormValue("expires_after[anchor]")
	expiresAfterSecondsStr := r.FormValue("expires_after[seconds]")
	if expiresAfterAnchor != "" && expiresAfterSecondsStr != "" {
		expiresAfterSeconds, err := strconv.ParseInt(expiresAfterSecondsStr, 10, 64)
		if err != nil {
			logger.Error(err, "failed to parse expires_after[seconds]")
			common.WriteInternalServerError(ctx, w)
			return
		}

		createdAt := time.Now()
		expiresAt = createdAt.Add(time.Duration(expiresAfterSeconds) * time.Second).Unix()
		logger.Info("file expiration set from request", "anchor", expiresAfterAnchor, "seconds", expiresAfterSeconds, "expiresAt", expiresAt)
	} else if c.config.FileTTLSeconds > 0 {
		// Use default TTL from config if expires_after not provided
		createdAt := time.Now()
		expiresAt = createdAt.Add(time.Duration(c.config.FileTTLSeconds) * time.Second).Unix()
		logger.Info("file expiration set from config default", "ttlSeconds", c.config.FileTTLSeconds, "expiresAt", expiresAt)
	}

	// Create a temporary file under /tmp
	tmpFile, err := os.CreateTemp("/tmp", "upload-*.tmp")
	if err != nil {
		logger.Error(err, "failed to create temp file")
		common.WriteInternalServerError(ctx, w)
		return
	}
	defer tmpFile.Close()

	// Store the file using the files client
	meta, err := c.filesClient.Store(ctx, tmpFile.Name(), maxFileSize, reader)
	if err != nil {
		logger.Error(err, "failed to store file")
		common.WriteInternalServerError(ctx, w)
		return
	}

	logger.Info("file stored successfully", "location", meta.Location, "size", meta.Size)

	// Construct create response
	fileObj := openai.FileObject{
		ID:        meta.ID,
		Bytes:     meta.Size,
		CreatedAt: meta.ModTime.Unix(),
		ExpiresAt: expiresAt,
		Filename:  filename,
		Object:    "file",
		Purpose:   openai.FileObjectPurpose(purpose),
		Status:    openai.FileObjectStatusUploaded,
	}

	common.WriteJSONResponse(ctx, w, http.StatusOK, fileObj)
}

func (c *FilesApiHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	common.WriteNotImplementedError(r.Context(), w)
}

func (c *FilesApiHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	common.WriteNotImplementedError(r.Context(), w)
}

func (c *FilesApiHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	common.WriteNotImplementedError(r.Context(), w)
}

func (c *FilesApiHandler) RetrieveFile(w http.ResponseWriter, r *http.Request) {
	common.WriteNotImplementedError(r.Context(), w)
}
