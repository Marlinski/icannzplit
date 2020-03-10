package main

import (
	"icannzplit/vpn/ipvanish"
	"icannzplit/zplit"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "/icannzplit",
			Usage: "specify icannzplit home directory",
		},
		cli.StringFlag{
			Name:  "ipvanish-user",
			Value: "**IPVANISH_USERNAME**",
			Usage: "specify ipvanish username",
		},
		cli.StringFlag{
			Name:  "ipvanish-pass",
			Value: "**IPVANISH_PASSWORD**",
			Usage: "specify ipvanish password",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Printf(" I C A N N    M U L T I P L E X ! \n")
		log.Printf(" =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- \n")

		// load ipvanish config
		err := ipvanish.Init(c.String("homedir"), c.String("ipvanish-user"), c.String("ipvanish-pass"), c.Int("ipvanish-table"))
		if err != nil {
			log.Fatalf("fail to setup ipvanish: %+v", err)
		}

		// load icannsplit config
		cfg := zplit.ConfigInit(c.String("homedir"))
		cfg.Save(c.String("homedir"))

		// build routing plan
		plan := cfg.BuildPlan()
		plan.Dump()

		// connect all the routes
		plan.Execute()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}