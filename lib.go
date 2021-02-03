package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"
	nutanix "github.com/routebyintuition/ntnx-go-sdk"
	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

// InputReader provides a reader for secure input and will provide non-secure
type InputReader interface {
	ReadInputSecure() (string, error)
}

// StdInputSecureReader provides a data type for mocking in test the secure input reader
type StdInputSecureReader struct{}

// ReadInputSecure will read in standard input without echo for use with passwords
func (ir StdInputSecureReader) ReadInputSecure() (string, error) {
	pwd, error := terminal.ReadPassword(int(syscall.Stdin))
	return string(pwd), error
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
func GetInputStringValue(ir InputReader, message string, minLen int, def string) (string, error) {
	fmt.Print(message)

	inputString, err := ir.ReadInputSecure()
	if err != nil {
		fmt.Println("error: ", err)
		return def, errors.New("could not read terminal input")
	}

	fmt.Println()

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

func sliceContains(sli []string, str string) bool {
	for _, sliceItem := range sli {
		if sliceItem == str {
			return true
		}
	}

	return false
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

// writeProfileFile writes a profile to the identified file destination
func writeProfileFile(fl string, dl string, pi *profileItem) error {
	yamlData, err := yaml.Marshal(&pi)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := os.Stat(dl); os.IsNotExist(err) {
		err := os.Mkdir(dl, 0760)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(fl, yamlData, 0600)
	if err != nil {
		return err
	}
	return nil
}
