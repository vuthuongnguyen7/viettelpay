package viettelpay

import (
	"fmt"
)

type BatchError struct {
	Code string `json:"batchErrorCode"`
	Desc string `json:"batchErrorDesc"`
}

func (e BatchError) Error() string {
	return fmt.Sprintf("ViettelPay(%s): %s", e.Code, e.Desc)
}

type Error struct {
	Code string `json:"errorCode"`
	Desc string `json:"errorDesc"`
}

func (e Error) Error() string {
	return fmt.Sprintf("ViettelPay(%s): %s", e.Code, e.Desc)
}
