package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Storage は Cloudflare R2 に画像を保存します（本番環境用）。
// R2 は S3 互換 API を持つため aws-sdk-go-v2 で操作できます。
type R2Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string // R2 の公開URL (例: "https://pub-xxx.r2.dev" またはカスタムドメイン)
	prefix    string // オブジェクトキーのプレフィックス (例: "microposts")
}

// NewR2Storage は R2Storage を作成します。
//
//   - accountID:       Cloudflare アカウントID
//   - accessKeyID:     R2 APIトークンのアクセスキーID
//   - secretAccessKey: R2 APIトークンのシークレット
//   - bucket:          R2 バケット名
//   - publicURL:       R2 パブリックアクセスURL (末尾スラッシュなし)
//   - prefix:          オブジェクトキーのプレフィックス ("microposts" 等、不要なら空文字)
func NewR2Storage(accountID, accessKeyID, secretAccessKey, bucket, publicURL, prefix string) (*R2Storage, error) {
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", strings.TrimSpace(accountID))

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})

	return &R2Storage{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
		prefix:    prefix,
	}, nil
}

func (s *R2Storage) objectKey(key string) string {
	if s.prefix != "" {
		return s.prefix + "/" + key
	}
	return key
}

func (s *R2Storage) Upload(ctx context.Context, key string, r io.Reader, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(s.objectKey(key)),
		Body:        r,
		ContentType: aws.String(contentType),
	})
	return err
}

func (s *R2Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.objectKey(key)),
	})
	return err
}

func (s *R2Storage) PublicURL(key string) string {
	if s.prefix != "" {
		return s.publicURL + "/" + s.prefix + "/" + key
	}
	return s.publicURL + "/" + key
}
