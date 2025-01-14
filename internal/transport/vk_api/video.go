package vk_api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"vk_poster/internal/utils"

	"github.com/SevereCloud/vksdk/v3/api/params"
)

func NewVideoUploader(postDir string, current_client *VkApi, postTime *int64) *Uploader {
	return &Uploader{
		Client:             current_client.Client,
		PostDir:            postDir,
		MediaDir:           filepath.Join(postDir, "media"),
		MessagePath:        filepath.Join(postDir, "message.txt"),
		GroupId:            current_client.GroupId,
		UploaderGroupID:    -current_client.GroupId, //positive group id,
		ProcessedVideoPath: filepath.Join(postDir, "media", "processed.mp4"),
		PostTime:           postTime,
	}
}

func (v *Uploader) UploadPostWithMedia() error {
	// get message from message.txt
	message, _ := utils.CleanText(v.MessagePath)

	// get video file path
	videoPath := v.getVideoPath()
	if videoPath == "" {
		return fmt.Errorf("no video file found in media directory")
	}
	log.Printf("Found video file: %s\n", videoPath)

	// upload video, get attachment
	attachment, err := v.uploadVideo(videoPath)
	if err != nil {
		return err
	}

	// post message with attachment
	err = v.postMessageWithAttachment(string(message), attachment)
	if err != nil {
		return err
	}

	return nil
}

// ****** Private functions ******

func (v *Uploader) getVideoPath() string {
	_, err := os.Stat(v.ProcessedVideoPath)
	if err != nil {
		log.Fatalf("processed video file not found for video %s", v.ProcessedVideoPath)
	}

	return v.ProcessedVideoPath
}

func (v *Uploader) uploadVideo(videoPath string) (string, error) {
	// check the video file in parent function
	videoFile, _ := os.Open(videoPath)
	defer videoFile.Close()

	// get video upload server
	current_params := params.NewVideoSaveBuilder().
		GroupID(v.UploaderGroupID).
		Repeat(true).
		Description(os.Getenv("VIDEO_DESCRIPTION")).
		Params
	videoResponse, err := v.Client.UploadVideo(current_params, videoFile)

	if err != nil {
		return "", fmt.Errorf("failed to upload video: %w", err)
	}

	log.Println("Uploaded video successfully")

	attachment := fmt.Sprintf("video%d_%d", videoResponse.OwnerID, videoResponse.VideoID)
	return attachment, nil
}

func (v *Uploader) postMessageWithAttachment(message string, attachment string) error {
	log.Println("Posting message with attachment...")

	wallPostParams := params.NewWallPostBuilder().
		OwnerID(v.GroupId).
		Message(message).
		Attachments([]string{attachment}).
		FromGroup(true).
		PublishDate(int(*v.PostTime)).
		Params

	log.Printf("WallPostParams: %v\n", wallPostParams)

	_, err := v.Client.WallPost(wallPostParams)
	if err != nil {
		return fmt.Errorf("failed to post on wall: %w", err)
	}
	log.Println("Posted message with attachment successfully")
	return nil
}
