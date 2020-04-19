package rspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"strconv"
	"time"
)

// Status stores response from /status endpoint
type Status struct {
	Message       string `json:"message"`
	RSpaceVersion string `json:"rspaceVersion"`
}

// configures pagination and verbosity for listings
type RecordListingConfig struct {
	SortOrder  string
	PageSize   int
	PageNumber int
	OrderBy    string
	Quiet      bool
}

// factory method to return a RecordListingConfig with default values
func NewRecordListingConfig() RecordListingConfig {
	return RecordListingConfig{
		PageSize:   20,
		OrderBy:    "lastModified",
		PageNumber: 0,
		SortOrder:  "desc",
		Quiet:      false,
	}
}

//toParams generates a url.Values object containing web params
func (config *RecordListingConfig) toParams() url.Values {
	params := url.Values{}
	params.Add("pageSize",  strconv.Itoa(config.PageSize))
	params.Add("pageNumber",  strconv.Itoa(config.PageNumber))
	params.Add("orderBy", config.OrderBy + " " + config.SortOrder)
	return params
}

type DocumentPost struct {
	Name   string         `json:"name"`
	Tags   string         `json:"tags"`
	FormId FormId         `json:"formId"`
	Fields []FieldContent `json:"fields"`
}
type FieldContent struct {
	Content string `json:"content"`
}
type FormId struct {
	Id int
}

// constructor for a new document
func DocumentPostNew(name string, tags string, formId int, content []string) *DocumentPost {
	id := FormId{formId}
	c := make([]FieldContent, 10)
	for _, v := range content {
		fmt.Println(v)
		//append(c, FieldContent{v})
	}
	post := DocumentPost{name, tags, id, c}
	return &post
}

type DocumentList struct {
	Documents  []DocumentInfo
	TotalHits  int
	PageNumber int
	Links      []Link `json:"_links"`
}
type UserList struct {
	Users      []UserInfo
	TotalHits  int
	PageNumber int
	Links      []Link `json:"_links"`
}

//Summary information about a Document
type DocumentInfo struct {
	*IndentifiableNamable
	Created        string
	LastModified   string
	ParentFolderId int
	Signed         bool
	Tags           string
	FormInfo       FormInfo `json:"form"`
	UserInfo       UserInfo `json:"owner"`
}

func (di *DocumentInfo) CreatedTime() (time.Time, error) {
	return parseTimestamp(di.Created)
}
func (di *DocumentInfo) LastModifiedTime() (time.Time, error) {
	return parseTimestamp(di.LastModified)
}

func parseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

// FileInfo holds metadata about Files
type FileInfo struct {
	*IndentifiableNamable
	ContentType string
	Size        int
	Caption     string
	Created     string
	Version     int
}

func (fi *FileInfo) CreatedTime() (time.Time, error) {
	return parseTimestamp(fi.Created)
}

type Folder struct {
	*IndentifiableNamable
	Created        string
	LastModified   string
	IsNotebook     bool `json:"notebook"`
	ParentFolderId int
}

func (f *Folder) CreatedTime() (time.Time, error) {
	return parseTimestamp(f.Created)
}
func (f *Folder) LastModifiedTime() (time.Time, error) {
	return parseTimestamp(f.LastModified)
}

type FolderTreeItem struct {
	*IndentifiableNamable
	Created      string
	LastModified string
	IsNotebook   bool `json:"notebook"`
	Type         string
}

func (f *FolderTreeItem) CreatedTime() (time.Time, error) {
	return parseTimestamp(f.Created)
}
func (f *FolderTreeItem) LastModifiedTime() (time.Time, error) {
	return parseTimestamp(f.LastModified)
}

type FolderList struct {
	Records    []FolderTreeItem
	TotalHits  int
	PageNumber int
	Links      []Link `json:"_links"`
}

type FolderPost struct {
	Name           string `json:"name,omitempty"`
	IsNotebook     bool   `json:"notebook,omitempty"`
	ParentFolderId int    `json:"parentFolderId,omitempty"`
}

type FileList struct {
	TotalHits  int
	PageNumber int
	Links      []Link `json:"_links"`
	Files      []FileInfo
}
// BasicInfo provides simple information common to many RSpace resources
type BasicInfo interface {
	GetName() string
	GetId() int
	GetGlobalId() string
}

type IndentifiableNamable struct {
	Id       int
	GlobalId string
	Name     string
}

func (item IndentifiableNamable) GetId() int {
	return item.Id;
}

func (item IndentifiableNamable) GetGlobalId() string {
	return item.GlobalId;
}

func (item IndentifiableNamable) GetName() string {
	return item.Name;
}

type Field struct {
	*IndentifiableNamable
	Type         string
	Content      string
	LastModified string
	Files        []FileInfo
}

func (f *Field) LastModifiedTime() (time.Time, error) {
	return parseTimestamp(f.LastModified)
}

//Full document including content
type Document struct {
	*DocumentInfo
	Fields []Field
}

