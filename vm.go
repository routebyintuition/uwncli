package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	nutanix "github.com/routebyintuition/ntnx-go-sdk"
	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/routebyintuition/ntnx-go-sdk/pe"
	"github.com/urfave/cli/v2"
)

func (n *NCLI) vmCreate(c *cli.Context) error {
	createInstruct := &pc.VMCreateRequest{}

	if isInputFromPipe() {
		createInstructStruct, err := processYAMLReader(os.Stdin, createInstruct)
		if err != nil {
			return err
		}
		outJSON, _ := json.MarshalIndent(createInstructStruct, "", "  ")
		fmt.Println("YAML process results: ", string(outJSON))

	} else if len(c.String("vm-yaml")) > 0 {
		dir, err := isDirectory(c.String("vm-yaml"))
		if err != nil {
			return err
		}
		if dir {
			return errors.New("path provided is a directory...single yaml file needed")
		}
		fh, err := os.Open(c.String("vm-yaml"))
		if err != nil {
			return err
		}
		fhr := bufio.NewReader(fh)
		createInstructStruct, err := processYAMLReader(fhr, createInstruct)
		outJSON, _ := json.MarshalIndent(createInstructStruct, "", "  ")
		fmt.Println("YAML process results: ", string(outJSON))

	} else {
		return errors.New("no yaml config provided via stdin or file")
	}

	err := n.ValidYAMLCreate(createInstruct)
	if err != nil {
		return err
	}

	getRes, _, err := n.con.PC.VM.Create(createInstruct)
	if err != nil {
		return err
	}

	outJSON, _ := json.MarshalIndent(getRes, "", "  ")
	fmt.Println("YAML process results: ", string(outJSON))

	return nil
}

func (n *NCLI) vmList(c *cli.Context) error {

	ListRequest := new(pc.VMListRequest)
	ListRequest.Length = 40

	var vmListLoop []pc.Entities
	totalMatches := 0
	offset := 0
	currentMatches := -1
	var err error
	var getRes *pc.VMListResponse

	for totalMatches > currentMatches {
		if currentMatches == -1 {
			currentMatches = 0
		}

		ListRequest.Offset = offset

		getRes, _, err = n.con.PC.VM.List(ListRequest)
		if err != nil {
			return err
		}

		currentMatches += *getRes.Metadata.Length
		totalMatches = *getRes.Metadata.TotalMatches
		offset += ListRequest.Length
		vmListLoop = append(vmListLoop, getRes.Entities...)
	}

	n.tr.SetHeader([]string{"Name", "UUID", "Powered", "Cluster"})

	n.tr.SetFooter([]string{"Total", "", "", strconv.Itoa(*getRes.Metadata.TotalMatches)})

	data := [][]string{}

	for _, entityValue := range vmListLoop {
		data = append(data, []string{*entityValue.Spec.Name, *entityValue.Metadata.UUID, *entityValue.Spec.Resources.PowerState, entityValue.Spec.ClusterReference.Name})
	}
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

func (n *NCLI) vmMemoryUpdate(c *cli.Context) error {
	if !IsValidUUID(c.Args().First()) {
		return errors.New("invalid UUID format")
	}

	if len(c.Args().Get(1)) == 0 {
		return errors.New("no memory value provided...should be an integer")
	}

	memVal, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return err
	}
	if memVal < 500 || memVal > 500000 {
		return errors.New("invalid memory value")
	}

	getRequest := &pc.VMGetRequest{UUID: c.Args().First()}
	getRes, _, err := n.con.PC.VM.Get(getRequest)
	if err != nil {
		return err
	}

	updateRequest := &pc.VMUpdateRequest{}
	updateRequestData := &pc.VMUpdateRequestData{}
	updateRequestData.Spec = getRes.Spec
	updateRequestData.APIVersion = &getRes.APIVersion
	updateRequestData.Metadata = &getRes.Metadata

	updateRequest.UUID = c.Args().First()
	updateRequest.Data = *updateRequestData

	mibValue := GetMibFromMB(memVal)
	updateRequest.Data.Spec.Resources.MemorySizeMib = &mibValue

	_, _, err = n.con.PC.VM.Update(updateRequest)
	if err != nil {
		return err
	}

	fmt.Printf("vm updated to %d memory\n", memVal)

	return nil
}

