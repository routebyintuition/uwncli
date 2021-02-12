package main

import (
	"strconv"

	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/urfave/cli/v2"
)

func (n *NCLI) listSubnets(c *cli.Context) error {

	ListRequest := new(pc.SubnetListRequest)

	var vmListLoop []pc.Entities
	totalMatches := 0
	offset := 0
	currentMatches := -1
	var err error
	getRes := &pc.SubnetListResponse{}

	for totalMatches > currentMatches {
		if currentMatches == -1 {
			currentMatches = 0
		}

		ListRequest.Offset = offset

		getRes, _, err = n.con.PC.Subnet.List(ListRequest)
		if err != nil {
			return err
		}

		currentMatches += *getRes.Metadata.Length
		totalMatches = *getRes.Metadata.TotalMatches
		offset += ListRequest.Length
		vmListLoop = append(vmListLoop, getRes.Entities...)
	}

	n.tr.SetHeader([]string{"Name", "UUID", "Network", "CIDR", "VLAN", "TYPE"})
	n.tr.SetFooter([]string{"", "", "", "", "TOTAL", strconv.Itoa(*getRes.Metadata.TotalMatches)})

	data := [][]string{}

	for _, entityValue := range vmListLoop {
		gwAddress := "0.0.0.0"
		gwPrefix := "0"

		if entityValue.Spec.Resources.IPConfig != nil {
			gwAddress = entityValue.Spec.Resources.IPConfig.DefaultGatewayIP
			gwPrefix = strconv.Itoa(entityValue.Spec.Resources.IPConfig.PrefixLength)
		}

		data = append(data, []string{*entityValue.Spec.Name, *entityValue.Metadata.UUID, gwAddress, gwPrefix, strconv.Itoa(*entityValue.Spec.Resources.VlanID), *entityValue.Spec.Resources.SubnetType})
	}
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

func (n *NCLI) getSubnetUUIDList() ([]string, error) {

	ListRequest := new(pc.SubnetListRequest)

	data := []string{}

	var vmListLoop []pc.Entities
	totalMatches := 0
	offset := 0
	currentMatches := -1
	var err error
	getRes := &pc.SubnetListResponse{}

	for totalMatches > currentMatches {
		if currentMatches == -1 {
			currentMatches = 0
		}

		ListRequest.Offset = offset

		getRes, _, err = n.con.PC.Subnet.List(ListRequest)
		if err != nil {
			return nil, err
		}

		currentMatches += *getRes.Metadata.Length
		totalMatches = *getRes.Metadata.TotalMatches
		offset += ListRequest.Length
		vmListLoop = append(vmListLoop, getRes.Entities...)
	}

	for _, entityValue := range vmListLoop {
		data = append(data, *entityValue.Metadata.UUID)
	}

	return data, nil
}
