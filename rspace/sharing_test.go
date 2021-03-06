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
	list, _ := webClient.ShareList(doc.Name, NewRecordListingConfig())
	initialCount := list.TotalHits
	grpShare := GroupShare{Id: group.Id, Permission: "edit", SharedFolderId: group.SharedFolderId}
	idsToShare := make([]int, 0)
	idsToShare = append(idsToShare, doc.Id)
	grps := make([]GroupShare, 0)
	grps = append(grps, grpShare)
	sharePost := SharePost{Groups: grps, ItemsToShare: idsToShare, Users: []UserShare{}}
	shared, err := webClient.Share(&sharePost)
	fmt.Println(shared)
	if err != nil {
		Log.Warning(err.Error())
	}
	assertIntEquals(t, doc.Id, shared.ShareInfos[0].ItemId, "")
	assertStringEquals(t, doc.Name, shared.ShareInfos[0].ItemName, "")

	// now list, should
	list, _ = webClient.ShareList(doc.Name, NewRecordListingConfig())
	assertIntEquals(t, initialCount+1, list.TotalHits, "")
	assertIntEquals(t, initialCount+1, len(list.Shares), "")

	// now unshare
	deleted, err := webClient.Unshare(shared.ShareInfos[0].Id)
	assertTrue(t, deleted, "should have been unshared")
	// original number of shared items
	list, _ = webClient.ShareList(doc.Name, NewRecordListingConfig())
	assertIntEquals(t, initialCount, list.TotalHits, "")
	assertIntEquals(t, initialCount, len(list.Shares), "")

}
