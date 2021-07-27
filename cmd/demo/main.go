package main

import (
	"context"
	"fmt"

	"giautm.dev/viettelpay"

	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

var reqs = []viettelpay.CheckAccount{
	{MSISDN: "84982612499", CustomerName: "Nguyen Thi Van Giang"},
	{MSISDN: "84362634580", CustomerName: "NGUY THI QUYNH"},
	{MSISDN: "84983647257", CustomerName: "Dinh Thi Quynh"},
	{MSISDN: "84968008909", CustomerName: "Cong Ly"},
}

var reqs2 = []viettelpay.RequestDisbursement{
	{
		TransactionID: viettelpay.GenOrderID(),
		SMSContent:    "hello@giautm.dev",
		MSISDN:        "84982612499",
		CustomerName:  "Nguyen Thi Van Giang",
		Amount:        1000,
	},
	// {MSISDN: "84362634580", CustomerName: "NGUY THI QUYNH"},
	// {MSISDN: "84983647257", CustomerName: "Dinh Thi Quynh"},
	// {MSISDN: "84968008909", CustomerName: "Cong Ly"},
}

func main() {
	ctx := context.Background()

	partnerAPI, err := viettelpay.ProvidePartnerAPI(ctx, nil)
	if err != nil {
		panic(err)
	}

	// result, err := partnerAPI.RequestPayment(ctx, "Test", reqs2...)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(result)

	result, err := partnerAPI.QueryRequests(ctx,
		"01FBK4322AXR0KG5AAQ2E73A6C",
		viettelpay.QueryByMSISDN("84336392248"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
