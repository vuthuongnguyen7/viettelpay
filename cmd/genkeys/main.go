package main

import (
	"fmt"
	"os"

	"giautm.dev/viettelpay"
)

func main() {
	prvKeyFile, err := os.Create("private.pem")
	checkError(err)
	defer prvKeyFile.Close()

	pubKeyFile, err := os.Create("public.pem")
	checkError(err)
	defer pubKeyFile.Close()

	err = viettelpay.GenerateKeysPEM(prvKeyFile, pubKeyFile, 1024)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
