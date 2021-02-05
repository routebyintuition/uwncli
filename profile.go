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
	KarbonUser    string
	KarbonPass    string
	Username      string
	Password      string
}

func (b *BCLI) configureDefaultProfile(c *cli.Context) error {

	ir := &StdInputSecureReader{}
	sr := &StdInputReader{}

	dirLocale, fileLocale := GetConfigLocale("default")

	if _, err := os.Stat(fileLocale); err == nil {
		errStr := fmt.Sprintf("default profile already created. first run: uwncli profile delete default")
		return errors.New(errStr)
	}

	var prismCenUser, prismCenPass, prismCenAddr, karbonAddr, karbonUser, karbonPass string

	prismCenUser, err := GetInputStringValue(sr, "prism central username (ex: admin): ", 3, "")
	if err != nil {
		return err
	}
	prismCenPass, err = GetInputStringValue(ir, "prism central password: ", 8, "")
	if err != nil {
		return err
	}
	prismCenAddr, err = GetInputStringValue(sr, "prism central address (ex: 10.0.0.1:9440): ", 6, "")
	if err != nil {
		return err
	}

	karbonAddr, err = GetInputStringValue(sr, fmt.Sprintf("nutanix karbon address [%s]: ", prismCenAddr), 6, prismCenAddr)
	if err != nil {
		return err
	}

	karbonUser, err = GetInputStringValue(sr, fmt.Sprintf("nutanix karbon username [%s]: ", prismCenUser), 3, prismCenUser)
	if err != nil {
		return err
	}

	karbonPass, err = GetInputStringValue(ir, "nutanix karbon password [pc password entered above by default]: ", 8, prismCenPass)
	if err != nil {
		return err
	}

	profileItem := &profileItem{}

	profileItem.PCAddress = prismCenAddr
	profileItem.Username = prismCenUser
	profileItem.Password = prismCenPass
	profileItem.KarbonAddress = karbonAddr
	profileItem.KarbonUser = karbonUser
	profileItem.KarbonPass = karbonPass

	err = writeProfileFile(fileLocale, dirLocale, profileItem)
	if err != nil {
		return err
	}

	fmt.Println("saved profile to: ", fileLocale)

	return nil
}

func (b *BCLI) createProfile(c *cli.Context) error {
	ir := &StdInputSecureReader{}
	sr := &StdInputReader{}

	profileName, err := GetInputStringValue(sr, "Profile name [default]: ", 0, "default")
	if err != nil {
		return err
	}

	dirLocale, fileLocale := GetConfigLocale(profileName)

	if _, err := os.Stat(fileLocale); err == nil {
		errStr := fmt.Sprintf("%s profile already created. first run: uwncli profile delete %s", profileName, profileName)
		return errors.New(errStr)
	}

	var prismCenUser, prismCenPass, prismCenAddr, karbonAddr, karbonUser, karbonPass string

	prismCenUser, err = GetInputStringValue(sr, "prism central username (ex: admin): ", 3, "")
	if err != nil {
		return err
	}
	prismCenPass, err = GetInputStringValue(ir, "prism central password: ", 8, "")
	if err != nil {
		return err
	}
	prismCenAddr, err = GetInputStringValue(sr, "prism central address (ex: 10.0.0.1:9440): ", 6, "")
	if err != nil {
		return err
	}

	karbonAddr, err = GetInputStringValue(sr, fmt.Sprintf("nutanix karbon address [%s]: ", prismCenAddr), 6, prismCenAddr)
	if err != nil {
		return err
	}

	karbonUser, err = GetInputStringValue(sr, fmt.Sprintf("nutanix karbon username [%s]: ", prismCenUser), 3, prismCenUser)
	if err != nil {
		return err
	}

	karbonPass, err = GetInputStringValue(ir, "nutanix karbon password [pc password entered above by default]: ", 8, prismCenPass)
	if err != nil {
		return err
	}

	profileItem := &profileItem{}

	profileItem.PCAddress = prismCenAddr
	profileItem.Username = prismCenUser
	profileItem.Password = prismCenPass
	profileItem.KarbonAddress = karbonAddr
	profileItem.KarbonUser = karbonUser
	profileItem.KarbonPass = karbonPass

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
