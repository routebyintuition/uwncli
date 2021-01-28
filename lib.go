package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	nutanix "github.com/routebyintuition/ntnx-go-sdk"
	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/crypto/ssh/terminal"
)

// ArgList is a list of values passed in used to determine if they are nil and have a min and max length
type ArgList struct {
	Name      string
	Value     interface{}
	MinLength int
	MaxLength int
}

func evalStringPtrSlice(list []ArgList) error {

	return nil
}

func isStringPtrNotNil(str *string) (bool, int) {
	if str == nil {
		return false, 0
	}

	return true, len(*str)
}

// setupConnection will setup the new prism central SDK connection
func setupConnection(c *cli.Context) (*nutanix.Client, error) {
	httpClient := &http.Client{Transport: &http.Transport{}}

	if c.Bool("skip-cert-verify") {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	pcURL := c.String("url")
	if pcURL == "" {
		pcURL = fmt.Sprintf("https://%s/api/nutanix/v3/", c.String("pcaddress"))
	}

	pcConfig := &pc.ServiceConfig{
		User: nutanix.String(c.String("username")),
		Pass: nutanix.String(c.String("password")),
		URL:  nutanix.String(pcURL),
	}

	con, err := nutanix.NewClient(httpClient, &nutanix.Config{PrismCentral: pcConfig})
	if err != nil {
		fmt.Println("error on NewClient: ", err)
		return nil, err
	}

	return con, err
}

// IsValidUUID validates UUID string
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// GetInputStringValue provides a command prompt to retrieve the user input and allows for default values upon user entry
func GetInputStringValue(message string, minLen int, def string) (string, error) {
	fmt.Print(message)

	inputByte, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("error: ", err)
		return def, errors.New("could not read terminal input")
	}

	fmt.Println()

	inputString := string(inputByte)
	input := strings.TrimSpace(inputString)

	if len(input) < minLen {
		fmt.Printf("invalid input length. less than %d characters \n", minLen)
		return def, errors.New("invalid input length")
	}

	if input == "" {
		return def, nil
	}

	return input, nil
}

// GetConfigLocale returns the full path to credential configuration file
func GetConfigLocale(profile string) (string, string) {
	home, _ := homedir.Dir()
	dirLocale := filepath.Join(home, ".nutanix")
	fileName := fmt.Sprintf("%s.credential", profile)
	fileLocale := filepath.Join(home, ".nutanix", fileName)

	return dirLocale, fileLocale
}

// GetMibFromMB returns the Mebibyte value from a provided Megabyte int
func GetMibFromMB(mb int) int {
	return int(math.Floor(float64(mb) * 0.95367431640625))
}

// defaultInputSource creates a default InputSourceContext.
func defaultInputSource() (altsrc.InputSourceContext, error) {
	return &altsrc.MapInputSource{}, nil
}

// NewYamlSourceFromProfileFunc creates a new Yaml InputSourceContext from a provided flag name and source context.
func NewYamlSourceFromProfileFunc(flagProfileName string) func(context *cli.Context) (altsrc.InputSourceContext, error) {
	return func(context *cli.Context) (altsrc.InputSourceContext, error) {
		if context.IsSet(flagProfileName) {
			profileName := context.String(flagProfileName)
			_, profilePath := GetConfigLocale(profileName)
			return altsrc.NewYamlSourceFromFile(profilePath)
		}

		_, profilePath := GetConfigLocale("default")
		if _, err := os.Stat(profilePath); err == nil {
			return altsrc.NewYamlSourceFromFile(profilePath)
		}

		return defaultInputSource()
	}
}
