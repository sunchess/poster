package vk_api

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v3/api"
)

type VkApi struct {
	AccessToken string
	Client      *api.VK
	GroupId     int
}

type Uploader struct {
	Api                *VkApi
	Client             *api.VK
	PostDir            string
	MediaDir           string
	MessagePath        string
	GroupId            int
	UploaderGroupID    int
	ProcessedVideoPath string
	PostTime           *int64
}

// ****** Public functions ******
func NewVkApi() *VkApi {
	groupId := convertToGroupId(os.Getenv("VK_GROUP_ID"))

	return &VkApi{
		AccessToken: os.Getenv("VK_ACCESS_TOKEN"),
		GroupId:     groupId,
	}
}

func (a *VkApi) Connect() {
	a.Client = api.NewVK(a.AccessToken)
	//debug a.Client
	fmt.Println(a.Client)
}

func convertToGroupId(groupId string) int {
	convertedGroupId, _ := strconv.Atoi(groupId)
	currentGroupId := -(convertedGroupId)

	return currentGroupId
}
