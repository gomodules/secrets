package types

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tamalsaha/ksm-xorm-demo/secret"
	"github.com/tamalsaha/ksm-xorm-demo/secret/provider/gce"
	"strings"
)

type Secret struct {
	ProviderName string `json:"p"`
	KeyInfo      string `json:"k"`
	Data         string `json:"-"` // Value
	Cipher       []byte `json:"c"`
}
func (s *Secret) FromDB(bytes []byte) error {
	if err := json.Unmarshal(bytes, s); err != nil {
		return err
	}

	secProvider, err := getSecretProvider(s.KeyInfo, s.ProviderName)
	if err != nil {
		return err
	}

	val, err := secProvider.Decrypt(s.Cipher)
	if err != nil {
		return err
	}
	s.Data = string(val)

	return nil
}

func (s *Secret) ToDB() ([]byte, error) {
	provider, err := getSecretProvider(s.KeyInfo, s.ProviderName)
	if err != nil {
		return nil, err
	}
	if s.Cipher, err = provider.Encrypt([]byte(s.Data)); err != nil {
		return nil, err
	}
	//providerName.keyInfo.Value
	return json.Marshal(s)
}

func (s *Secret) String() string  {
	return fmt.Sprintf("%v:%v:*",s.ProviderName, s.KeyInfo)
}


func getSecretProvider(key, name string) (secret.Interface, error)  {
	switch strings.ToLower(name) {
	case gce.ProviderType:
		return gce.New(key)
	}
	return nil, errors.Errorf("Unknown provider %v", name)
}
