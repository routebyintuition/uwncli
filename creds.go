package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
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

type profileList map[string]profileItem

func saveCredentials(c *cli.Context) error {

	profile := make(profileList)

	profileName, err := GetInputStringValue("Profile name [default]: ", 0, "default")
	if err != nil {
		return err
	}

	dirLocale, fileLocale := GetConfigLocale(profileName)

	if _, err := os.Stat(fileLocale); err == nil {
		confFile, err := ioutil.ReadFile(fileLocale)
		if err != nil {
			log.Printf("could not open config file: %v", err)
		}
		err = yaml.Unmarshal(confFile, profile)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}

		fmt.Println("unmarshalled existing config: ", profile)
	}

	var prismCenUser, prismCenPass, prismCenAddr string

	prismCenUser, err = GetInputStringValue("prism central username [admin]: ", 0, "admin")
	if err != nil {
		return err
	}
	prismCenPass, err = GetInputStringValue("prism central password []: ", 8, "")
	if err != nil {
		return err
	}
	prismCenAddr, err = GetInputStringValue("prism central address [10.0.0.1:9440]: ", 6, "10.0.0.1")

	profileItem := &profileItem{}

	profileItem.PCAddress = prismCenAddr
	profileItem.Username = prismCenUser
	profileItem.Password = prismCenPass

	yamlData, err := yaml.Marshal(&profileItem)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := os.Stat(dirLocale); os.IsNotExist(err) {
		err := os.Mkdir(dirLocale, 0760)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(fileLocale, yamlData, 0600)
	if err != nil {
		return err
	}

	fmt.Println("saved profile to: ", fileLocale)

	return nil
}