type FormInfo struct {
	Id       int
	GlobalId string
	Name     string
	StableId string
	Version  int
}

type UserInfo struct {
	Id           int
	Username     string
	Email        string
	FirstName    string
	LastName     string
	HomeFolderId int
}

type GroupInfo struct {
	Id             int
	Name           string
	Type           string
	SharedFolderId int
	Members        []struct {
		Id       int
		Username string
		Role     string
	}
}

type Link struct {
	Link url.URL
	Rel  string
}

func (l *Link) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string
	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "link" {
			u, err := url.Parse(v)
			if err != nil {
				return err
			}
			l.Link = *u
		}
		if strings.ToLower(k) == "rel" {
			l.Rel = v
		}
	}
	return nil
}

//RSpaceError encapsulates server or client side errors leading to a request being rejected.
type RSpaceError struct {
	Status       string
	HttpCode     int
	InternalCode int
	Message      string
	Errors       []string
	Timestamp    string `json:"iso8601Timestamp"`
}

func (f *RSpaceError) CreatedTime() (time.Time, error) {
	return parseTimestamp(f.Timestamp)
}

func (rsError *RSpaceError) String() string {
	if rsError.HttpCode >= 400 && rsError.HttpCode < 500 {
		return formatErrorMsg(rsError, "Client")
	} else if rsError.HttpCode > 500 {
		return formatErrorMsg(rsError, "Server")
	} else {
		return formatErrorMsg(rsError, "Unknown")
	}
}

func (rsError *RSpaceError) Error() string {
	return rsError.String()
}

func formatErrorMsg(rsError *RSpaceError, errType string) string {
	concatenateErrM := strings.Join(rsError.Errors, "\n")
	rc := fmt.Sprintf("%s error:httpCode=%d, status=%s, internalCode=%d, timestamp=%s,  message=%s\nErrors: %s",
		errType, rsError.HttpCode, rsError.Status, rsError.InternalCode, rsError.Timestamp, rsError.Message, concatenateErrM)
	return rc
}

type Email string
type UserRoleType int

const (
	user UserRoleType = iota
	pi
	admin
	sysadmin
)

var userRoles = [4]string{"ROLE_USER", "ROLE_PI", "ROLE_ADMIN", "ROLE_SYSADMIN"}

//
type UserPost struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Affiliation string `json:"affiliation,omitempty"`
	ApiKey      string `json:"apiKey,omitempty"`
}

func (upost *UserPost) String() string {
	pwordToPrint := "not set..."
	if len(upost.Password) > 0 {
		pwordToPrint = "..."
	}
	return fmt.Sprintf("username=%s,email=%s,firstName=%s,lastName=%s,password=%s,role=%s,affiliation=%s,apiKey=%s",
		upost.Username, upost.Email, upost.FirstName, upost.LastName, pwordToPrint, upost.Role, upost.Affiliation, upost.ApiKey)
}

// Use this to build a UserPost object to create a new user from.
type UserPostBuilder struct {
	Username    string
	Email       Email
	FirstName   string
	LastName    string
	Password    string
	Role        UserRoleType
	Affiliation string
	ApiKey      string
}

func (b *UserPostBuilder) username(username string) *UserPostBuilder {
	b.Username = username
	return b
}
func (b *UserPostBuilder) password(password string) *UserPostBuilder {
	b.Password = password
	return b
}
func (b *UserPostBuilder) email(emailAddress Email) *UserPostBuilder {
	b.Email = emailAddress
	return b
}
func (b *UserPostBuilder) firstName(firstName string) *UserPostBuilder {
	b.FirstName = firstName
	return b
}
func (b *UserPostBuilder) lastName(lastName string) *UserPostBuilder {
	b.LastName = lastName
	return b
}
func (b *UserPostBuilder) role(role UserRoleType) *UserPostBuilder {
	b.Role = role
	return b
}
func (b *UserPostBuilder) affiliation(affiliation string) *UserPostBuilder {
	b.Affiliation = affiliation
	return b
}
func (b *UserPostBuilder) apiKey(apiKey string) *UserPostBuilder {
	b.ApiKey = apiKey
	return b
}
func (b *UserPostBuilder) build() (*UserPost, error) {
	rc := UserPost{}
	if len(b.Username) < 6 {
		return nil, errors.New("username must be >= 6 characters")
	}
	if len(b.Password) < 6 {
		return nil, errors.New("Password must be >= 8 characters")
	}
	if len(b.FirstName) == 0 {
		return nil, errors.New("Please supply first name")
	}
	if len(b.LastName) == 0 {
		return nil, errors.New("Please supply last name")
	}
	if len(string(b.Email)) < 3 {
		return nil, errors.New("Please supply valid email address")
	}
	rc.FirstName = b.FirstName
	rc.Password = b.Password
	rc.Username = b.Username
	rc.LastName = b.LastName
	rc.Email = string(b.Email)
	rc.Role = userRoles[b.Role]
	rc.Affiliation = b.Affiliation
	rc.ApiKey = b.ApiKey
	return &rc, nil
}

