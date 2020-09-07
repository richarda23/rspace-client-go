package rspace

import (
	"fmt"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestCreateContainer(t *testing.T) {
	filePath := "testdata/containers.yaml"
	con := &ContainerPost{}

	data, _ := ioutil.ReadFile(filePath)
	err := yaml.Unmarshal([]byte(data), con)
	if err != nil {
		return
	}
	con.ParentContainer = nil
	res, _ := webClient.inventoryS.CreateContainers(con)
	fmt.Println(len(res.Containers))

}
