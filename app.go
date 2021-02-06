package main

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func getFlags() []cli.Flag {
	output := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "pcurl",
			Aliases:     []string{"purl"},
			Usage:       "Prism Central URL",
			DefaultText: "https://10.0.0.1:9440/api/nutanix/v3/",
			EnvVars:     []string{"NUTANIX_PC_URL"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "pcaddress",
			Aliases:     []string{"pca"},
			Usage:       "Prism Central Address",
			DefaultText: "10.0.0.1:9440",
			EnvVars:     []string{"NUTANIX_PC_ADDRESS"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "peurl",
			Aliases:     []string{"element_url"},
			Usage:       "Prism Element URL",
			DefaultText: "https://10.0.0.11:9440/PrismGateway/services/rest/v2.0/",
			EnvVars:     []string{"NUTANIX_PE_URL"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "peaddress",
			Aliases:     []string{"element_address"},
			Usage:       "Prism Element Address (any CVM IP/hostname)",
			DefaultText: "10.0.0.11:9440",
			EnvVars:     []string{"NUTANIX_PE_ADDRESS"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "username",
			Aliases:     []string{"u", "user"},
			Usage:       "Prism Central Username",
			DefaultText: "<username>",
			EnvVars:     []string{"NUTANIX_PC_USER"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "password",
			Aliases:     []string{"p", "pass"},
			Usage:       "Prism Central Passowrd",
			DefaultText: "<password>",
			EnvVars:     []string{"NUTANIX_PC_PASS"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "karbonaddress",
			Aliases:     []string{"ka", "kaddr", "kaddress"},
			Usage:       "Karbon Address <IP or hostname with port>",
			DefaultText: "<karbon address>",
			EnvVars:     []string{"NUTANIX_KARBON_ADDRESS"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "karbonuser",
			Aliases:     []string{"ku", "kuser"},
			Usage:       "Karbon Username",
			DefaultText: "<karbon username>",
			EnvVars:     []string{"NUTANIX_KARBON_USER"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "karbonpass",
			Aliases:     []string{"kp", "kpass"},
			Usage:       "Karbon Passowrd",
			DefaultText: "<karbon password>",
			EnvVars:     []string{"NUTANIX_KARBON_PASS"},
		}),
		&cli.StringFlag{
			Name:        "profile",
			Aliases:     []string{"pro"},
			Value:       "default",
			Usage:       "stored credential profile to use",
			DefaultText: "<default>",
			EnvVars:     []string{"NUTANIX_PROFILE"},
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
	}
	return output
}
