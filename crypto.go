package viettelpay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
)

func (s *partnerAPI) Sign(data []byte) ([]byte, error) {
	hashed := sha1.Sum(data)
	return s.PartnerPrivateKey.Sign(rand.Reader, hashed[:], crypto.SHA1)
}

func (s *partnerAPI) Verify(data, signature []byte) error {
	hashed := sha1.Sum(data)
	return rsa.VerifyPKCS1v15(s.ViettelPublicKey, crypto.SHA1, hashed[:], signature)
}

func (s *partnerAPI) Encrypt(msg []byte) (string, error) {
	keySize := s.ViettelPublicKey.Size()
	maxLength := keySize - 42
	dataLength := len(msg)
	iterations := dataLength / maxLength

	data := ""
	for i := 0; i <= iterations; i++ {
		last := (i + 1) * maxLength
		if last > dataLength {
			last = dataLength
		}

		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, s.ViettelPublicKey, msg[i*maxLength:last])
		if err != nil {
			return "", err
		}

		reverseBytes(bytes)
		data += base64.StdEncoding.EncodeToString(bytes)
	}

	return data, nil
}

func reverseBytes(in []byte) []byte {
	for l, r := 0, len(in)-1; l < r; {
		in[l], in[r] = in[r], in[l]
		l++
		r--
	}

	return in
}
