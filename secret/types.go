package secret


type Interface interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)

}

