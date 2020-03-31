package rspace

import (
	"fmt"
	"strings"
	"errors"
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
	Links      []Link `json: "_links"`
}

//Summary information about a Document
type DocumentInfo struct {
	*IndentifiableNamable
	Created        string
	LastModified   string
	ParentFolderId int
	Signed         bool
	Tags           string
	FormInfo       FormInfo
	UserInfo       UserInfo
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

type Folder struct {
	*IndentifiableNamable
	Created        string
	LastModified   string
	IsNotebook     bool `json :"notebook"`
	ParentFolderId int
}

type FolderTreeItem struct {
	*IndentifiableNamable
	Created      string
	LastModified string
	IsNotebook   bool `json :"notebook"`
	Type         string
}
type FolderList struct {
	Records    []FolderTreeItem
	TotalHits  int
	PageNumber int
	Links      []Link `json: "_links"`
}

type FolderPost struct {
	Name           string `json:"name"`
	IsNotebook     bool   `json:"notebook"`
	ParentFolderId int    `json:"parentFolderId"`
}

type FileList struct {
	TotalHits  int
	PageNumber int
	Links      []Link `json: "_links"`
	Files      []FileInfo
}

type IndentifiableNamable struct {
	Id       int
	GlobalId string
	Name     string
}

type Field struct {
	*IndentifiableNamable
	Type         string
	Content      string
	LastModified string
	Files        []FileInfo
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
	Id           int
	Name     string
	Type        string
	SharedFolderId    int
	Members [] struct{
		Id int
		Username string
		Role string
	 }
}

type Link struct {
	Link string
	Rel  string
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
	pwordToPrint:="not set..."
	if len(upost.Password) > 0 {
		pwordToPrint="..."
	}
	return fmt.Sprintf("username=%s,email=%s,firstName=%s,lastName=%s,password=%s,role=%s,affiliation=%s,apiKey=%s",
	 upost.Username, upost.Email,upost.FirstName, upost.LastName, pwordToPrint, upost.Role, upost.Affiliation, upost.ApiKey)
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
func (b *UserPostBuilder) build() (*UserPost, error){
	rc := UserPost{}
	if len (b.Username) < 6 {
		return nil, errors.New("username must be >= 6 characters")
	}	
	if len (b.Password) < 6 {
		return nil, errors.New("Password must be >= 8 characters")
	}	
	if len (b.FirstName) ==0  {
		return nil, errors.New("Please supply first name")
	}	
	if len (b.LastName) ==0  {
		return nil, errors.New("Please supply last name")
	}
	if len (string(b.Email)) < 3   {
		return nil, errors.New("Please supply valid email address")
	}
	rc.FirstName=b.FirstName
	rc.Password=b.Password
	rc.Username=b.Username
	rc.LastName=b.LastName
	rc.Email=string(b.Email)
	rc.Role=userRoles[b.Role]
	rc.Affiliation=b.Affiliation
	rc.ApiKey=b.ApiKey
	return &rc, nil
}
// GroupPost is serialized to JSON. Client code  should use GroupPostNew to create this object.
type GroupPost struct {
	DisplayName string `json:"displayName"`
	Members []UserGroupPost `json:"members"`
}

//GroupPostNew performs validated construction of a GroupPost object
func GroupPostNew (name string, userGroups []UserGroupPost) (*GroupPost, error) {
	rc := GroupPost {}
	if len (name) == 0  {
		return nil, errors.New("Please supply a name for the group")
	}
	rc.DisplayName = name
	if len (userGroups) == 0 {
		return nil, errors.New("Please supply at least 1 group member")
	}
	var piExists bool
	for _, upost := range userGroups {
		if  find(userInGroupRoles, upost.RoleInGroup) < 0 {
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

var userInGroupRoles = []string {"DEFAULT", "RS_LAB_ADMIN", "PI"}

//UserGroup post defines a single user's membership role within a group.
type UserGroupPost struct {
	Username string `json:"username"`
	RoleInGroup string `json:"roleInGroup"`
}

func find(slice []string, val string) int {
    for i, item := range slice {
        if item == val {
            return i
        }
    }
    return -1
}
