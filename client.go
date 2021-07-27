package viettelpay

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type EnvelopeBase struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ServiceCode string `json:"serviceCode"`
	OrderID     string `json:"orderId"`
	Data        []byte `json:"data"`
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

func (e EnvelopeResponseData) CheckError() *ViettelPayError {
	if e.ErrorCode != "00" {
		return &ViettelPayError{
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

	fmt.Println(res.Return_)
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

	if err := envResData.CheckError(); err != nil {
		// Always return error for PXX code
		if len(err.Code) == 3 && err.Code[0] == 'P' {
			return err
		}
	}

	return UnmarshalGzipJSON(bytes.NewReader(envResData.Data), result)
}
