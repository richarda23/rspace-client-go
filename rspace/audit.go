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

func (as *ActivityService) auditUrl() string {
	return as.BaseUrl.String()+ "/activity"
}

// Activities queries the audit trail for activities, by user, date or activity type
func (as *ActivityService) Activities(q *ActivityQuery, pgCrit RecordListingConfig) (*ActivityList, error) {
	time.Sleep(as.Delay)
	urlStr := as.auditUrl()
	var encodedParams string
	pgCrit.OrderBy="date"
	var params url.Values = pgCrit.toParams()
	if q != nil {
		if len(q.Users) > 0 {
			params.Add("usernames", strings.Join(q.Users, ","))
		}
		if len(q.Domains) > 0 {
			params.Add("domains", strings.ToUpper(strings.Join(q.Domains, ",")))
		}
		if len(q.Actions) > 0 {
			params.Add("actions", strings.ToUpper(strings.Join(q.Actions, ",")))
		}
		if len(q.Oid) > 0 {
			params.Add("oid", q.Oid)
		}
		if !q.DateFrom.IsZero() {
			params.Add("dateFrom", q.DateFrom.Format("2006-01-02"))
		}
		if !q.DateTo.IsZero() {
			params.Add("dateTo", q.DateTo.Format("2006-01-02"))
		}
		encodedParams = params.Encode()
	}
	if len(encodedParams) > 0 {
		urlStr = urlStr + "?" + encodedParams
	}
	data, err := as.doGet(urlStr)
	if err != nil {
		return nil, err
	}
	var result = ActivityList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
