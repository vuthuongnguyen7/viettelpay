package viettelpay

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"gocloud.dev/runtimevar"
)

var ErrNoPEM = errors.New("no PEM data is found")

type KeysOption struct {
	ViettelPublicKey  string
	PartnerPrivateKey string
}

func readPEM(ctx context.Context, varURL string) ([]byte, error) {
	env, err := runtimevar.OpenVariable(ctx, varURL)
	if err != nil {
		return nil, err
	}
	defer env.Close()

	snap, err := env.Latest(ctx)
	if err != nil {
		return nil, err
	}

	if val, ok := snap.Value.([]byte); ok {
		if block, _ := pem.Decode(val); block != nil {
			return block.Bytes, nil
		}
	}

	return nil, ErrNoPEM
}

func partnerKey(ctx context.Context, varURL string) (*rsa.PrivateKey, error) {
	derBytes, err := readPEM(ctx, varURL)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(derBytes)
}

func viettelKey(ctx context.Context, varURL string) (*rsa.PublicKey, error) {
	derBytes, err := readPEM(ctx, varURL)
	if err != nil {
		return nil, err
	}

	if key, err := x509.ParsePKIXPublicKey(derBytes); err != nil {
		return nil, err
	} else if p, ok := key.(*rsa.PublicKey); ok {
		return p, nil
	}

	return nil, errors.New("invalid key type")
}
