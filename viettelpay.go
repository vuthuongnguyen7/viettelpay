package viettelpay

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	"giautm.dev/viettelpay/soap"
	"github.com/oklog/ulid"
)

func DefaultGenID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}

type CheckAccountRequest struct {
	MSISDN       string `json:"msisdn"`
	CustomerName string `json:"customerName"`
}

type CheckAccountResponse struct {
	MSISDN       string `json:"msisdn"`
	CustomerName string `json:"customerName"`
	Package      string `json:"package"`
	Code         string `json:"errorCode"`
	Message      string `json:"errorMsg"`
}

type PartnerAPI interface {
	Process(ctx context.Context, cmd string, request, response interface{}) error
	CheckAccount(ctx context.Context, checks ...CheckAccountRequest) ([]CheckAccountResponse, error)
}

type options struct {
	password    string
	username    string
	serviceCode string
	genID       func() string
}

// A Option sets options such as credentials, tls, etc.
type Option func(*options)

// WithAuth is an Option to set BasicAuth
func WithAuth(username, password, serviceCode string) Option {
	return func(o *options) {
		o.username = username
		o.password = password
		o.serviceCode = serviceCode
	}
}

var ns2Opt = soap.WithNS2("http://partnerapi.bankplus.viettel.com/")

var defaultOptions = options{
	username:    os.Getenv("VIETTELPAY_USERNAME"),
	password:    os.Getenv("VIETTELPAY_PASSWORD"),
	serviceCode: os.Getenv("VIETTELPAY_SERVICE_CODE"),
	genID:       DefaultGenID,
}

type partnerAPI struct {
	client *soap.Client
	opts   *options

	PartnerPrivateKey *rsa.PrivateKey
	ViettelPublicKey  *rsa.PublicKey
}

func NewPartnerAPI(ctx context.Context, url string, opt ...Option) (_ PartnerAPI, err error) {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	var (
		prikey *rsa.PrivateKey
		pubkey *rsa.PublicKey
	)
	if prikey, err = partnerKey(ctx, "file:///workspaces/viettelpay/keys/partner-private-key.pem?decoder=bytes"); err != nil {
		return nil, err
	}

	if pubkey, err = viettelKey(ctx, "file:///workspaces/viettelpay/keys/viettel-public-key.pem?decoder=bytes"); err != nil {
		return nil, err
	}

	return &partnerAPI{
		client: soap.NewClient(url, ns2Opt),
		opts:   &opts,

		PartnerPrivateKey: prikey,
		ViettelPublicKey:  pubkey,
	}, nil
}

func (s *partnerAPI) CheckAccount(ctx context.Context, checks ...CheckAccountRequest) ([]CheckAccountResponse, error) {
	var results []CheckAccountResponse
	if err := s.Process(ctx, "VTP305", checks, &results); err != nil {
		return nil, err
	}
	return results, nil
}

type EnvelopeRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ServiceCode string `json:"serviceCode"`
	OrderID     string `json:"orderId"`
	Data        []byte `json:"data"`
}

type EnvelopeResponse struct {
	Data      json.RawMessage `json:"data"`
	Signature []byte          `json:"signature"`
}

type EnvelopeResponseData struct {
	OrderID string `json:"orderId"`
	Data    []byte `json:"data"`

	Username        string `json:"username"`
	ServiceCode     string `json:"serviceCode"`
	RealServiceCode string `json:"realServiceCode"`

	RequestId string `json:"requestId"`
	TransDate string `json:"transDate"`

	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
}

func (e EnvelopeResponseData) CheckError() error {
	if e.ErrorCode != "00" {
		return &ViettelPayError{
			Code: e.ErrorCode,
			Desc: e.ErrorDesc,
		}
	}

	return nil
}

func (s *partnerAPI) Process(ctx context.Context, cmd string, data, result interface{}) error {
	passwordEncrypted, err := s.Encrypt(([]byte)(s.opts.password))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err = MarshalGzipJSON(buf, data); err != nil {
		return err
	}

	dataJSON, err := json.Marshal(&EnvelopeRequest{
		Data:        buf.Bytes(),
		OrderID:     s.opts.genID(),
		Password:    passwordEncrypted,
		ServiceCode: s.opts.serviceCode,
		Username:    s.opts.username,
	})
	if err != nil {
		return err
	}

	signature, err := s.Sign(dataJSON)
	if err != nil {
		return err
	}

	res, err := s.call(ctx, &Process{
		Cmd:       cmd,
		Data:      string(dataJSON),
		Signature: base64.StdEncoding.EncodeToString(signature),
	})
	if err != nil {
		return err
	}

	var envRes EnvelopeResponse
	err = json.NewDecoder(bytes.NewBufferString(res.Return_)).
		Decode(&envRes)
	if err != nil {
		return err
	}
	if err = s.Verify(envRes.Data, envRes.Signature); err != nil {
		return err
	}

	var envResData EnvelopeResponseData
	if err = json.Unmarshal(envRes.Data, &envResData); err != nil {
		return err
	}

	if err = envResData.CheckError(); err != nil {
		return err
	}

	return UnmarshalGzipJSON(bytes.NewReader(envResData.Data), result)
}
