package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/urfave/cli/v2"
)

func (n *NCLI) imageList(c *cli.Context) error {

	ListRequest := new(pc.ImageListRequest)
	ListRequest.Length = 1000

	getRes, _, err := n.con.PC.Image.List(ListRequest)
	if err != nil {
		return err
	}

	n.tr.SetHeader([]string{"Name", "Type", "UUID", "Status", "Size", "Source"})

	n.tr.SetFooter([]string{"", "", "", "", "Total", strconv.Itoa(len(getRes.Entities))})

	data := [][]string{}

	for _, entityValue := range getRes.Entities {
		var eType, eName string

		if entityValue.Spec.Name != nil {
			eName = fmt.Sprintf(*entityValue.Spec.Name)
		} else {
			eName = ""
		}
		if entityValue.Spec.Resources.ImageType != nil {
			eType = fmt.Sprintf("%v", *entityValue.Spec.Resources.ImageType)
		} else {
			eType = ""
		}

		eUUID := fmt.Sprintf(*entityValue.Metadata.UUID)
		eState := fmt.Sprintf(entityValue.Status.State)

		data = append(data, []string{eName, eType, eUUID, eState, strconv.Itoa(*entityValue.Status.Resources.SizeBytes), *entityValue.Status.Resources.SourceURI})
	}
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

func (n *NCLI) imageCreate(c *cli.Context) error {

	CreateRequest := pc.ImageCreateRequest{}

	iName := c.String("image-name")
	if len(iName) < 3 {
		return errors.New("image-name is undefined or less than 3 characters")
	}

	iDesc := c.String("image-description")
	if len(iDesc) < 3 {
		return errors.New("image-description is undefined or less than 3 characters")
	}

	iType := c.String("image-type")
	if iType != "ISO_IMAGE" && iType != "DISK_IMAGE" {
		return errors.New("image-type must be either DISK_IMAGE or ISO_IMAGE")
	}

	iSource := c.String("image-source")
	if len(iSource) < 3 {
		return errors.New("image-source is undefined or less than 3 characters")
	}

	iKind := "image"
	iCatMap := false

	spec := pc.Spec{
		Name:        &iName,
		Description: &iDesc,
	}
	res := pc.Resources{
		ImageType: &iType,
		SourceURI: &iSource,
	}
	meta := pc.Metadata{
		Kind:                 &iKind,
		UseCategoriesMapping: &iCatMap,
	}

	CreateRequest.Spec = &spec
	CreateRequest.Spec.Resources = &res
	CreateRequest.Metadata = &meta

	getRes, _, err := n.con.PC.Image.Create(&CreateRequest)
	if err != nil {
		return err
	}

	n.tr.SetHeader([]string{"Name", "UUID", "Description", "Status"})

	out, _ := json.MarshalIndent(getRes, "", "  ")
	fmt.Println("cluster list result: ", string(out))

	data := [][]string{}
	data = append(data, []string{*getRes.Spec.Name, *getRes.Metadata.UUID, *getRes.Spec.Description, getRes.Status.State})
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

// GetImageUUIDList returns a string slice containing all image UUIDs
func (n *NCLI) GetImageUUIDList() ([]string, error) {

	ListRequest := new(pc.ImageListRequest)
	ListRequest.Length = 40

	data := []string{}

	var listLoop []pc.Entities
	totalMatches := 0
	offset := 0
	currentMatches := -1
	var err error
	getRes := &pc.ImageListResponse{}

	for totalMatches > currentMatches {
		if currentMatches == -1 {
			currentMatches = 0
		}

		ListRequest.Offset = offset

		getRes, _, err = n.con.PC.Image.List(ListRequest)
		if err != nil {
			return nil, err
		}

		currentMatches += *getRes.Metadata.Length
		totalMatches = *getRes.Metadata.TotalMatches
		offset += ListRequest.Length
		listLoop = append(listLoop, getRes.Entities...)
	}

	for _, entityValue := range listLoop {
		if entityValue.Metadata.UUID != nil {
			data = append(data, fmt.Sprintf(*entityValue.Metadata.UUID))
		}
	}

	return data, nil
}
