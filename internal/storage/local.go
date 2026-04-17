package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage はローカルファイルシステムに画像を保存します（開発環境用）。
type LocalStorage struct {
	dir     string // 保存先ディレクトリ (例: "web/static/images/microposts")
	baseURL string // 公開URLのベース (例: "http://localhost:8080/static/images/microposts")
}

// NewLocalStorage は LocalStorage を作成し、保存先ディレクトリを準備します。
func NewLocalStorage(dir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create dir %s: %w", dir, err)
	}
	return &LocalStorage{dir: dir, baseURL: baseURL}, nil
}

func (s *LocalStorage) Upload(_ context.Context, key string, r io.Reader, _ string) error {
	dst, err := os.Create(filepath.Join(s.dir, key))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, r)
	return err
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	err := os.Remove(filepath.Join(s.dir, key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *LocalStorage) PublicURL(key string) string {
	return s.baseURL + "/" + key
}
