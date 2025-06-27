package hasher

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	// кодируем salt + hash вместе
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s:%s", saltEncoded, hashEncoded), nil
}

func Verify(hashedPwd, pwd string) error {
	parts := bytes.SplitN([]byte(hashedPwd), []byte{':'}, 2)
	if len(parts) != 2 {
		return errors.New("invalid hash format")
	}
	salt, err := base64.RawStdEncoding.DecodeString(string(parts[0]))
	if err != nil {
		return err
	}
	expected, err := base64.RawStdEncoding.DecodeString(string(parts[1]))
	if err != nil {
		return err
	}

	got := argon2.IDKey([]byte(pwd), salt, 1, 64*1024, 4, 32)
	if subtle.ConstantTimeCompare(expected, got) != 1 {
		return errors.New("логин или пароль неправильный")
	}
	return nil
}
