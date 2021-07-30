package viettelpay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

type KeyStore interface {
	Sign(data []byte) (signature []byte, err error)
	Verify(data, signature []byte) (err error)
	Decrypt(msg []byte) (string, error)
	Encrypt(msg []byte) (string, error)
}

type keyStore struct {
	partnerPrivateKey *rsa.PrivateKey
	viettelPublicKey  *rsa.PublicKey
}

func NewKeyStore(partnerPriKey, viettelPubKey []byte) (_ KeyStore, err error) {
	keys := &keyStore{}

	if keys.partnerPrivateKey, err = x509.ParsePKCS1PrivateKey(partnerPriKey); err != nil {
		return nil, err
	}

	if key, err := x509.ParsePKIXPublicKey(viettelPubKey); err != nil {
		return nil, err
	} else if rsaKey, ok := key.(*rsa.PublicKey); ok {
		keys.viettelPublicKey = rsaKey
	} else {
		return nil, errors.New("invalid key type")
	}

	return keys, nil
}

func (s *keyStore) Sign(data []byte) ([]byte, error) {
	hashed := sha1.Sum(data)
	return s.partnerPrivateKey.Sign(rand.Reader, hashed[:], crypto.SHA1)
}

func (s *keyStore) Verify(data, signature []byte) error {
	hashed := sha1.Sum(data)
	return rsa.VerifyPKCS1v15(s.viettelPublicKey, crypto.SHA1, hashed[:], signature)
}

func (s *keyStore) Decrypt(msg []byte) (string, error) {
	return Decrypt(msg, s.partnerPrivateKey)
}

func (s *keyStore) Encrypt(msg []byte) (string, error) {
	return Encrypt(msg, s.viettelPublicKey)
}

func Decrypt(msg []byte, privateKey *rsa.PrivateKey) (string, error) {
	result := bytes.NewBuffer(nil)

	keySize := privateKey.Size()
	base64BlockSize := base64.StdEncoding.EncodedLen(keySize)
	iterations := len(msg) / base64BlockSize

	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(base64BlockSize))
	r := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(msg))
	for i := 0; i < iterations; i++ {
		n, err := r.Read(ciphertext)
		if err != nil {
			return "", err
		}

		reverseBytes(ciphertext[:n])
		plaintext, err := privateKey.Decrypt(rand.Reader, ciphertext[:n], nil)
		if err != nil {
			return "", err
		}

		result.Write(plaintext)
	}

	return result.String(), nil
}

func Encrypt(msg []byte, publicKey *rsa.PublicKey) (string, error) {
	keySize := publicKey.Size()
	maxLength := keySize - 42
	dataLength := len(msg)
	iterations := dataLength / maxLength

	result := ""
	for i := 0; i <= iterations; i++ {
		last := (i + 1) * maxLength
		if last > dataLength {
			last = dataLength
		}

		ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg[i*maxLength:last])
		if err != nil {
			return "", err
		}

		reverseBytes(ciphertext)
		result += base64.StdEncoding.EncodeToString(ciphertext)
	}

	return result, nil
}

func reverseBytes(in []byte) []byte {
	for l, r := 0, len(in)-1; l < r; {
		in[l], in[r] = in[r], in[l]
		l++
		r--
	}

	return in
}
