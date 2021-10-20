package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"giautm.dev/viettelpay"
	"github.com/urfave/cli/v2"

	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

// main entry point for the temporal server
func main() {
	app := buildCLI()
	_ = app.Run(os.Args)
}

// buildCLI is the main entry point for the temporal server
func buildCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "viettelpay"
	app.Usage = "Viettel Pay Tools"
	app.ArgsUsage = " "
	app.Flags = []cli.Flag{}

	app.Commands = []*cli.Command{
		{
			Name:  "verify",
			Usage: "Verify VTP account",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "msisdn",
					Aliases:  []string{"m"},
					Usage:    "MSISDN",
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

				partnerAPI, err := initialClient(ctx)
				if err != nil {
					return err
				}

				orderID := viettelpay.GenOrderID()
				results, err := partnerAPI.CheckAccount(ctx, orderID, viettelpay.CheckAccount{
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

				return cli.Exit("All services are stopped.", 0)
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
			Before: func(c *cli.Context) error {
				if c.Args().Len() > 0 {
					return cli.Exit("ERROR: start command doesn't support arguments. Use --service flag instead.", 1)
				}
				return nil
			},
			Action: func(c *cli.Context) error {
				ctx := c.Context
				orderID := c.String("orderID")

				partnerAPI, err := initialClient(ctx)
				if err != nil {
					return err
				}

				results, err := partnerAPI.QueryRequests(ctx, orderID, nil)
				for _, r := range results {
					fmt.Printf("%s - %s - %s \n", r.TransactionID, r.ErrorCode, r.ErrorMsg)
				}
				if errors.Is(err, viettelpay.ErrBatchDisbSuccess) {
					fmt.Println("Chi thành công")
				} else if err != nil {
					// Panic for other error
					return cli.Exit(fmt.Sprintf("Unable to query result. Error: %v", err), 1)
				}

				return cli.Exit("All services are stopped.", 0)
			},
		},
	}
	return app
}

func initialClient(ctx context.Context) (viettelpay.PartnerAPI, error) {
	cfg, err := viettelpay.ProvideConfig(ctx)
	if err != nil {
		return nil, cli.Exit(fmt.Sprintf("Failed to read config. Error: %v", err), 1)
	}

	partnerAPI, err := viettelpay.ProvidePartnerAPI(cfg, nil)
	if err != nil {
		return nil, cli.Exit(fmt.Sprintf("Failed to initial partner api. Error: %v", err), 1)
	}

	return partnerAPI, nil
}
