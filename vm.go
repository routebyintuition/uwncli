package main

import (
	"errors"
	"strconv"

	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/urfave/cli/v2"
)

func (n *NCLI) vmList(c *cli.Context) error {

	ListRequest := new(pc.VMListRequest)
	ListRequest.Length = 1000

	getRes, _, err := n.con.PC.VM.List(ListRequest)
	if err != nil {
		return err
	}

	n.tr.SetHeader([]string{"Name", "UUID"})

	n.tr.SetFooter([]string{"Total", strconv.Itoa(len(getRes.Entities))})

	data := [][]string{}

	for _, entityValue := range getRes.Entities {
		data = append(data, []string{*entityValue.Spec.Name, *entityValue.Metadata.UUID})
	}
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

func (n *NCLI) vmGet(c *cli.Context) error {
	if !IsValidUUID(c.Args().First()) {
		return errors.New("invalid UUID format")
	}

	getRequest := &pc.VMGetRequest{UUID: c.Args().First()}
	getRes, _, err := n.con.PC.VM.Get(getRequest)
	if err != nil {
		return err
	}

	data := [][]string{}

	name := *getRes.Spec.Name

	for _, diskItem := range *getRes.Status.Resources.DiskList {
		data = append(data, []string{name, diskItem.UUID, strconv.Itoa(diskItem.DiskSizeBytes), diskItem.DeviceProperties.DiskAddress.AdapterType})
	}

	n.tr.SetHeader([]string{"VM", "Disk UUID", "Size (Bytes)", "Disk Type"})

	n.tr.SetFooter([]string{"", "", "Total", strconv.Itoa(len(*getRes.Status.Resources.DiskList))})

	n.tr.SetAutoMergeCells(true)
	n.tr.SetRowLine(true)
	n.tr.AppendBulk(data)
	n.tr.Render()
	return nil
}

func (n *NCLI) vmDiskList(c *cli.Context) error {
	if !IsValidUUID(c.Args().First()) {
		return errors.New("invalid UUID format")
	}

	getRequest := &pc.VMGetRequest{UUID: c.Args().First()}
	getRes, _, err := n.con.PC.VM.Get(getRequest)

	if err != nil {
		return err
	}

	data := [][]string{}

	name := *getRes.Spec.Name

	for _, diskItem := range *getRes.Status.Resources.DiskList {
		data = append(data, []string{diskItem.UUID, strconv.Itoa(diskItem.DiskSizeBytes), diskItem.DeviceProperties.DiskAddress.AdapterType})
	}

	n.tr.SetHeader([]string{"Disk UUID", "Size (Bytes)", "Disk Type"})

	n.tr.SetFooter([]string{name, "Total", strconv.Itoa(len(*getRes.Status.Resources.DiskList))})

	//n.tr.SetAutoMergeCells(true)
	//n.tr.SetRowLine(true)
	n.tr.AppendBulk(data)
	n.tr.Render()
	return nil
}
