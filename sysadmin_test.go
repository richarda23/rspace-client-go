package rspace

import (
	"fmt"
	"testing"
	"time"
	//"errors"
)

var sysads *SysadminService = &SysadminService{
	BaseService: BaseService{
		Delay: time.Duration(100) * time.Millisecond}}


func TestUserNew(t *testing.T) {

	// given
	userPost := createRandomUser(pi)
	var got = sysads.UserNew(userPost)
	if got.Id == 0 {
		fail(t, "Id was nil but should be set")
	}
	assertStringEquals(t, userPost.Username, got.Username,"")
}
func createRandomUser(userRole UserRoleType) *UserPost {
	uname := randomAlphanumeric(8)
	pwd := randomAlphanumeric(8)
	var email Email = Email(fmt.Sprintf("%s@somewhere.com", uname))
	firstName := randomAlphanumeric(3)
	lastName := randomAlphanumeric(8)
	userBuilder := UserPostBuilder{}
	userPost,_ := userBuilder.username(uname).password(pwd).email(email).firstName(firstName).lastName(lastName).role(userRole).build()
	return userPost
}
func TestGroupNew(t *testing.T) {

	// given a PI user
	userPiPost := createRandomUser(pi)
	var user *UserInfo = sysads.UserNew(userPiPost)

	//create a group
	var userGroupPosts []UserGroupPost = make ([]UserGroupPost,0,5)
	userGroupPosts = append(userGroupPosts, UserGroupPost{user.Username, "PI"})
	groupPost,err := GroupPostNew("groupname", userGroupPosts)
	var group *GroupInfo;
	group,err = sysads.GroupNew(groupPost)
	assertNil(t, err, "")
	assertNil(t, group, "")

}

