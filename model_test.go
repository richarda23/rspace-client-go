package rspace

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestUserPost(t *testing.T) {
	builder := &UserPostBuilder{}
	userpost :=builder.username("user1234").password("secret23").firstName("first").lastName("last").email("a@b.com").role(user).affiliation("u-somewhere").apiKey("abcdefg").build()
	
	fmt.Println(userpost)
	json, _ := json.Marshal(userpost)
	fmt.Println(string(json))
}