// vmVDiskGet returns the VDISK list for an identified VM
func (n *NCLI) vmVDiskGet(c *cli.Context) error {
	if !IsValidUUID(c.Args().First()) {
		return errors.New("invalid UUID format")
	}

	getRequest := &pe.VMGetRequest{
		Params: &pe.VMGetRequestParams{
			UUID: c.Args().First(),
		},
		Query: &pe.VMGetRequestQuery{
			IncludeVMDiskConfig: true,
		},
	}
	getRes, _, err := n.con.PE.VM.Get(getRequest)
	if err != nil {
		return err
	}

	data := [][]string{}

	name := *getRes.Name

	vdiskCount := 0
	if getRes.VMDiskInfo != nil {
		vdiskCount = len(*getRes.VMDiskInfo)
	}

	if getRes.VMDiskInfo != nil {
		for _, diskItem := range *getRes.VMDiskInfo {
			diskSizeStr := BytesToHumanReadable(int64(*diskItem.Size))

			nfsLocation := "nfs"

			if diskItem.DiskAddress.NdfsFilepath != nil {
				nfsLocation = *diskItem.DiskAddress.NdfsFilepath
			}

			diskType := "DISK"
			if diskItem.IsCdrom != nil {
				if *diskItem.IsCdrom {
					diskType = "CDROM"
				}
			}

			data = append(data, []string{"PC", *diskItem.DiskAddress.DeviceUUID, fmt.Sprintf("nfs://127.0.0.1%s", nfsLocation), diskSizeStr, diskType})
			data = append(data, []string{"PE", *diskItem.DiskAddress.VmdiskUUID, fmt.Sprintf("nfs://127.0.0.1%s", nfsLocation), diskSizeStr, diskType})
		}
	}

	n.tr.SetHeader([]string{"", fmt.Sprintf("Disk UUID - %s", name), "NFS LOCATION", "Size", "Disk Type"})

	data = append(data, []string{"", "", "", "TOTAL", strconv.Itoa(vdiskCount)})

	n.tr.SetAutoMergeCells(true)
	n.tr.SetRowLine(true)
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
		sizeKB := fmt.Sprintf("%.1f", float64(diskItem.DiskSizeBytes/1000))
		sizeMB := fmt.Sprintf("%.1f", float64(diskItem.DiskSizeBytes/1000000))
		sizeGB := fmt.Sprintf("%.1f", float64(diskItem.DiskSizeBytes/1000000000))

		diskSizeStr := strconv.Itoa(diskItem.DiskSizeBytes)

		if sizeKB != "0.0" {
			diskSizeStr = fmt.Sprintf("%s KB", sizeKB)
		}

		if sizeMB != "0.0" {
			diskSizeStr = fmt.Sprintf("%s KB", sizeMB)
		}

		if sizeGB != "0.0" {
			diskSizeStr = fmt.Sprintf("%s GB", sizeGB)
		}

		data = append(data, []string{name, diskItem.UUID, diskSizeStr, diskItem.DeviceProperties.DeviceType})
	}

	n.tr.SetHeader([]string{"VM", "Disk UUID", "Size", "Disk Type"})

	// n.tr.SetFooter([]string{"", "", "Total", strconv.Itoa(len(*getRes.Status.Resources.DiskList))})
	data = append(data, []string{name, "", "TOTAL", strconv.Itoa(len(*getRes.Status.Resources.DiskList))})

	data = append(data, []string{name, "", "POWER STATE", *getRes.Spec.Resources.PowerState})
	data = append(data, []string{name, "NETWORK UUID", "ADDRESS", "NETWORK"})

	for _, networkItem := range *getRes.Status.Resources.NicList {
		networkIP := "UNDEFINED"
		if len(networkItem.IPEndpointList) > 0 {
			networkIP = networkItem.IPEndpointList[0].IP
		}
		data = append(data, []string{name, networkItem.UUID, networkIP, networkItem.SubnetReference.Name})
	}

	n.tr.SetAutoMergeCells(true)
	n.tr.SetRowLine(true)
	n.tr.AppendBulk(data)
	n.tr.Render()
	return nil
}

func (n *NCLI) vmSetPowerState(c *cli.Context) error {
	if !IsValidUUID(c.Args().First()) {
		return errors.New("invalid UUID format")
	}
	powerState := strings.ToUpper(c.Args().Get(1))
	// powerOptions := []string{"ON", "OFF", "POWERCYCLE", "RESET", "PAUSE", "SUSPEND", "RESUME", "ACPI_SHUTDOWN", "ACPI_REBOOT"}
	powerOptions := []string{"ON", "OFF"}
	if !stringSliceContains(powerOptions, powerState) {
		return errors.New("invalid memory state. <ON, OFF>")
	}

	getRequest := &pc.VMGetRequest{UUID: c.Args().First()}
	getRes, _, err := n.con.PC.VM.Get(getRequest)
	if err != nil {
		return err
	}

	updateRequest := &pc.VMUpdateRequest{}
	updateRequestData := &pc.VMUpdateRequestData{}
	updateRequestData.Spec = getRes.Spec
	updateRequestData.APIVersion = &getRes.APIVersion
	updateRequestData.Metadata = &getRes.Metadata

	updateRequest.UUID = c.Args().First()
	updateRequest.Data = *updateRequestData

	updateRequest.Data.Spec.Resources.PowerState = nutanix.String(powerState)

	_, _, err = n.con.PC.VM.Update(updateRequest)
	if err != nil {
		return err
	}

	fmt.Println("virtual machine updated to power state: ", powerState)

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

	n.tr.AppendBulk(data)
	n.tr.Render()
	return nil
}
