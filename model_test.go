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
	userpost, err = builder.username("user1234").password("secret23").firstName("first").lastName("last").email("a@b.com").role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, userpost, "UserPost was nil")

	// key, affiliation are optional
	builder = &UserPostBuilder{}
	userpost, err = builder.username("user1234").password("secret23").firstName("first").lastName("last").email("a@b.com").role(user).build()
	assertNotNil(t, userpost, "UserPost was nil")

	builder = &UserPostBuilder{}
	userpost, err = builder.build()
	assertNotNil(t, err, "error was nil")

	builder = &UserPostBuilder{}
	userpost, err = builder.username("abc").password("secret23").firstName("first").lastName("last").email("a@b.com").role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, err, "error was nil")

	builder = &UserPostBuilder{}
	tooShortPwd := "secret2"
	userpost, err = builder.username("user1234").password(tooShortPwd).firstName("first").lastName("last").email("a@b.com").role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, userpost, "error was nil")

	builder = &UserPostBuilder{}
	tooShortEmail := Email("@")
	userpost, err = builder.username("user1234").password(tooShortPwd).firstName("first").lastName("last").email(tooShortEmail).role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	assertNotNil(t, userpost, "error was nil")

}

func TestGroupPost(t *testing.T) {
	var groupPost *GroupPost
	var err error
	// success
	var userGroupPosts []UserGroupPost = make([]UserGroupPost, 0, 5)
	userGroupPosts = append(userGroupPosts, UserGroupPost{"username1", "PI"})
	groupPost, err = GroupPostNew("groupname", userGroupPosts)
	assertNotNil(t, groupPost, "Group post was nil")
	fmt.Println(userGroupPosts)
	assertIntEquals(t, 1, len(groupPost.Members), "")

	// name required
	groupPost, err = GroupPostNew("", userGroupPosts)
	assertNotNil(t, err, "expected error,  was nil")

	//at least 1 group member required
	groupPost, err = GroupPostNew("name", make([]UserGroupPost, 0, 5))
	assertNotNil(t, err, "expected error,  was nil")

	groupPost, err = GroupPostNew("name", []UserGroupPost{UserGroupPost{"some user", "DEFAULT"}})
	assertNotNil(t, err, "expected error,  was nil")

}
func TestActivityQueryBuilder(t *testing.T) {
	var err error
	var q *ActivityQuery 
	var builder  *ActivityQueryBuilder
	builder = &ActivityQueryBuilder{}
	// valid global id
	q, err = builder.Oid("GL1234").Build()
	assertNotNil(t, q, "query should not be nil")
	assertNil(t, err, "err should be  nil")

	// invalid global id
	q = nil
	q, err = builder.Oid("GL???4").Build()
	fmt.Println(q)
	assertNotNil(t, err, "err should not be  nil")

	domain := "RECORD"
	action := "READ"
	action2 := "COPY"
	user := "bob"
	builder = &ActivityQueryBuilder{}
	q,err = builder.Domain(domain).Action(action).Action(action2).User(user).Build()
	assertStringEquals(t, "bob", q.Users[0], "")
	assertStringEquals(t, "READ", q.Actions[0], "")
	assertStringEquals(t, "COPY", q.Actions[1], "")
	assertStringEquals(t, "RECORD", q.Domains[0], "")
	assertTrue(t, len(q.Oid)==0, "OID should be empty")
}

