package viettelpay

import (
	"context"
	"encoding/xml"
)

type Process struct {
	XMLName xml.Name `xml:"ns2:process"`

	Cmd       string `xml:"cmd,omitempty" json:"cmd,omitempty"`
	Data      string `xml:"data,omitempty" json:"data,omitempty"`
	Signature string `xml:"signature,omitempty" json:"signature,omitempty"`
}

type ProcessResponse struct {
	XMLName xml.Name `xml:"http://partnerapi.bankplus.viettel.com/ processResponse"`

	Return_ string `xml:"return,omitempty" json:"return,omitempty"`
}

func (s *partnerAPI) call(ctx context.Context, request *Process) (*ProcessResponse, error) {
	response := new(ProcessResponse)
	err := s.client.CallContext(ctx, "''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
