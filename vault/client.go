package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type ClientOption struct {
	address string
	token   string
	timeout time.Duration
}

func NewClientOption() *ClientOption {
	return &ClientOption{}
}

func (o *ClientOption) WithAddress(address string) *ClientOption {
	o.address = address
	return o
}

func (o *ClientOption) WithToken(token string) *ClientOption {
	o.token = token
	return o
}

func (o *ClientOption) WithTimeout(timeout time.Duration) *ClientOption {
	o.timeout = timeout
	return o
}

type Client struct {
	client *vault.Client
}

func NewClient(opt *ClientOption) (*Client, error) {
	client, err := vault.New(vault.WithAddress(opt.address), vault.WithRequestTimeout(opt.timeout))
	if err != nil {
		return nil, err
	}

	if err := client.SetToken(opt.token); err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

type Serializable interface {
	Serialize() (map[string]any, error)
	Deserialize(data map[string]any) error
}

type MountOption struct {
	mountPath            string
	requestSpecificToken string
}

type MountedSecret[T Serializable] struct {
	client *vault.Client
	option *MountOption
}

func Mount[T Serializable](client *Client, option MountOption) (*MountedSecret[T], error) {
	return &MountedSecret[T]{
		client: client.client,
		option: &option,
	}, nil
}

func (m *MountedSecret[T]) Write(ctx context.Context, path string, data T) error {
	d, err := data.Serialize()
	if err != nil {
		return err
	}

	opt := []vault.RequestOption{
		vault.WithMountPath(m.option.mountPath),
	}
	if m.option.requestSpecificToken != "" {
		opt = append(opt, vault.WithToken(m.option.requestSpecificToken))
	}

	if _, err := m.client.Secrets.KvV2Write(ctx, path, schema.KvV2WriteRequest{
		Data: d,
	}, opt...); err != nil {
		return err
	}

	return nil
}

func (m *MountedSecret[T]) Read(ctx context.Context, path string) (T, error) {
	result := *new(T)

	opt := []vault.RequestOption{
		vault.WithMountPath(m.option.mountPath),
	}
	if m.option.requestSpecificToken != "" {
		opt = append(opt, vault.WithToken(m.option.requestSpecificToken))
	}

	resp, err := m.client.Secrets.KvV2Read(ctx, path, opt...)
	if err != nil {
		return result, err
	}

	if err := result.Deserialize(resp.Data.Data); err != nil {
		return result, err
	}

	return result, nil
}
