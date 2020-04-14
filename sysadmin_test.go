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
	userPost := createRandomUser(pi)
	got, err := webClient.SysadminS.UserNew(userPost)
	if err != nil {
		Log.Error(err)
	}
	if got.Id == 0 {
		fail(t, "Id was nil but should be set")
	}
	assertStringEquals(t, userPost.Username, got.Username, "")
}
func TestUsers(t *testing.T) {
	// all users were created before a time in the future
	userList, e := webClient.SysadminS.Users(time.Now().AddDate(1, 0, 0), time.Now().AddDate(1, 0, 0))
	if e != nil {
		fmt.Println(e)
	}
	assertTrue(t, userList.TotalHits > 0, "Expected some users but was 0")
	// no users created 10 years ago
	userList2, _ := webClient.SysadminS.Users(time.Time{}, time.Now().AddDate(-10, 0, 0))
	assertIntEquals(t, 0, userList2.TotalHits, "")
}

func TestGroupNew(t *testing.T) {
	// given a PI user
	userPiPost := createRandomUser(pi)
	var err error
	var user *UserInfo
	user, err = webClient.SysadminS.UserNew(userPiPost)
	if err != nil {
		Log.Error(err)
	}
	//create a group
	var userGroupPosts []UserGroupPost = make([]UserGroupPost, 0, 5)
	userGroupPosts = append(userGroupPosts, UserGroupPost{user.Username, "PI"})
	groupPost, err := GroupPostNew("groupname", userGroupPosts)
	var group *GroupInfo
	group, err = webClient.SysadminS.GroupNew(groupPost)
	if err != nil {
		Log.Error(err)
	}
	assertNil(t, err, "")
	assertNotNil(t, group, "")
	assertStringEquals(t, "groupname", group.Name, "")
	assertIntEquals(t, 1, len(group.Members), "")
}

func createRandomUser(userRole UserRoleType) *UserPost {
	uname := randomAlphanumeric(8)
	pwd := randomAlphanumeric(8)
	var email Email = Email(fmt.Sprintf("%s@somewhere.com", uname))
	firstName := randomAlphanumeric(3)
	lastName := randomAlphanumeric(8)
	userBuilder := UserPostBuilder{}
	userPost, _ := userBuilder.affiliation("somewhere").username(uname).password(pwd).email(email).firstName(firstName).lastName(lastName).role(userRole).build()
	return userPost
}
