package rspace

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
}

// factory method to return a RecordListingConfig with default values
func NewRecordListingConfig() RecordListingConfig {
	return RecordListingConfig{
		PageSize:   20,
		OrderBy:    "lastModified",
		PageNumber: 0,
		SortOrder:  "desc",
	}
}

//toParams generates a url.Values object containing web params
func (config *RecordListingConfig) toParams() url.Values {
	params := url.Values{}
	params.Add("pageSize", strconv.Itoa(config.PageSize))
	params.Add("pageNumber", strconv.Itoa(config.PageNumber))
	params.Add("orderBy", config.OrderBy+" "+config.SortOrder)
	return params
}

// DocumentPost contains optional data for creating or updating an RSpace document
type DocumentPost struct {
	Name   string         `json:"name,omitempty"`
	Tags   string         `json:"tags,omitempty"`
	FormID FormId         `json:"form,omitempty"`
	Fields []FieldContent `json:"fields"`
}
type FieldContent struct {
	Content string `json:"content"`
	//Id will be null when posting, can be set if editing a specific field
	Id int `json:"id,omitempty"`
}
type FormId struct {
	Id int `json:"id,omitempty"`
}

// constructor for a new document
func DocumentPostNew(name string, tags string, formId int, content []string) *DocumentPost {
	id := FormId{formId}
	c := make([]FieldContent, 0)
	for _, v := range content {
		c = append(c, FieldContent{Content: v})
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
	*IdentifiableNamable
	Created        string
	LastModified   string
	ParentFolderId int
	Signed         bool
	Tags           string
	Form           Form     `json:"form"`
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
	*IdentifiableNamable
	ContentType string
	Size        int
	Caption     string
	Created     string
	Version     int
}

func (fi *FileInfo) CreatedTime() (time.Time, error) {
	return parseTimestamp(fi.Created)
}

// FileUploadConfig configures file upload. filePath is required, other
// properties are optional.
type FileUploadConfig struct {
	Caption  string
	FilePath string
	FolderId int
}

type Folder struct {
	*IdentifiableNamable
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
	*IdentifiableNamable
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

type IdentifiableNamable struct {
	Id       int
	GlobalId string
	Name     string
}

func (item IdentifiableNamable) GetId() int {
	return item.Id
}

func (item IdentifiableNamable) GetGlobalId() string {
	return item.GlobalId
}

func (item IdentifiableNamable) GetName() string {
	return item.Name
}

type Field struct {
	*IdentifiableNamable
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

type UserInfo struct {
	Id           int
	Username     string
	Email        string
	FirstName    string
	LastName     string
	HomeFolderId int
}

type GroupList struct {
	Groups []*GroupInfo
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

type SharePost struct {
	ItemsToShare []int        `json:"itemsToShare"`
	Groups       []GroupShare `json:"groups"`
	Users        []UserShare  `json:"users"`
}
type GroupShare struct {
	Id             int    `json:"id"`
	Permission     string `json:"permission"`
	SharedFolderId int    `json:"sharedFolderId"`
}

type UserShare struct {
	Id         int    `json:"id"`
	Permission string `json:"permission"`
}

//ShareInfoList represents the results of posting a /share request
type ShareInfoList struct {
	ShareInfos   []*ShareResult
	FailedShares []int
	Links        []string `json:"_links"`
}

type SharedItemList struct {
	Shares     []*ShareResult
	TotalHits  int
	PageNumber int
	Links      []string `json:"_links"`
}

//ShareResult represents a successful share of a single resource
type ShareResult struct {
	Id         int
	ItemId     int    `json:"sharedItemId"`
	ItemName   string `json:"shareItemName"`
	TargetType string `json:"sharedTargetType"`
	Permission string
	Links      []string `json:"_links"`
}

type Link struct {
	Link string
	Rel  string
}

type Email string

type Permission int

const (
	Read Permission = iota
	Edit
)

type UserRoleType int

const (
	User UserRoleType = iota
	Pi
	Admin
	Sysadmin
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
	username    string
	email       Email
	firstName   string
	lastName    string
	password    string
	role        UserRoleType
	affiliation string
	apiKey      string
}

func (b *UserPostBuilder) Username(username string) *UserPostBuilder {
	b.username = username
	return b
}
func (b *UserPostBuilder) Password(password string) *UserPostBuilder {
	b.password = password
	return b
}
func (b *UserPostBuilder) Email(emailAddress Email) *UserPostBuilder {
	b.email = emailAddress
	return b
}
func (b *UserPostBuilder) FirstName(firstName string) *UserPostBuilder {
	b.firstName = firstName
	return b
}
func (b *UserPostBuilder) LastName(lastName string) *UserPostBuilder {
	b.lastName = lastName
	return b
}
func (b *UserPostBuilder) Role(role UserRoleType) *UserPostBuilder {
	b.role = role
	return b
}
func (b *UserPostBuilder) Affiliation(affiliation string) *UserPostBuilder {
	b.affiliation = affiliation
	return b
}
func (b *UserPostBuilder) ApiKey(apiKey string) *UserPostBuilder {
	b.apiKey = apiKey
	return b
}
func (b *UserPostBuilder) Build() (*UserPost, error) {
	rc := UserPost{}
	if len(b.username) < 6 {
		return nil, errors.New("username must be >= 6 characters")
	}
	if len(b.password) < 6 {
		return nil, errors.New("Password must be >= 8 characters")
	}
	if len(b.firstName) == 0 {
		return nil, errors.New("Please supply first name")
	}
	if len(b.lastName) == 0 {
		return nil, errors.New("Please supply last name")
	}
	if len(string(b.email)) < 3 {
		return nil, errors.New("Please supply valid email address")
	}
	rc.FirstName = b.firstName
	rc.Password = strings.TrimSpace(b.password)
	rc.Username = b.username
	rc.LastName = b.lastName
	rc.Email = string(b.email)
	rc.Role = userRoles[b.role]
	rc.Affiliation = b.affiliation
	rc.ApiKey = b.apiKey
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

// ActivityList encapsulates search results for audit events
type ActivityList struct {
	Activities []Activity
	Links      []Link `json:"_links"`
	TotalHits  int
	PageNumber int
}

// Activity holds information abuot a particular audit event.
// The Payload field holds arbitrary data that is specific for each event type.
type Activity struct {
	Username, FullName, Domain, Action string
	Timestamp                          string
	Payload                            interface{}
}

// TimestampTime offers the timestamp of the audit event.
func (a *Activity) TimestampTime() (time.Time, error) {
	return parseTimestamp(a.Timestamp)
}

// GlobalId is  a Unique identifier for an RSpace object, e.g. 'GL1234' or 'SD5678'
type GlobalId string

// ActivityQuery encapsulates a query to the /activities endpoint. Either use directly or use the ActivityQueryBuilder, which
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

// Build generates an ActivityQuery from the builder, that is validated and ready to
// send.
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

// FormList holds the results of listing Forms
type FormList struct {
	TotalHits  int
	PageNumber int
	Forms      []Form
	Links      []Link `json:"_links"`
}

// Form holds basic information about a Form
type Form struct {
	*IdentifiableNamable
	Version   int
	FormState string
	StableId  string
	Links     []Link `json:"_links"`
	Tags      string
}

type ExportFormat int

const (
	XML_FORMAT ExportFormat = iota
	HTML_FORMAT
)

func (scope ExportFormat) String() string {
	switch scope {
	case XML_FORMAT:
		return "xml"
	case HTML_FORMAT:
		return "html"
	}
	return ""
}

type ExportScope int

const (
	USER_EXPORT_SCOPE ExportScope = iota
	GROUP_EXPORT_SCOPE
)

func (scope ExportScope) String() string {
	switch scope {
	case USER_EXPORT_SCOPE:
		return "user"
	case GROUP_EXPORT_SCOPE:
		return "group"
	}
	return ""
}

type ExportPost struct {
	// XML, HTML
	Format ExportFormat
	// user,group
	Scope ExportScope
	// The id of a user or group to export
	Id int
}

// NewExportPost generates an ExportPost of type USER_EXPORT_SCOPE
// to HTML format, set as defaults
func NewExportPost() ExportPost {
	return ExportPost{HTML_FORMAT, USER_EXPORT_SCOPE, 0}
}

// Job holds information about a long-running request
type Job struct {
	Status          string
	Id              int
	PercentComplete float32
	Links           []Link `json:"_links"`
	Result          JobResult
}

type JobResult struct {
	ExpiryDate string
	Size       int
	Checksum   string
	Algorithm  string
}

// DownloadLink returns a URL of the link to access the export,
// if the Job is in COMPLETED state and has a download link.
// Otherwise returns nil
func (job *Job) DownloadLink() *url.URL {
	if job.IsCompleted() {
		if len(job.Links) > 0 {
			fmt.Println("links" + job.Links[0].Link)
			urlO, _ := url.Parse(job.Links[0].Link)
			return urlO
		}
	}
	return nil
}

func (job *Job) IsCompleted() bool {
	return job.Status == "COMPLETED"
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
