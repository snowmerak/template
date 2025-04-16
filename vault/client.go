package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type ClientOption struct {
	address  string
	token    string
	roleId   string
	secretId string
	timeout  time.Duration
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

func (o *ClientOption) WithRoleId(roleId, secretId string) *ClientOption {
	o.roleId = roleId
	o.secretId = secretId
	return o
}

func (o *ClientOption) WithTimeout(timeout time.Duration) *ClientOption {
	o.timeout = timeout
	return o
}

type Client struct {
	client *vault.Client
}

func NewClient(ctx context.Context, opt *ClientOption) (*Client, error) {
	client, err := vault.New(vault.WithAddress(opt.address), vault.WithRequestTimeout(opt.timeout))
	if err != nil {
		return nil, err
	}

	switch len(opt.token) {
	case 0:
		resp, err := client.Auth.AppRoleLogin(ctx, schema.AppRoleLoginRequest{
			RoleId:   opt.roleId,
			SecretId: opt.secretId,
		})
		if err != nil {
			return nil, err
		}

		if err := client.SetToken(resp.Auth.ClientToken); err != nil {
			return nil, err
		}
	default:
		if err := client.SetToken(opt.token); err != nil {
			return nil, err
		}
	}

	return &Client{client: client}, nil
}

type Serializable[T any] interface {
	Serialize() (map[string]any, error)
	Deserialize(data map[string]any) error
	Init() T
}

type MountOption struct {
	mountPath            string
	requestSpecificToken string
}

func NewMountOption() *MountOption {
	return &MountOption{}
}

func (o *MountOption) WithMountPath(mountPath string) *MountOption {
	o.mountPath = mountPath
	return o
}

func (o *MountOption) WithRequestSpecificToken(token string) *MountOption {
	o.requestSpecificToken = token
	return o
}

type MountedSecret[T Serializable[T]] struct {
	client *vault.Client
	option *MountOption
}

func Mount[T Serializable[T]](client *Client, option *MountOption) (*MountedSecret[T], error) {
	return &MountedSecret[T]{
		client: client.client,
		option: option,
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

	result = result.Init()
	if err := result.Deserialize(resp.Data.Data); err != nil {
		return result, err
	}

	return result, nil
}

func (m *MountedSecret[T]) Delete(ctx context.Context, path string) error {
	opt := []vault.RequestOption{
		vault.WithMountPath(m.option.mountPath),
	}
	if m.option.requestSpecificToken != "" {
		opt = append(opt, vault.WithToken(m.option.requestSpecificToken))
	}

	if _, err := m.client.Secrets.KvV2Delete(ctx, path, opt...); err != nil {
		return err
	}

	return nil
}
