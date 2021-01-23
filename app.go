package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func getFlags() []cli.Flag {
	output := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "pcurl",
			Aliases: []string{"purl"},
			// Value:       "https://localhost:9440/api/nutanix/v3/",
			Usage:       "Prism Central URL",
			DefaultText: "https://10.0.0.1:9440/api/nutanix/v3/",
			EnvVars:     []string{"NUTANIX_PC_URL"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "pcaddress",
			Aliases: []string{"pca"},
			// Value:       "10.0.0.1",
			Usage:       "Prism Central Address",
			DefaultText: "10.0.0.1:9440",
			EnvVars:     []string{"NUTANIX_PC_ADDRESS"},
		}),
		&cli.StringFlag{
			Name:  "load-profile",
			Usage: "<filename>",
			// Value:       fileLocale,
			// DefaultText: fileLocale,
			EnvVars: []string{"NUTANIX_LOAD_PROFILE"},
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u", "user"},
			// Value:       "{pc username}",
			Usage:       "Prism Central Username",
			DefaultText: "<username>",
			EnvVars:     []string{"NUTANIX_PC_USER"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p", "pass"},
		}),
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"conf", "co"},
			Value:       "~/.nutanix/credentials",
			Usage:       "Credentials file location",
			DefaultText: "<credentials file and path>",
		},
		&cli.StringFlag{
			Name:        "profile",
			Aliases:     []string{"pro"},
			Value:       "default",
			Usage:       "stored credential profile to use",
			DefaultText: "<default>",
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

func (n *NCLI) listProfiles(c *cli.Context) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	profileFolder := fmt.Sprintf("%s/.nutanix/", home)
	profileFileList, err := ioutil.ReadDir(profileFolder)
	if err != nil {
		return err
	}

	n.tr.SetHeader([]string{"Profile", "Last Modified", "Location"})

	data := [][]string{}
	profileCount := 0

	for _, profileFileItem := range profileFileList {
		if !profileFileItem.IsDir() && strings.HasSuffix(profileFileItem.Name(), ".credential") {
			profileCount++
			data = append(data, []string{strings.TrimSuffix(profileFileItem.Name(), ".credential"), profileFileItem.ModTime().String(), fmt.Sprintf("%s%s", profileFolder, profileFileItem.Name())})
			fmt.Println("profile: ", strings.TrimSuffix(profileFileItem.Name(), ".credential"))
			fmt.Println("modified: ", profileFileItem.ModTime())
		}
	}

	n.tr.SetFooter([]string{"Total", "", strconv.Itoa(profileCount)})
	n.tr.SetAutoWrapText(false)
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}
