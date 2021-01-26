package main

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	nutanix "github.com/routebyintuition/ntnx-go-sdk"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// NCLI is the baseline type to house SDK/connection config
type NCLI struct {
	con *nutanix.Client
	tr  *tablewriter.Table
}

func main() {

	ncli := &NCLI{}
	ncli.tr = tablewriter.NewWriter(os.Stdout)

	flags := getFlags()

	app := &cli.App{
		Name:                 "Unikum Wunderbar Nutanix CLI",
		Usage:                "Built for the Unikum und Wunderbar",
		EnableBashCompletion: true,
		Before:               altsrc.InitInputSourceWithContext(flags, NewYamlSourceFromProfileFunc("profile")),
		Flags:                flags,
		Commands: []*cli.Command{
			{
				Name:    "configure",
				Aliases: []string{"conf"},
				Usage:   "configure stored credentials",
				Action:  saveCredentials,
			},
			{
				Name:   "list-profiles",
				Usage:  "list saved profiles",
				Action: ncli.listProfiles,
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name: "vm",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "retrieve all VMs",
						Action: ncli.vmList,
					},
					{
						Name:   "get",
						Usage:  "get one VM by UUID",
						Action: ncli.vmGet,
					},
					{
						Name:   "disklist",
						Usage:  "get disk list of VM by UUID",
						Action: ncli.vmDiskList,
					},
					{
						Name:   "update-memory",
						Usage:  "<UUID> <memory in MB integer>",
						Action: ncli.vmMemoryUpdate,
					},
				},
				Category: "vm",
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name: "image",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list all images",
						Action: ncli.imageList,
					},
					{
						Name:   "create",
						Usage:  "create a new image",
						Action: ncli.imageCreate,
					},
				},
				Category: "image",
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name: "cluster",
				Subcommands: []*cli.Command{
					{
						Name:  "all",
						Usage: "retrieve all VMs",
						Action: func(c *cli.Context) error {
							fmt.Println("new task template: ", c.Args().First())
							fmt.Println("user: ", c.String("user"))
							vmList(ncli.con.PC, c.String("vmname"), 10)
							return nil
						},
					},
				},
				Category: "cluster",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