// GroupPost is serialized to JSON. Client code  should use GroupPostNew to create this object.
type GroupPost struct {
	DisplayName string          `json:"displayName"`
	Members     []UserGroupPost `json:"members"`
}

//GroupPostNew performs validated construction of a GroupPost object
func GroupPostNew(name string, userGroups []UserGroupPost) (*GroupPost, error) {
	rc := GroupPost{}
	if len(name) == 0 {
		return nil, errors.New("Please supply a name for the group")
	}
	rc.DisplayName = name
	if len(userGroups) == 0 {
		return nil, errors.New("Please supply at least 1 group member")
	}
	var piExists bool
	for _, upost := range userGroups {
		if find(userInGroupRoles, upost.RoleInGroup) < 0 {
			return nil, errors.New("Please supply a valid group role for this user")
		}
		if upost.RoleInGroup == "PI" {
			piExists = true
		}
	}
	if !piExists {
		return nil, errors.New("There must be exactly 1 PI in the group")
	}

	rc.Members = userGroups
	return &rc, nil
}

var userInGroupRoles = []string{"DEFAULT", "RS_LAB_ADMIN", "PI"}

//UserGroup post defines a single user's membership role within a group.
type UserGroupPost struct {
	Username    string `json:"username"`
	RoleInGroup string `json:"roleInGroup"`
}

type ActivityList struct {
	Activities []Activity
	Links      []Link `json:"_links"`
	TotalHits  int
	PageNumber int
}
type Activity struct {
	Username, FullName, Domain, Action string
	Timestamp                          string
	Payload                            interface{}
}

func (a *Activity) TimestampTime() (time.Time, error) {
	return parseTimestamp(a.Timestamp)
}

// GlobalId is  a Unique identifier for an RSpace object, e.g. 'GL1234' or 'SD5678'
type GlobalId string

// ActivityQuery encapsulates a auery to the /activities endpoint. Either use directly or use the ActivityQueryBuilder, which
// provides more convenient construction and validation
type ActivityQuery struct {
	Domains  []string
	Actions  []string
	Oid      string
	Users    []string
	DateFrom time.Time
	DateTo   time.Time
}

// ActivityQueryBuilder provides convenient methods to construct a query to the /activities endpoint
type ActivityQueryBuilder struct {
	domains  []string
	actions  []string
	oid      GlobalId
	users    []string
	dateFrom time.Time
	dateTo   time.Time
}

func (b *ActivityQueryBuilder) Domain(domain string) *ActivityQueryBuilder {
	b.domains = makeStringSlice(b.domains, domain)
	return b
}

func (b *ActivityQueryBuilder) Action(action string) *ActivityQueryBuilder {
	b.actions = makeStringSlice(b.actions, action)
	return b
}

func (b *ActivityQueryBuilder) User(user string) *ActivityQueryBuilder {
	b.users = makeStringSlice(b.users, user)
	return b
}

//DateFrom specifies a lower bound on the time stamp of an activity
func (b *ActivityQueryBuilder) DateFrom(dateFrom time.Time) *ActivityQueryBuilder {
	b.dateFrom = dateFrom
	return b
}

//DateTo specifies an upper bound on the time stamp of an activity.
func (b *ActivityQueryBuilder) DateTo(dateTo time.Time) *ActivityQueryBuilder {
	b.dateTo = dateTo
	return b
}

//Oid restricts the search to activities involving the specific item.
func (b *ActivityQueryBuilder) Oid(oid GlobalId) *ActivityQueryBuilder {
	b.oid = oid
	return b
}
func (b *ActivityQueryBuilder) Build() (*ActivityQuery, error) {
	rc := ActivityQuery{}
	rc.Domains = b.domains
	rc.Actions = b.actions
	rc.Users = b.users
	if !b.dateFrom.IsZero() && !b.dateTo.IsZero() && b.dateFrom.After(b.dateTo) {
		return nil, errors.New(fmt.Sprintf("from Date cannot be before to date"))
	}
	rc.DateFrom = b.dateFrom
	rc.DateTo = b.dateTo
	if len(b.oid) > 0 {
		if match, _ := regexp.MatchString("[A-Z]{2}\\d+", string(b.oid)); match == true {
			rc.Oid = string(b.oid)
		} else {
			return nil, errors.New(fmt.Sprintf("'%s' is not a valid global ID", b.oid))
		}
	}
	return &rc, nil
}

type FormList struct {
	TotalHits  int
	PageNumber int
	Forms      []Form
	Links      []Link `json:"_links"`
}

type Form struct {
	*IndentifiableNamable
	Version   int
	FormState string
	StableId  string
	Links     []Link `json:"_links"`
	Tags      string
}

func makeStringSlice(existingSl []string, toAdd string) []string {
	if len(existingSl) == 0 {
		existingSl = make([]string, 0, 0)
	}
	return append(existingSl, toAdd)
}

func find(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}
