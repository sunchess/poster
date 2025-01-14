package vk_api

import (
	"fmt"
	"log"
	"path/filepath"
	"vk_poster/internal/utils"

	"github.com/SevereCloud/vksdk/v3/api/params"
)

func (a *VkApi) PostWithText(postDir string, currentSession *VkApi) error {
	messagePath := filepath.Join(postDir, "message.txt")

	message, _ := utils.CleanText(messagePath)

	wallPostParams := params.NewWallPostBuilder().
		OwnerID(currentSession.GroupId).
		Message(string(message)).
		FromGroup(true).
		Params

	_, err := a.Client.WallPost(wallPostParams)
	if err != nil {
		return fmt.Errorf("failed to post on wall: %w", err)
	}
	log.Println("Posted message successfully")

	return nil
}
