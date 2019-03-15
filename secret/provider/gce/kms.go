package gce

import (
	cloudkms "cloud.google.com/go/kms/apiv1"
	"context"
	"fmt"
	"github.com/tamalsaha/ksm-xorm-demo/secret"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	ProviderType = "gce"
)

type Cryptographic struct {
	ctx context.Context
	client *cloudkms.KeyManagementClient
	KeyName string

}

func New(keyName string) (secret.Interface, error)  {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Cryptographic{
		client:client,
		ctx: ctx,
		KeyName:keyName,
	}, nil
}

func (c *Cryptographic) Encrypt(text []byte) ([]byte, error)  {
	fmt.Println(string(text))
	req := &kmspb.EncryptRequest{
		Name: c.KeyName,
		Plaintext: text,
	}

	resp, err := c.client.Encrypt(c.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Ciphertext, nil
}

func (c *Cryptographic) Decrypt(cipher []byte) ([]byte, error)  {
	req := &kmspb.DecryptRequest{
		Name: c.KeyName,
		Ciphertext: cipher,
	}

	resp, err := c.client.Decrypt(c.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}
