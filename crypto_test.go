package viettelpay_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"strings"
	"testing"

	"giautm.dev/viettelpay"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("failed to generate key: %w", err)
		return
	}

	tests := []struct {
		name    string
		wantDst string
		wantErr bool
	}{
		{
			name:    "simple block",
			wantDst: "abcxyz",
		},
		{
			name:    "two blocks",
			wantDst: strings.Repeat("abcxyz", 100),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := bytes.NewBufferString(tt.wantDst)

			buf := bytes.NewBuffer(nil)
			err = viettelpay.Encrypt(buf, src, src.Len(), &key.PublicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			dst := bytes.NewBuffer(nil)
			err = viettelpay.Decrypt(dst, buf, buf.Len(), key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotDst := dst.String(); gotDst != tt.wantDst {
				t.Errorf("Decrypt() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}
