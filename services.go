package rspace

import (
	"time"
)

var activityService *ActivityService = &ActivityService{
	BaseService: baseService(),
}

var ds *DocumentService = &DocumentService{
	BaseService: baseService(),
}
var sysads *SysadminService = &SysadminService{
	BaseService: baseService(),
}

func baseService() BaseService {
	return BaseService{
		Delay: time.Duration(100) * time.Millisecond}
}
