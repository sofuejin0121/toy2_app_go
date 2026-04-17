package storage

import (
	"context"
	"io"
)

// ImageStorage は画像の保存先を抽象化するインターフェースです。
// 開発環境ではローカルファイルシステム、本番環境では Cloudflare R2 に切り替えます。
type ImageStorage interface {
	// Upload は画像を key (ファイル名) で保存します。
	Upload(ctx context.Context, key string, r io.Reader, contentType string) error
	// Delete は key で指定された画像を削除します。存在しない場合は nil を返します。
	Delete(ctx context.Context, key string) error
	// PublicURL は key から公開アクセス可能な URL を返します。
	PublicURL(key string) string
}
