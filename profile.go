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

type profileItem struct {
	PCAddress     string
	PCURL         string
	PEAddress     string
	PEURL         string
	KarbonAddress string
	KarbonURL     string
	Username      string
	Password      string
}

func (b *BCLI) configureDefaultProfile(c *cli.Context) error {

	ir := &StdInputSecureReader{}

	dirLocale, fileLocale := GetConfigLocale("default")

	if _, err := os.Stat(fileLocale); err == nil {
		errStr := fmt.Sprintf("default profile already created. first run: uwncli profile delete default")
		return errors.New(errStr)
	}

	var prismCenUser, prismCenPass, prismCenAddr string

	prismCenUser, err := GetInputStringValue(ir, "prism central username [admin]: ", 0, "admin")
	if err != nil {
		return err
	}
	prismCenPass, err = GetInputStringValue(ir, "prism central password []: ", 8, "")
	if err != nil {
		return err
	}
	prismCenAddr, err = GetInputStringValue(ir, "prism central address [10.0.0.1:9440]: ", 6, "10.0.0.1")

	profileItem := &profileItem{}

	profileItem.PCAddress = prismCenAddr
	profileItem.Username = prismCenUser
	profileItem.Password = prismCenPass

	err = writeProfileFile(fileLocale, dirLocale, profileItem)
	if err != nil {
		return err
	}

	fmt.Println("saved profile to: ", fileLocale)

	return nil
}

func (b *BCLI) createProfile(c *cli.Context) error {
	ir := &StdInputSecureReader{}

	profileName, err := GetInputStringValue(ir, "Profile name [default]: ", 0, "default")
	if err != nil {
		return err
	}

	dirLocale, fileLocale := GetConfigLocale(profileName)

	if _, err := os.Stat(fileLocale); err == nil {
		errStr := fmt.Sprintf("%s profile already created. first run: uwncli profile delete %s", profileName, profileName)
		return errors.New(errStr)
	}

	var prismCenUser, prismCenPass, prismCenAddr string

	prismCenUser, err = GetInputStringValue(ir, "prism central username [admin]: ", 0, "admin")
	if err != nil {
		return err
	}
	prismCenPass, err = GetInputStringValue(ir, "prism central password []: ", 8, "")
	if err != nil {
		return err
	}
	prismCenAddr, err = GetInputStringValue(ir, "prism central address [10.0.0.1:9440]: ", 6, "10.0.0.1")

	profileItem := &profileItem{}

	profileItem.PCAddress = prismCenAddr
	profileItem.Username = prismCenUser
	profileItem.Password = prismCenPass

	err = writeProfileFile(fileLocale, dirLocale, profileItem)
	if err != nil {
		return err
	}

	fmt.Println("saved profile to: ", fileLocale)

	return nil
}

func (b *BCLI) listProfiles(c *cli.Context) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	profileFolder := fmt.Sprintf("%s/.nutanix/", home)
	profileFileList, err := ioutil.ReadDir(profileFolder)
	if err != nil {
		return err
	}

	b.tr.SetHeader([]string{"Profile", "Last Modified", "Location"})

	data := [][]string{}
	profileCount := 0

	for _, profileFileItem := range profileFileList {
		if !profileFileItem.IsDir() && strings.HasSuffix(profileFileItem.Name(), ".credential") {
			profileCount++
			data = append(data, []string{strings.TrimSuffix(profileFileItem.Name(), ".credential"), profileFileItem.ModTime().String(), fmt.Sprintf("%s%s", profileFolder, profileFileItem.Name())})
		}
	}

	b.tr.SetFooter([]string{"Total", "", strconv.Itoa(profileCount)})
	b.tr.SetAutoWrapText(false)
	b.tr.AppendBulk(data)
	b.tr.Render()

	return nil
}

func (b *BCLI) deleteProfile(c *cli.Context) error {
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
