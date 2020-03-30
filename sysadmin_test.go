package rspace

import (
	"fmt"
	"testing"
	"time"
)

var sysads *SysadminService = &SysadminService{
	BaseService: BaseService{
		Delay: time.Duration(100) * time.Millisecond}}


func TestUserNew(t *testing.T) {

	// given
	uname := randomAlphanumeric(8)
	pwd := randomAlphanumeric(8)
	var email Email = Email(fmt.Sprintf("%s@somewhere.com", uname))
	firstName := "Bob"
	lastName := "Smith"
	userBuilder := UserPostBuilder{}
	userPost := userBuilder.affiliation("somwhere").username(uname).password(pwd).email(email).firstName(firstName).lastName(lastName).role(pi).build()
	Log.Info(userPost)
	var got = sysads.UserNew(userPost)
	if got.Id == 0 {
		fail(t, "Id was nill but should be set")
	}
	assertStringEquals(t, uname, got.Username,"")
}

