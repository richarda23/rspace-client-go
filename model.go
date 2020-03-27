package rspace

import (
	"fmt"
)

// Status stores response from /status endpoint
type Status struct {
	Message       string `json: "message"`
	RSpaceVersion string `json: "rspaceVersion"`
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
	Created        string
	LastModified   string
	IsNotebook     bool `json :"notebook"`
	Type string
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

type Link struct {
	Link string
	Rel  string
}
