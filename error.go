package viettelpay

import (
	"fmt"
)

type ViettelPayError struct {
	Code string `json:"errorCode"`
	Desc string `json:"errorDesc"`
}

func (e ViettelPayError) Error() string {
	return fmt.Sprintf("ViettelPay(%s): %s", e.Code, e.Desc)
}
