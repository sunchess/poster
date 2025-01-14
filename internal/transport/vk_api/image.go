package vk_api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vk_poster/internal/utils"

	"github.com/SevereCloud/vksdk/v3/api/params"
)

func NewImageUploader(postDir string, current_client *VkApi, postTime *int64) *Uploader {
	return &Uploader{
		Client:          current_client.Client,
		PostDir:         postDir,
		MediaDir:        filepath.Join(postDir, "media"),
		MessagePath:     filepath.Join(postDir, "message.txt"),
		GroupId:         current_client.GroupId,
		UploaderGroupID: -current_client.GroupId, //positive group id
		PostTime:        postTime,
	}
}

func (u *Uploader) UploadOnlyImageOnWall() error {
	imagePaths := u.getImagePaths()
	log.Printf("Found image files: %v\n", imagePaths)

	_, err := u.uploadImage(imagePaths[0])

	return err
}

func (u *Uploader) UploadPostWithImages() error {
	// get message from message.txt
	message, _ := utils.CleanText(u.MessagePath)

	// Get image file paths
	imagePaths := u.getImagePaths()
	if len(imagePaths) == 0 {
		return fmt.Errorf("no image files found in media directory")
	}
	log.Printf("Found image files: %v\n", imagePaths)

	// Upload images and get attachments
	var attachments []string
	for _, imagePath := range imagePaths {
		attachment, err := u.uploadImage(imagePath)
		if err != nil {
			return err
		}
		attachments = append(attachments, attachment)
	}

	// Post message with attachments
	if err := u.postMessageWithAttachments(string(message), attachments); err != nil {
		return err
	}

	return nil
}

// ****** Private functions ******

func (u *Uploader) getImagePaths() []string {
	files, err := os.ReadDir(u.MediaDir)
	if err != nil {
		log.Fatalf("failed to read media directory: %v", err)
	}

	var imagePaths []string
	for _, file := range files {
		if strings.EqualFold(filepath.Ext(file.Name()), ".jpg") || strings.EqualFold(filepath.Ext(file.Name()), ".jpeg") || strings.EqualFold(filepath.Ext(file.Name()), ".png") {
			imagePaths = append(imagePaths, filepath.Join(u.MediaDir, file.Name()))
		}
	}

	return imagePaths
}

func (u *Uploader) uploadImage(imagePath string) (string, error) {
	// Check the image file in parent function
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %w", err)
	}
	defer imageFile.Close()

	saveResponse, err := u.Client.UploadGroupWallPhoto(u.UploaderGroupID, imageFile)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	attachment := fmt.Sprintf("photo%d_%d", saveResponse[0].OwnerID, saveResponse[0].ID)
	return attachment, nil
}

func (u *Uploader) postMessageWithAttachments(message string, attachments []string) error {
	log.Println("Posting message with attachments...")

	wallPostParams := params.NewWallPostBuilder().
		OwnerID(u.GroupId).
		Message(message).
		Attachments(attachments).
		FromGroup(true).
		PublishDate(int(*u.PostTime)).
		Params

	_, err := u.Client.WallPost(wallPostParams)
	if err != nil {
		return fmt.Errorf("failed to post on wall: %w", err)
	}
	log.Println("Posted message with attachments successfully")
	return nil
}
