package main

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"
	nutanix "github.com/routebyintuition/ntnx-go-sdk"
	"github.com/routebyintuition/ntnx-go-sdk/karbon"
	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/routebyintuition/ntnx-go-sdk/pe"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

// InputReader provides a reader for secure input and will provide non-secure
type InputReader interface {
	ReadInput() (string, error)
}

// StdInputSecureReader provides a data type for mocking in test the secure input reader
type StdInputSecureReader struct{}

// StdInputReader data type for mocking that should not be used for passwords
type StdInputReader struct{}

// ReadInput will read in standard input without echo for use with passwords
func (ir StdInputSecureReader) ReadInput() (string, error) {
	pwd, error := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(pwd), error
}

// ReadInput will read in standard input but should not be used for passwords
func (ir StdInputReader) ReadInput() (string, error) {
	var input string
	fmt.Scanln(&input)

	return input, nil
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

	peURL := c.String("peurl")
	if peURL == "" {
		peURL = fmt.Sprintf("https://%s/PrismGateway/services/rest/v2.0/", c.String("peaddress"))
	}

	karbonURL := c.String("karbonurl")
	if karbonURL == "" {
		karbonURL = fmt.Sprintf("https://%s/karbon/", c.String("karbonaddress"))
	}

	pcConfig := &pc.ServiceConfig{
		User: nutanix.String(c.String("username")),
		Pass: nutanix.String(c.String("password")),
		URL:  nutanix.String(pcURL),
	}

	peConfig := &pe.ServiceConfig{
		User: nutanix.String(c.String("username")),
		Pass: nutanix.String(c.String("password")),
		URL:  nutanix.String(peURL),
	}

	krbnConfig := &karbon.ServiceConfig{
		User: nutanix.String(c.String("karbonuser")),
		Pass: nutanix.String(c.String("karbonpass")),
		URL:  nutanix.String(karbonURL),
	}

	con, err := nutanix.NewClient(httpClient, &nutanix.Config{
		PrismCentral: pcConfig,
		Karbon:       krbnConfig,
		PrismElement: peConfig,
	})
	if err != nil {
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

	inputString, err := ir.ReadInput()
	if err != nil {
		return def, errors.New("could not read terminal input")
	}

	input := strings.TrimSpace(inputString)

	if len(input) < minLen && len(def) < minLen {
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

// BytesToHumanReadable converts byte sizes to human readable
func BytesToHumanReadable(size int64) string {
	sizeInt := float64(size)

	sizeKB := fmt.Sprintf("%.1f", float64(sizeInt/1000))
	sizeMB := fmt.Sprintf("%.1f", float64(sizeInt/1000000))
	sizeGB := fmt.Sprintf("%.1f", float64(sizeInt/1000000000))
	sizeTB := fmt.Sprintf("%.1f", float64(sizeInt/1000000000000))

	diskSizeStr := strconv.Itoa(int(size)) + " Bytes"

	if sizeKB != "0.0" {
		diskSizeStr = fmt.Sprintf("%s KB", sizeKB)
	}

	if sizeMB != "0.0" {
		diskSizeStr = fmt.Sprintf("%s KB", sizeMB)
	}

	if sizeGB != "0.0" {
		diskSizeStr = fmt.Sprintf("%s GB", sizeGB)
	}

	if sizeTB != "0.0" {
		diskSizeStr = fmt.Sprintf("%s TB", sizeTB)
	}

	return diskSizeStr
}

// isInputFromPipe determines if the input is from a pipe or file
// thanks to dev.to article
func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func processYAMLReader(r io.Reader, i interface{}) (interface{}, error) {
	err := yaml.NewDecoder(r).Decode(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

// isDirectory checks to see if path is a directory
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// ValidYAMLCreate validates a create VM YAML file
func (n *NCLI) ValidYAMLCreate(pcc *pc.VMCreateRequest) error {
	if len(*pcc.Spec.Name) < 3 {
		return errors.New("invalid length of VM name")
	}

	subnetExists := false
	subnetList, err := n.getSubnetUUIDList()

	if err != nil {
		return err
	}
	for _, v := range *pcc.Spec.Resources.NicList {
		if stringSliceContains(subnetList, v.SubnetReference.UUID) {
			subnetExists = true
		}
	}
	if !subnetExists {
		return errors.New("subnet UUID not defined on cluster")
	}

	imageExists := false
	imageList, err := n.GetImageUUIDList()
	if err != nil {
		return err
	}
	for _, v := range *pcc.Spec.Resources.DiskList {
		if v.DataSourceReference != nil {
			if stringSliceContains(imageList, v.DataSourceReference.UUID) {
				imageExists = true
			}
		} else {
			imageExists = true
		}
	}
	if !imageExists {
		return errors.New("disk image UUID not defined on cluster")
	}

	return errors.New("function not ready for prime time")
}

// stringSliceContains checks whether a string slice contains a specific string
func stringSliceContains(slice []string, str string) bool {
	for _, val := range slice {
		if val == str {
			return true
		}
	}
	return false
}

// isBase64 determines of a string is base64 encoded...used for cloud_init
func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
