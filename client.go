package viettelpay

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
)

type EnvelopeBase struct {
	Password    string `json:"password"`
	ServiceCode string `json:"serviceCode"`
	Username    string `json:"username"`

	Data    []byte `json:"data"`
	OrderID string `json:"orderId"`
}

func (e *EnvelopeBase) SetData(val []byte) {
	e.Data = val
}
func (e *EnvelopeBase) SetUsername(val string) {
	e.Username = val
}
func (e *EnvelopeBase) SetPassword(val string) {
	e.Password = val
}
func (e *EnvelopeBase) SetServiceCode(val string) {
	e.ServiceCode = val
}

type EnvelopeResponse struct {
	Data      json.RawMessage `json:"data"`
	Signature []byte          `json:"signature"`
}

type EnvelopeResponseData struct {
	Data    []byte `json:"data"`
	OrderID string `json:"orderId"`

	RealServiceCode string `json:"realServiceCode"`
	ServiceCode     string `json:"serviceCode"`
	Username        string `json:"username"`

	RequestId string `json:"requestId"`
	TransDate string `json:"transDate"`

	BatchErrorCode string `json:"batchErrorCode"`
	BatchErrorDesc string `json:"batchErrorDesc"`

	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
}

func (e EnvelopeResponseData) CheckError() error {
	if e.ErrorCode != "00" {
		return &ViettelPayError{
			Code: e.ErrorCode,
			Desc: e.ErrorDesc,
		}
	} else if e.BatchErrorCode != "" {
		return &ViettelPayBatchError{
			Code: e.ErrorCode,
			Desc: e.ErrorDesc,
		}
	}

	return nil
}

func (s *partnerAPI) Process(ctx context.Context, req Request, result interface{}) error {
	passwordEncrypted, err := s.Encrypt(([]byte)(s.opts.password))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err = MarshalGzipJSON(buf, req.Data()); err != nil {
		return err
	}

	envReq := req.Envelope()
	envReq.SetData(buf.Bytes())
	envReq.SetPassword(passwordEncrypted)
	envReq.SetServiceCode(s.opts.serviceCode)
	envReq.SetUsername(s.opts.username)

	envReqJSON, err := json.Marshal(envReq)
	if err != nil {
		return err
	}

	signature, err := s.Sign(envReqJSON)
	if err != nil {
		return err
	}

	res, err := s.call(ctx, &Process{
		Cmd:       req.Command(),
		Data:      string(envReqJSON),
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

	// NOTE: VTP also return data in case errors happen.
	// So, we unmarshal data first then check error late.
	err = UnmarshalGzipJSON(bytes.NewReader(envResData.Data), result)
	if err != nil {
		return err
	}

	return envResData.CheckError()
}
