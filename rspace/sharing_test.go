package rspace

import (
	"fmt"
	"testing"
)

// Requires sysadmin permission to create the group
func TestSharingNew(t *testing.T) {
	// given a PI user
	var err error
	var pi *UserInfo
	userPiPost := createRandomUser(Pi)
	pi, err = webClient.UserNew(userPiPost)
	if err != nil {
		t.Fatalf(err.Error())
	}
	//create a group
	var userGroupPosts []UserGroupPost = make([]UserGroupPost, 0, 5)
	userGroupPosts = append(userGroupPosts, UserGroupPost{pi.Username, "PI"},
		UserGroupPost{"sysadmin1", "DEFAULT"})
	groupPost, err := GroupPostNew("groupname", userGroupPosts)
	var group *GroupInfo
	group, err = webClient.GroupNew(groupPost)
	assertNotNil(t, group, "")

	doc, _ := webClient.NewEmptyBasicDocument("toShare", "")
	grpShare := GroupShare{Id: group.Id, Permission: "edit", SharedFolderId: group.SharedFolderId}
	idsToShare := make([]int, 0)
	idsToShare = append(idsToShare, doc.Id)
	grps := make([]GroupShare, 0)
	grps = append(grps, grpShare)
	sharePost := SharePost{Groups: grps, ItemsToShare: idsToShare, Users: []UserShare{}}
	shared, err := webClient.Share(&sharePost)
	if err != nil {
		Log.Warning(err.Error())
	}
	fmt.Println(shared.ShareInfos[0])
}
