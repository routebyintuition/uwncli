package main

import (
	"errors"
	"strconv"

	"github.com/routebyintuition/ntnx-go-sdk/pe"
	"github.com/urfave/cli/v2"
)

// performs a prism element disk list to v2 API
func (n *NCLI) diskList(c *cli.Context) error {
	ListRequest := new(pe.DiskListRequest)

	getRes, _, err := n.con.PE.Disk.List(ListRequest)
	if err != nil {
		return err
	}

	n.tr.SetHeader([]string{"Disk UUID", "Tier", "Size", "Status", "Host", "Online"})

	n.tr.SetFooter([]string{"", "", "", "", "TOTAL", strconv.Itoa(*getRes.Metadata.TotalEntities)})

	data := [][]string{}

	for _, entityValue := range getRes.Entities {
		data = append(data, []string{*entityValue.DiskUUID, *entityValue.StorageTierName, strconv.Itoa(int(*entityValue.DiskSize)), *entityValue.DiskStatus, *entityValue.HostName, strconv.FormatBool(*entityValue.Online)})
	}
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

// vDiskList lists all vdisks with details
func (n *NCLI) vDiskList(c *cli.Context) error {
	ListRequest := new(pe.DiskVirtualListRequest)

	getRes, _, err := n.con.PE.Disk.ListVDisk(ListRequest)
	if err != nil {
		return err
	}

	n.tr.SetHeader([]string{"vDisk UUID", "Attached", "Disk Capacity", "VM Disk Address", "Storage Container UUID"})

	n.tr.SetFooter([]string{"", "", "", "TOTAL", strconv.Itoa(*getRes.Metadata.TotalEntities)})

	data := [][]string{}

	for _, entityValue := range getRes.Entities {
		attachedVM := "None"
		if entityValue.AttachedVMUUID != nil {
			attachedVM = *entityValue.AttachedVMUUID
		}
		diskVMAddress := "None"
		if entityValue.DiskAddress != nil {
			diskVMAddress = *entityValue.DiskAddress
		}

		data = append(data, []string{*entityValue.UUID, attachedVM, BytesToHumanReadable(*entityValue.DiskCapacityInBytes), diskVMAddress, *entityValue.StorageContainerUUID})
	}
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

// vDiskGetByUUID returns vdisk details based upon vdisk UUID
func (n *NCLI) vDiskGetByUUID(uuid string) (*pe.DiskVirtualGetResponse, error) {
	if !IsValidUUID(uuid) {
		return nil, errors.New("invalid UUID format on vdisk get")
	}

	ListRequest := &pe.DiskVirtualGetRequest{Query: &pe.DiskVirtualGetRequestQuery{UUID: uuid}}

	getRes, _, err := n.con.PE.Disk.GetVDisk(ListRequest)
	if err != nil {
		return nil, err
	}

	return getRes, nil
}
