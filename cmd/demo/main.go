package main

import (
	"context"
	"errors"
	"fmt"

	"giautm.dev/viettelpay"

	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

var reqs = []viettelpay.CheckAccount{
	{
		MSISDN:       "84365233899",
		CustomerName: "0365233899",
	},
	// {MSISDN: "84362634580", CustomerName: "NGUY THI QUYNH"},
	// {MSISDN: "84983647257", CustomerName: "Dinh Thi Quynh"},
	// {MSISDN: "84968008909", CustomerName: "Cong Ly"},
}

var reqs2 = []viettelpay.RequestDisbursement{
	{
		TransactionID: viettelpay.GenOrderID(),
		SMSContent:    "giautm",
		MSISDN:        "84365233899",
		CustomerName:  "0365233899",
		Amount:        1000,
		Note:          "giautm note",
	},
	// {MSISDN: "84362634580", CustomerName: "NGUY THI QUYNH"},
	// {MSISDN: "84983647257", CustomerName: "Dinh Thi Quynh"},
	// {MSISDN: "84968008909", CustomerName: "Cong Ly"},
}

func main() {
	ctx := context.Background()

	cfg, err := viettelpay.ProvideConfig(ctx)
	if err != nil {
		panic(err)
	}

	partnerAPI, err := viettelpay.ProvidePartnerAPI(cfg, nil)
	if err != nil {
		panic(err)
	}

	// result, err := partnerAPI.CheckAccount(ctx, viettelpay.GenOrderID(), reqs...)
	// fmt.Println(result)
	// if err != nil {
	// 	panic(err)
	// }

	// orderID := viettelpay.GenOrderID()
	// result2, err := partnerAPI.RequestDisbursement(ctx, orderID, "Test", reqs2...)
	// fmt.Println(result2)
	// if err != nil {
	// 	panic(err)
	// }

	result, err := partnerAPI.QueryRequests(ctx,
		"01FBRYWSNEWB265WEHHEHCDRH4",
		nil,
	)
	fmt.Println(result)

	var batchErr *viettelpay.ViettelPayBatchError
	if errors.As(err, &batchErr) {
		if batchErr.Code == "DISB_SUCCESS" {
			fmt.Println("Chi thành công")
		}
	} else if err != nil {
		// Panic for other error
		panic(err)
	}
}
