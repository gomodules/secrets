package types

import (
	"context"
	"encoding/json"
	"fmt"
	"gocloud.dev/secrets"

)

type Secret struct {
	Url string `json:"u"`
	Data         string `json:"-"` // Value
	Cipher       []byte `json:"c"`
}
func (s *Secret) FromDB(bytes []byte) error {
	if err := json.Unmarshal(bytes, s); err != nil {
		return err
	}

	ctx := context.Background()
	k, err := secrets.OpenKeeper(ctx, s.Url)
	if err != nil {
		return err
	}

	val, err := k.Decrypt(ctx, s.Cipher)
	if err != nil {
		return err
	}
	s.Data = string(val)

	return nil
}

func (s *Secret) ToDB() ([]byte, error) {
	ctx := context.Background()
	k, err := secrets.OpenKeeper(ctx, s.Url)
	if err != nil {
		return nil, err
	}
	if s.Cipher, err = k.Encrypt(ctx, []byte(s.Data)); err != nil {
		return nil, err
	}
	//providerName.keyInfo.Value
	return json.Marshal(s)
}

func (s *Secret) String() string  {
	return fmt.Sprintf("%v:%v",s.Url, s.Data)
}

