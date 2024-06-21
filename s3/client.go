package s3

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/url"
	"time"
)

type Config struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
	Bucket    string
	UseSSL    bool
}

type Client struct {
	config Config
	conn   *minio.Client
}

func New(cfg Config) (*Client, error) {
	cli, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		config: cfg,
		conn:   cli,
	}, nil
}

type PutOptions minio.PutObjectOptions

func (c *Client) PutObject(ctx context.Context, objectName string, reader io.Reader, option PutOptions) error {
	if _, err := c.conn.PutObject(ctx, c.config.Bucket, objectName, reader, -1, minio.PutObjectOptions(option)); err != nil {
		return err
	}

	return nil
}

func (c *Client) PresignedPutObject(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	u, err := c.conn.PresignedPutObject(ctx, c.config.Bucket, objectName, expires)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

type GetOptions minio.GetObjectOptions

func (c *Client) GetObject(ctx context.Context, objectName string, option GetOptions) (io.ReadCloser, error) {
	reader, err := c.conn.GetObject(ctx, c.config.Bucket, objectName, minio.GetObjectOptions(option))
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		reader.Close()
	})

	return reader, nil
}

func (c *Client) PresignedGetObject(ctx context.Context, objectName string, expires time.Duration, reqParams url.Values) (string, error) {
	u, err := c.conn.PresignedGetObject(ctx, c.config.Bucket, objectName, expires, reqParams)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

type RemoveOption minio.RemoveObjectOptions

func (c *Client) RemoveObject(ctx context.Context, objectName string, option RemoveOption) error {
	if err := c.conn.RemoveObject(ctx, c.config.Bucket, objectName, minio.RemoveObjectOptions(option)); err != nil {
		return err
	}

	return nil
}

type ListOption minio.ListObjectsOptions

func (c *Client) ListObjects(ctx context.Context, option ListOption) <-chan minio.ObjectInfo {
	return c.conn.ListObjects(ctx, c.config.Bucket, minio.ListObjectsOptions(option))
}
