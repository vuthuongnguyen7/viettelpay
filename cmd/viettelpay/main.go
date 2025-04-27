package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	vtp "giautm.dev/viettelpay"
	"github.com/urfave/cli/v2"

	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

func main() {
	app := buildCLI()
	_ = app.Run(os.Args)
}

func buildCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "viettelpay"
	app.Usage = "Viettel Pay Tools"
	app.ArgsUsage = " "
	app.Flags = []cli.Flag{}

	app.Commands = []*cli.Command{
		{
			Name:  "pham hoang huong",
			Usage: "pham quang hoc",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "01007784025",
					Aliases:  []string{"m"},
					Usage:    "841007784025",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "Customer Name",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				ctx := c.Context
				msisdn := c.String("msisdn")
				customerName := c.String("name")

				client, err := initialClient(ctx)
				if err != nil {
					return err
				}

				orderID := vtp.GenOrderID()
				results, err := client.CheckAccount(ctx, orderID, vtp.CheckAccount{
					MSISDN:       msisdn,
					CustomerName: customerName,
				})
				for _, r := range results {
					fmt.Printf("%s - %s: %s\n", r.MSISDN, r.ErrorCode, r.ErrorDesc)
				}
				if err != nil {
					// Panic for other error
					return cli.Exit(fmt.Sprintf("Unable to query result. Error: %v", err), 1)
				}

				return nil
			},
		},
		{
			Name:      "query",
			Usage:     "Query result of Request Disbursement",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "orderID",
					Aliases:  []string{"o"},
					Usage:    "Order ID to query the result",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				ctx := c.Context
				orderID := c.String("orderID")

				client, err := initialClient(ctx)
				if err != nil {
					return err
				}

				results, err := client.QueryRequests(ctx, orderID, nil)
				for _, r := range results {
					fmt.Printf("%s - %v\n", r.TransactionID, r.Err())
				}

				var batchErr *vtp.BatchError
				if errors.As(err, &batchErr) {
					if batchErr.Is(vtp.ErrBatchDisbSuccess) {
						fmt.Println("Chi thành công")
					} else {
						fmt.Println(batchErr.Error())
					}
				} else if err != nil {
					// Panic for other error
					return cli.Exit(fmt.Sprintf("Unable to query result. Error: %v", err), 1)
				}

				return nil
			},
		},
	}
	return app
}

func initialClient(ctx context.Context) (vtp.PartnerAPI, error) {
	cfg, err := vtp.ProvideConfig(ctx)
	if err != nil {
		return nil, cli.Exit(fmt.Sprintf("Failed to read config. Error: %v", err), 1)
	}

	partnerAPI, err := vtp.ProvidePartnerAPI(cfg, nil)
	if err != nil {
		return nil, cli.Exit(fmt.Sprintf("Failed to initial partner api. Error: %v", err), 1)
	}

	return partnerAPI, nil
}
