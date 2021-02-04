package main

import (
	"strconv"

	"github.com/routebyintuition/ntnx-go-sdk/karbon"
	"github.com/urfave/cli/v2"
)

func (n *NCLI) krbnListClusters(c *cli.Context) error {

	cListRequest := new(karbon.ClusterListRequest)

	getRes, _, err := n.con.Karbon.Cluster.List(cListRequest)
	if err != nil {
		return err
	}
	clusterCount := len(*getRes)

	n.tr.SetHeader([]string{"Name", "UUID", "Address", "Version"})

	n.tr.SetFooter([]string{"", "", "Total", strconv.Itoa(clusterCount)})

	data := [][]string{}

	for _, entityValue := range *getRes {
		data = append(data, []string{entityValue.Name, entityValue.UUID, entityValue.KubeapiServerIpv4Address, entityValue.Version})
	}
	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}
