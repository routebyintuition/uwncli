package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"

	nutanix "github.com/routebyintuition/ntnx-go-sdk"
	"github.com/routebyintuition/ntnx-go-sdk/pc"
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
func setupConnection(skipVerify bool) (*nutanix.Client, error) {
	httpClient := &http.Client{Transport: &http.Transport{}}

	if skipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	con, err := nutanix.NewClient(httpClient, &nutanix.Config{PrismCentral: new(pc.ServiceConfig)})
	if err != nil {
		fmt.Println("error on NewClient: ", err)
		return nil, err
	}

	return con, err
}

func vmList(PC *pc.Client, vmName string, count int) {
	ListRequest := new(pc.VMListRequest)
	ListRequest.Length = count

	getRes, _, err := PC.VM.List(ListRequest)
	if err != nil {
		fmt.Println("cluster list error: ", err)
		return
	}

	fmt.Println("VM Count: ", len(getRes.Entities))
	for index, entityValue := range getRes.Entities {
		fmt.Printf("VM: %d - %s \n", index, entityValue.Spec.Name)
		fmt.Printf("\t UUID: %s \n", entityValue.Metadata.UUID)
	}
}

// IsValidUUID validates UUID string
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
