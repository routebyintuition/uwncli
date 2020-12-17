package main

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	nutanix "github.com/routebyintuition/ntnx-go-sdk"
	"github.com/urfave/cli/v2"
)

// NCLI is the baseline type to house SDK/connection config
type NCLI struct {
	con *nutanix.Client
	tr  *tablewriter.Table
}

func main() {

	ncli := &NCLI{}
	ncli.tr = tablewriter.NewWriter(os.Stdout)

	app := &cli.App{
		Name:                 "Unikum Wunderbar Nutanix CLI",
		Usage:                "Built for those who are Unikum Wunderbar",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "url",
				Aliases:     []string{"w"},
				Value:       "https://localhost:9440/api/nutanix/v3/",
				Usage:       "Prism Central URL",
				DefaultText: "https://localhost:9440/api/nutanix/v3/",
				EnvVars:     []string{"NUTANIX_PC_URL"},
			},
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u", "username"},
				Value:       "{pc username}",
				Usage:       "Prism Central Username",
				DefaultText: "<username>",
				EnvVars:     []string{"NUTANIX_PC_USER"},
			},
			&cli.StringFlag{
				Name:        "pass",
				Aliases:     []string{"p", "password"},
				Value:       "{pc password}",
				Usage:       "Prism Central Password",
				DefaultText: "<password>",
				EnvVars:     []string{"NUTANIX_PC_PASS"},
			},
			&cli.StringFlag{
				Name:        "image-name",
				Aliases:     []string{"iname", "in"},
				Value:       "",
				Usage:       "Image Name",
				DefaultText: "<image name>",
			},
			&cli.StringFlag{
				Name:        "image-description",
				Aliases:     []string{"idesc", "id"},
				Value:       "",
				Usage:       "Image Description",
				DefaultText: "<image description>",
			},
			&cli.StringFlag{
				Name:        "image-source",
				Aliases:     []string{"isrc", "is"},
				Value:       "",
				Usage:       "Image Source",
				DefaultText: "<image source>",
			},
			&cli.StringFlag{
				Name:        "image-type",
				Aliases:     []string{"itype", "it"},
				Value:       "",
				Usage:       "Image Type - DISK_IMAGE or ISO_IMAGE",
				DefaultText: "<image type>",
			},
			&cli.BoolFlag{
				Name:    "skip-cert-verify",
				Aliases: []string{"skipverify", "scv"},
				Value:   false,
				Usage:   "Image Source",
			},
		},
		Commands: []*cli.Command{
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c.Bool("skip-cert-verify"))
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
				},
				Category: "vm",
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c.Bool("skip-cert-verify"))
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
					ncli.con, err = setupConnection(c.Bool("skip-cert-verify"))
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
