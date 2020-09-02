package rspace

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ContainerPost struct {
	Name            string          `yaml:"name" json:"name"`
	Containers      []ContainerPost `yaml:"containers"`
	ParentContainer *ParentRef      `json:"parentContainer,omitempty"`
}

type ParentRef struct {
	Id int `json:"id"`
}
type Container struct {
	IdentifiableNamable
	ParentId int `json:"parentContainer"`
}
type InventoryService struct {
	BaseService
}

func (is *InventoryService) CreateContainers(toCreate *ContainerPost) {
	result, err := is.createContainer(toCreate)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range toCreate.Containers {
		copy := v
		copy.ParentContainer = &ParentRef{Id: result.Id}
		is.CreateContainers(&copy)
	}
}

func (is *InventoryService) createContainer(toCreate *ContainerPost) (*Container, error) {
	url := is.containerUrl()
	data, err := is.doPostJsonBody(toCreate, url)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = Container{}
	json.Unmarshal(data, &result)
	fmt.Println(result)

	return &result, nil
}
func (fs *InventoryService) containerUrl() string {
	return strings.Replace(fs.BaseUrl.String(), "api/v1", "api/inventory/v1", 1) + "/containers"
}
