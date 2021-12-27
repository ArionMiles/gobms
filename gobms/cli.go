package main

import (
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"

	"github.com/ArionMiles/gobms/pkg/bms"
)

func listRegions() error {
	regions, err := bms.ListRegions()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Region", "Code"})

	for _, region := range *regions {
		t.AppendRow([]interface{}{region[0].Name, region[0].Code})
	}

	t.Render()
	return nil

}

func listTheaters(regionCode, regionName string) error {
	client, err := bms.NewClient(regionCode, regionName)
	if err != nil {
		return err
	}

	theaters, err := client.ListTheaters()
	if err != nil {
		return err
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Theater", "Code"})

	for _, theater := range theaters {
		t.AppendRow([]interface{}{theater.Name, theater.Code})
	}

	t.Render()
	return nil
}

func main() {
	app := &cli.App{
		Name:  "gobms",
		Usage: "CLI for gobms",
		Commands: []*cli.Command{
			{
				Name:  "listRegions",
				Usage: "List all available regions on BookMyShow",
				Action: func(c *cli.Context) error {
					err := listRegions()
					return err
				},
			},
			{
				Name:  "listTheaters",
				Usage: "List all available theaters in a region",
				Action: func(c *cli.Context) error {
					err := listTheaters(c.String("regionCode"), c.String("regionName"))
					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "regionCode",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "regionName",
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
