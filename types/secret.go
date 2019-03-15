package types

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tamalsaha/ksm-xorm-demo/secret"
	"github.com/tamalsaha/ksm-xorm-demo/secret/provider/gce"
	"strings"
)

type Secret struct {
	ProviderName string
	KeyInfo      string
	Value        string
}

func (s *Secret) FromDB(bytes []byte) error {
	if err := s.decode(bytes); err != nil {
		return err
	}

	secProvider, err := getSecretProvider(s.KeyInfo, s.ProviderName)
	if err != nil {
		return err
	}

	val, err := secProvider.Decrypt([]byte(s.Value))
	if err != nil {
		return err
	}
	s.Value = string(val)

	return nil
}

func (s *Secret) ToDB() ([]byte, error) {
	provider, err := getSecretProvider(s.KeyInfo, s.ProviderName)
	if err != nil {
		return nil, err
	}
	data, err := provider.Encrypt([]byte(s.Value))
	if err != nil {
		return nil, err
	}
	//providerName.keyInfo.Value
	ret := fmt.Sprintf("%v.%v.%v", s.ProviderName,
		base64.StdEncoding.EncodeToString([]byte(s.KeyInfo)),
		base64.StdEncoding.EncodeToString(data))

	return []byte(ret), nil
}

func (s *Secret) String() string  {
	bytes, err := s.ToDB()
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (s *Secret) decode(bytes []byte) error  {
	data := string(bytes)
	res := strings.Split(data, ".")
	//providerName.keyInfo.Value

	s.ProviderName = res[0]

	keyData, err := base64.StdEncoding.DecodeString(res[1])
	if err != nil {
		return err
	}
	s.KeyInfo = string(keyData)

	valData, err := base64.StdEncoding.DecodeString(res[2])
	if err != nil {
		return err
	}

	s.Value = string(valData)

	return nil
}

func getSecretProvider(key, name string) (secret.Interface, error)  {
	switch strings.ToLower(name) {
	case gce.ProviderType:
		return gce.New(key)
	}
	return nil, errors.Errorf("Unknown provider %v", name)
}