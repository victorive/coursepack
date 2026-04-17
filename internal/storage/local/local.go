package local

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

func (ls *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	filePath := filepath.Join(ls.basePath, key)
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return key, nil
}

func (ls *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	filePath := filepath.Join(ls.basePath, key)

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

func (ls *LocalStorage) Delete(ctx context.Context, key string) error {
	filePath := filepath.Join(ls.basePath, key)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (ls *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	filePath := filepath.Join(ls.basePath, key)
	_, err := os.Stat(filePath)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func (ls *LocalStorage) GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	hash := sha256.Sum256([]byte(key + time.Now().Add(expiry).String()))
	token := fmt.Sprintf("%x", hash)

	return fmt.Sprintf("/api/v1/files/download?key=%s&token=%s", key, token[:16]), nil
}

func (ls *LocalStorage) GetSize(ctx context.Context, key string) (int64, error) {
	filePath := filepath.Join(ls.basePath, key)

	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("file not found")
		}
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return stat.Size(), nil
}

func (ls *LocalStorage) GetContentType(ctx context.Context, key string) (string, error) {
	ext := filepath.Ext(key)
	contentType := mime.TypeByExtension(ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return contentType, nil
}

func (ls *LocalStorage) GetFilePath(key string) string {
	return filepath.Join(ls.basePath, key)
}

func (ls *LocalStorage) GetBasePath() string {
	return ls.basePath
}
