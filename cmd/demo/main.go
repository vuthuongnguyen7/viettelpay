package main

import (
	"context"
	"fmt"

	"giautm.dev/viettelpay"

	_ "gocloud.dev/runtimevar/filevar"
)

var reqs = []viettelpay.CheckAccountRequest{
	{MSISDN: "84982612499", CustomerName: "Nguyen Thi Van Giang"},
	// {MSISDN: "84362634580", CustomerName: "NGUY THI QUYNH"},
	// {MSISDN: "84983647257", CustomerName: "Dinh Thi Quynh"},
	// {MSISDN: "84968008909", CustomerName: "Cong Ly"},
}

func main() {
	ctx := context.Background()

	partnerAPI, err := viettelpay.NewPartnerAPI(ctx,
		"https://wallet.viettelpay.vn/uat/PartnerWS/PartnerAPI?wsdl",
		viettelpay.WithAuth("HCMHYT8888", "HCMHYT8888@123", "HCMHYT8888"),
	)
	if err != nil {
		panic(err)
	}

	result, err := partnerAPI.CheckAccount(ctx, reqs...)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
