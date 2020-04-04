package rspace

import (
	"fmt"
	"testing"
	"time"
)

var activityService *ActivityService = &ActivityService{
	BaseService: BaseService{
		Delay: time.Duration(100) * time.Millisecond}}

func TestActivityGet(t *testing.T) {
	data, err := activityService.Activities()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
