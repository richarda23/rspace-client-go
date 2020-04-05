package rspace

import (
	"encoding/json"
	//	"fmt"
	"net/url"
	"strings"
	"time"
)

type ActivityService struct {
	BaseService
}

func auditUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/activity"
}

// Activities queries the audit trail for activities, by user, date or activity type
func (fs *ActivityService) Activities(q *ActivityQuery) (*ActivityList, error) {
	time.Sleep(fs.Delay)
	urlStr := auditUrl()
	var encodedParams string
	if q != nil {
		params := url.Values{}
		if len(q.Users) > 0 {
			params.Add("users", strings.Join(q.Users, ","))
		}
		if len(q.Domains) > 0 {
			params.Add("domains", strings.Join(q.Domains, ","))
		}
		if len(q.Actions) > 0 {
			params.Add("actions", strings.Join(q.Actions, ","))
		}
		if len(q.Oid) > 0 {
			params.Add("oid", q.Oid)
		}
		if !q.DateFrom.IsZero() {
			params.Add("dateFrom", q.DateFrom.Format("2006-02-01"))
		}
		if !q.DateTo.IsZero() {
			params.Add("dateTo", q.DateTo.Format("2006-02-01"))
		}
		encodedParams = params.Encode()
	}
	if len(encodedParams) > 0 {
		urlStr = urlStr + "?" + encodedParams
	}
	data, err := DoGet(urlStr)
	if err != nil {
		return nil, err
	}
	var result = ActivityList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
