package viettelpay

import (
	"context"
	"encoding/pem"
	"errors"

	"github.com/sethvargo/go-envconfig"
	"gocloud.dev/runtimevar"
)

var ErrNoPEM = errors.New("no PEM data is found")

// BlockPEM is a slice of bytes where the information is PEM in the
// environment variable.
type BlockPEM []byte

// EnvDecode implements env.Decoder.
func (b *BlockPEM) EnvDecode(val string) error {
	if block, _ := pem.Decode(([]byte)(val)); block != nil {
		*b = block.Bytes
		return nil
	}

	return ErrNoPEM
}

// Bytes returns the underlying bytes.
func (b BlockPEM) Bytes() []byte {
	return []byte(b)
}

type Config struct {
	BaseURL     string `env:"BASE_URL"`
	Username    string `env:"USERNAME"`
	Password    string `env:"PASSWORD"`
	ServiceCode string `env:"SERVICE_CODE"`

	PartnerPrivateKey BlockPEM `env:"PARTNER_PRIVATE_KEY"`
	ViettelPublicKey  BlockPEM `env:"VIETTEL_PUBLIC_KEY"`
}

func ProvidePartnerAPI(ctx context.Context, client HTTPClient) (PartnerAPI, error) {
	var config Config

	l := envconfig.PrefixLookuper("VIETTELPAY_", envconfig.OsLookuper())
	err := envconfig.ProcessWith(ctx, &config, l, resolveSecretFunc)
	if err != nil {
		return nil, err
	}

	keyStore, err := NewKeyStore(config.PartnerPrivateKey, config.ViettelPublicKey)
	if err != nil {
		return nil, err
	}

	return NewPartnerAPI(
		config.BaseURL,
		WithAuth(config.Username, config.Password, config.ServiceCode),
		WithHTTPClient(client),
		WithKeyStore(keyStore),
	)
}

func resolveSecretFunc(ctx context.Context, key, value string) (string, error) {
	v, err := runtimevar.OpenVariable(ctx, value)
	if err != nil {
		return "", err
	}

	snap, err := v.Latest(ctx)
	if err != nil {
		return "", err
	}
	if s, ok := snap.Value.(string); ok {
		return s, nil
	} else if b, ok := snap.Value.([]byte); ok {
		return string(b), nil
	}

	return value, nil
}
