package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

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
		}
	}

	n.tr.SetFooter([]string{"Total", "", strconv.Itoa(profileCount)})
	n.tr.SetAutoWrapText(false)
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

func (n *NCLI) deleteProfile(c *cli.Context) error {
	if len(c.Args().First()) < 2 {
		return errors.New("must enter a valid profile name")
	}

	_, fileLocale := GetConfigLocale(c.Args().First())

	err := os.Remove(fileLocale)
	if err != nil {
		return err
	}
	fmt.Println("deleted profile: ", c.Args().First())

	return nil
}
