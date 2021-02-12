package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/routebyintuition/ntnx-go-sdk/pc"
	"github.com/routebyintuition/ntnx-go-sdk/pe"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// clusterList returns a list of all PC clusters
func (n *NCLI) clusterList(c *cli.Context) error {
	ListRequest := new(pc.ClusterListRequest)

	getRes, _, err := n.con.PC.Cluster.List(ListRequest)
	if err != nil {
		return err
	}

	n.tr.SetHeader([]string{"Cluster Name", "UUID", "Hypervisor", "Version"})

	data := [][]string{}

	entityCount := 0

	for _, entityValue := range getRes.Entities {
		if entityValue.Status.Resources.Nodes != nil {
			entityCount++
			hypervisorType := "AHV"
			hypervisorVersion := "0.0"

			for _, hList := range entityValue.Status.Resources.Nodes.HypervisorServerList {
				if hList.IP != "127.0.0.1" {
					hypervisorType = hList.Type
					hypervisorVersion = hList.Version
				}
			}
			data = append(data, []string{entityValue.Status.Name, *entityValue.Metadata.UUID, hypervisorType, hypervisorVersion})
		}
	}

	n.tr.SetFooter([]string{"", "", "TOTAL", strconv.Itoa(entityCount)})

	n.tr.AppendBulk(data)
	n.tr.Render()

	return nil
}

// clusterGet returns details of cluster from PE
func (n *NCLI) clusterGet(c *cli.Context) error {
	ListRequest := new(pe.ClusterGetRequest)

	getRes, _, err := n.con.PE.Cluster.Get(ListRequest)
	if err != nil {
		return err
	}

	yamlDoc, err := yaml.Marshal(&getRes)
	if err != nil {
		log.Fatalf("error on yaml converstion: %v", err)
	}
	fmt.Printf("%s\n", string(yamlDoc))

	return nil
}
