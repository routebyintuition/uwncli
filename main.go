package main

import (
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

// BCLI (base CLI) is used for non-API calls but allows the table writer setup
type BCLI struct {
	tr *tablewriter.Table
}

func main() {

	ncli := &NCLI{}
	ncli.tr = tablewriter.NewWriter(os.Stdout)

	bcli := &BCLI{}
	bcli.tr = tablewriter.NewWriter(os.Stdout)

	flags := getFlags()

	app := &cli.App{
		Name:                 "Unikum und Wunderbar Nutanix CLI",
		Usage:                "uwncli [flags] [command] [subcommand]",
		EnableBashCompletion: true,
		Before:               altsrc.InitInputSourceWithContext(flags, NewYamlSourceFromProfileFunc("profile")),
		Flags:                flags,
		Commands: []*cli.Command{
			{
				Name:    "configure",
				Aliases: []string{"conf"},
				Usage:   "configure stored credentials",
				Action:  bcli.configureDefaultProfile,
			},
			{
				Name:  "profile",
				Usage: "stored profile specific commands. use `uwncli profile help` to view options",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list all stored profiles",
						Action: bcli.listProfiles,
					},
					{
						Name:   "delete",
						Usage:  "uwncli delete <profile name>",
						Action: bcli.deleteProfile,
					},
					{
						Name:   "create",
						Usage:  "create a new profile",
						Action: bcli.createProfile,
					},
				},
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name:  "vm",
				Usage: "virtual machine specific commands. use `uwncli vm help` to view options",
				Subcommands: []*cli.Command{
					{
						Name:     "list",
						Usage:    "retrieve all VMs",
						Action:   ncli.vmList,
						Category: "get",
					},
					{
						Name:     "get",
						Usage:    "<VM UUID>",
						Action:   ncli.vmGet,
						Category: "get",
					},
					{
						Name:     "disklist",
						Usage:    "<VM UUID>",
						Action:   ncli.vmDiskList,
						Category: "put",
					},
					{
						Name:     "update-memory",
						Usage:    "<VM UUID> <memory in MB integer>",
						Action:   ncli.vmMemoryUpdate,
						Category: "put",
					},
					{
						Name:     "update-power",
						Usage:    "<VM UUID> <ON|OFF|POWERCYCLE|RESET|PAUSE|SUSPEND|RESUME|ACPI_SHUTDOWN|ACPI_REBOOT>",
						Action:   ncli.vmSetPowerState,
						Category: "put",
					},
				},
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name:  "image",
				Usage: "image specific commands. use `uwncli image help` to view options",
				Subcommands: []*cli.Command{
					{
						Name:     "list",
						Usage:    "list all images",
						Action:   ncli.imageList,
						Category: "image",
					},
					{
						Name:     "create",
						Usage:    "create a new image",
						Action:   ncli.imageCreate,
						Category: "image",
					},
				},
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name:  "cluster",
				Usage: "cluster specific commands. use `uwncli cluster help` to view options",
				Subcommands: []*cli.Command{
					{
						Name:     "all",
						Usage:    "retrieve all VMs",
						Category: "cluster",
					},
				},
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name:  "subnet",
				Usage: "subnet specific commands. use `uwncli subnet help` to view options",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list all subnets",
						Action: ncli.listSubnets,
					},
				},
			},
			{
				Before: func(c *cli.Context) error {
					var err error
					ncli.con, err = setupConnection(c)
					return err
				},
				Name:  "karbon",
				Usage: "karbon specific commands. use `uwncli karbon help` to view options",
				Subcommands: []*cli.Command{
					{
						Name:  "cluster",
						Usage: "karbon cluster commands. use `uwncli karbon cluster help",
						Subcommands: []*cli.Command{
							{
								Name:   "list",
								Usage:  "list all karbon clusters",
								Action: ncli.krbnListClusters,
							},
						},
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
