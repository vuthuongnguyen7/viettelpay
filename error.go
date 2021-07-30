package viettelpay

import (
	"fmt"
)

type BatchError struct {
	Code string `json:"batchErrorCode"`
	Desc string `json:"batchErrorDesc"`
}

var _ error = (*BatchError)(nil)

func (e *BatchError) Is(target error) bool {
	t, ok := target.(*BatchError)
	if !ok {
		return false
	}

	return e.Code == t.Code
}

func (e BatchError) Error() string {
	return fmt.Sprintf("ViettelPay(%s): %s", e.Code, e.Desc)
}

type Error struct {
	Code string `json:"errorCode"`
	Desc string `json:"errorDesc"`
}

var _ error = (*Error)(nil)

func (e Error) Error() string {
	return fmt.Sprintf("ViettelPay(%s): %s", e.Code, e.Desc)
}

var (
	ErrBatchWaitDisb     = &BatchError{Code: "WAIT_DISB"}
	ErrBatchCancelDisb   = &BatchError{Code: "CANCEL_DISB"}
	ErrBatchDisbursement = &BatchError{Code: "DISBURSEMENT"}
	ErrBatchDisbTimeout  = &BatchError{Code: "DISB_TIMEOUT"}
	ErrBatchDisbSuccess  = &BatchError{Code: "DISB_SUCCESS"}
	ErrBatchDisbFailed   = &BatchError{Code: "DISB_FAILED"}
)
