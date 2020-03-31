package rspace

import (
	"testing"
	//"errors"
	"fmt"
	//"encoding/json"
)

//validate User creation
func TestUserPost(t *testing.T) {
	var builder = &UserPostBuilder{}
	var err error = nil
	var userpost *UserPost = nil
	userpost,err =builder.username("user1234").password("secret23").firstName("first").lastName("last").email("a@b.com").role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, userpost, "UserPost was nil")

	// key, affiliation are optional
	builder = &UserPostBuilder{}
	userpost,err =builder.username("user1234").password("secret23").firstName("first").lastName("last").email("a@b.com").role(user).build()
	assertNotNil(t, userpost, "UserPost was nil")

	builder  = &UserPostBuilder{}
	userpost,err = builder.build()
	assertNotNil(t, err, "error was nil")

	builder  = &UserPostBuilder{}
	userpost,err =  builder.username("abc").password("secret23").firstName("first").lastName("last").email("a@b.com").role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, err, "error was nil")

	builder  = &UserPostBuilder{}
	tooShortPwd := "secret2"
	userpost,err =builder.username("user1234").password(tooShortPwd).firstName("first").lastName("last").email("a@b.com").role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, userpost, "error was nil")

	builder  = &UserPostBuilder{}
	tooShortEmail := Email("@")
	userpost,err =builder.username("user1234").password(tooShortPwd).firstName("first").lastName("last").email(tooShortEmail).role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, userpost, "error was nil")
	
}

func TestGroupPost(t *testing.T) {
	var groupPost *GroupPost;
	var err error;
	// success
	var userGroupPosts []UserGroupPost = make ([]UserGroupPost,0,5)
	userGroupPosts = append(userGroupPosts, UserGroupPost{"username1", "PI"})
	groupPost,err = GroupPostNew("groupname", userGroupPosts)
	assertNotNil(t, groupPost, "Group post was nil")
	fmt.Println(userGroupPosts)
	assertIntEquals(t,1,len(groupPost.Members),"")

	// name required
	groupPost,err = GroupPostNew("", userGroupPosts)
	assertNotNil(t, err, "expected error,  was nil")

	//at least 1 group member required
	groupPost,err = GroupPostNew("name", make ([]UserGroupPost,0,5))
	assertNotNil(t, err, "expected error,  was nil")

	groupPost,err = GroupPostNew("name", []UserGroupPost{ UserGroupPost{"some user", "DEFAULT"} })
	assertNotNil(t, err, "expected error,  was nil")

}
