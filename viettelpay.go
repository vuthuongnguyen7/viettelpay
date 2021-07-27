package viettelpay

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"
	"time"

	"giautm.dev/viettelpay/soap"
	ulid "github.com/oklog/ulid/v2"
)

type CheckAccount struct {
	MSISDN       string `json:"msisdn"`
	CustomerName string `json:"customerName"`
}

type CheckAccountResponse struct {
	CheckAccount
	Package string `json:"package"`
	Code    string `json:"errorCode"`
	Message string `json:"errorMsg"`
}

type RequestDisbursement struct {
	TransactionID string `json:"transId"`
	MSISDN        string `json:"msisdn"`
	CustomerName  string `json:"customerName"`
	Amount        int64  `json:"amount"`
	SMSContent    string `json:"smsContent"`
	Note          string `json:"note"`
}

type RequestDisbursementResponse struct {
	RequestDisbursement
	Code    string `json:"errorCode"`
	Message string `json:"errorMsg"`
}

type RequestDisbursementEnvelope struct {
	EnvelopeBase
	TotalAmount        int64  `json:"totalAmount"`
	TotalTransactions  int    `json:"totalTrans"`
	TransactionContent string `json:"transContent"`
}

type QueryRequestEnvelope struct {
	EnvelopeBase
	QueryType string `json:"queryType,omitempty"`
	QueryData string `json:"queryData,omitempty"`
}

type QueryRequests interface {
	Type() string
	Data() string
}

type QueryByTransaction string

func (q QueryByTransaction) Type() string {
	return "TRANS_ID"
}

func (q QueryByTransaction) Data() string {
	return string(q)
}

type QueryByMSISDN string

func (q QueryByMSISDN) Type() string {
	return "MSISDN"
}

func (q QueryByMSISDN) Data() string {
	return string(q)
}

type PartnerAPI interface {
	Process(ctx context.Context, req Request, response interface{}) error

	CheckAccount(ctx context.Context, orderID string, checks ...CheckAccount) ([]CheckAccountResponse, error)
	RequestDisbursement(ctx context.Context, orderID string, transactionContent string, reqs ...RequestDisbursement) ([]RequestDisbursementResponse, error)
	QueryRequests(ctx context.Context, orderID string, query QueryRequests) ([]RequestDisbursementResponse, error)
}

func GenOrderID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}

// HTTPClient is a client which can make HTTP requests
// An example implementation is net/http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type options struct {
	password    string
	username    string
	serviceCode string

	keyStore   KeyStore
	httpClient soap.HTTPClient
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

// WithHTTPClient is an Option to set the HTTP client to use
func WithHTTPClient(c HTTPClient) Option {
	return func(o *options) {
		o.httpClient = c
	}
}

// WithKeyStore is an Option to set BasicAuth
func WithKeyStore(keyStore KeyStore) Option {
	return func(o *options) {
		o.keyStore = keyStore
	}
}

var ns2Opt = soap.WithNS2("http://partnerapi.bankplus.viettel.com/")

var defaultOptions = options{}

type partnerAPI struct {
	client *soap.Client

	password    string
	username    string
	serviceCode string

	keyStore KeyStore
}

func NewPartnerAPI(url string, opt ...Option) (_ PartnerAPI, err error) {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	if opts.keyStore == nil {
		return nil, errors.New("missing keyStore option")
	}

	return &partnerAPI{
		client: soap.NewClient(url, ns2Opt, soap.WithHTTPClient(opts.httpClient)),

		keyStore:    opts.keyStore,
		username:    opts.username,
		password:    opts.password,
		serviceCode: opts.serviceCode,
	}, nil
}

func (s *partnerAPI) CheckAccount(ctx context.Context, orderID string, checks ...CheckAccount) ([]CheckAccountResponse, error) {
	results := []CheckAccountResponse{}
	env := &EnvelopeBase{}
	env.OrderID = orderID
	err := s.Process(ctx, NewRequest("VTP305", checks, env), &results)
	return results, err
}

func (s *partnerAPI) RequestDisbursement(ctx context.Context, orderID string, transactionContent string, reqs ...RequestDisbursement) ([]RequestDisbursementResponse, error) {
	env := &RequestDisbursementEnvelope{
		TotalTransactions:  len(reqs),
		TransactionContent: transactionContent,
	}
	env.OrderID = orderID
	for _, v := range reqs {
		env.TotalAmount += v.Amount
	}

	results := []RequestDisbursementResponse{}
	err := s.Process(ctx, NewRequest("VTP306", reqs, env), &results)
	return results, err
}

var emptyArray = []interface{}{}

func (s *partnerAPI) QueryRequests(ctx context.Context, orderID string, query QueryRequests) ([]RequestDisbursementResponse, error) {
	env := &QueryRequestEnvelope{}
	env.OrderID = orderID
	if query != nil {
		env.QueryType, env.QueryData = query.Type(), query.Data()
	}

	results := []RequestDisbursementResponse{}
	err := s.Process(ctx, NewRequest("VTP307", emptyArray, env), &results)
	return results, err
}
