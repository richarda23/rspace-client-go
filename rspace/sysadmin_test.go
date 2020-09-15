package rspace

import (
	"fmt"
	"testing"
	"time"
	//"github.com/op/go-logging"
	//"errors"
)

func TestUserNew(t *testing.T) {
	// given
	userPost := createRandomUser(Pi)
	got, err := webClient.UserNew(userPost)
	if err != nil {
		Log.Error(err)
		t.Fatalf("unexpected error " + err.Error())
	}
	if got.Id == 0 {
		fail(t, "Id was nil but should be set")
	}
	assertStringEquals(t, userPost.Username, got.Username, "")
}
func TestUsers(t *testing.T) {
	userPost := createRandomUser(Pi)
	got, _ := webClient.UserNew(userPost)
	// default ordering is creationDate desc, so new user should be first:
	var cfg RecordListingConfig = NewRecordListingConfig()
	// all users were created before a time in the future
	cfg.OrderBy = "creationDate"
	userList, e := webClient.Users(time.Now().AddDate(1, 0, 0), time.Now().AddDate(1, 0, 0), cfg)

	assertStringEquals(t, got.Username, userList.Users[0].Username, "new user should be first in list")

	userList, e = webClient.Users(time.Now().AddDate(1, 0, 0), time.Now().AddDate(1, 0, 0), cfg)
	if e != nil {
		t.Fatalf("unexpected error : " + e.Error())
	}
	assertTrue(t, userList.TotalHits > 0, "Expected some users but was 0")
	// no users created 10 years ago
	userList2, _ := webClient.Users(time.Time{}, time.Now().AddDate(-10, 0, 0), cfg)
	assertIntEquals(t, 0, userList2.TotalHits, "")
}

func TestGroupNew(t *testing.T) {
	// given a PI user
	userPiPost := createRandomUser(Pi)
	var err error
	var user *UserInfo
	user, err = webClient.UserNew(userPiPost)
	if err != nil {
		Log.Error(err)
	}
	groups, _ := webClient.Groups()
	initialGroupCount := len(groups.Groups)
	//create a group
	var userGroupPosts []UserGroupPost = make([]UserGroupPost, 0, 5)
	userGroupPosts = append(userGroupPosts, UserGroupPost{user.Username, "PI"},
		UserGroupPost{"sysadmin1", "DEFAULT"})
	groupPost, err := GroupPostNew("groupname", userGroupPosts)
	var group *GroupInfo
	group, err = webClient.GroupNew(groupPost)
	assertNil(t, err, "")
	assertNotNil(t, group, "")
	assertStringEquals(t, "groupname", group.Name, "")
	assertIntEquals(t, 2, len(group.Members), "")

	groups, _ = webClient.Groups()
	assertIntEquals(t, initialGroupCount+1, len(groups.Groups), "")
}

func createRandomUser(userRole UserRoleType) *UserPost {
	uname := randomAlphanumeric(8)
	pwd := randomAlphanumeric(8)
	var email Email = Email(fmt.Sprintf("%s@somewhere.com", uname))
	firstName := randomAlphanumeric(3)
	lastName := randomAlphanumeric(8)
	userBuilder := UserPostBuilder{}
	userPost, _ := userBuilder.Affiliation("somewhere").Username(uname).Password(pwd).Email(email).FirstName(firstName).LastName(lastName).Role(userRole).Build()
	return userPost
}
