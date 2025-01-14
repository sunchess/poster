package services

import (
	"log"
	"os"
	"sync/atomic"
	"time"
	"vk_poster/internal/transport/vk_api"
	"vk_poster/internal/utils"
)

type VkPostingService struct {
	PostingService
	Connection *vk_api.VkApi
}

func (vp VkPostingService) ProcessingPost(postDir string) {
	postTime := atomic.LoadInt64(vp.PostingService.PostPublishTime)

	//next post time in PostPublishGap
	defer func(postTime *int64) {
		*postTime = *postTime + int64(vp.PostingService.PostPublishGap)
		atomic.StoreInt64(vp.PostingService.PostPublishTime, *postTime)
	}(&postTime)

	api_client := vp.Connection

	//continue if post text file exist and ends with POSTFIX
	hasSign := utils.IsLastLineContainsString(postDir, os.Getenv("NOT_ADV_POST_SIGN_POSTFIX"))
	if !hasSign {
		log.Printf("There is no POSTFIX in the last line of the text file for post: %s", postDir)
		vp.PostingService.Db.SetPostedPostDir(postDir, vp.PostingService.Scope)
		return
	}

	//if media directory has images then upload post with images
	if utils.HasImages(postDir) {
		for i := 0; i < 2; i++ {
			vk_api_image_uploader := vk_api.NewImageUploader(postDir, api_client, &postTime)
			err := vk_api_image_uploader.UploadPostWithImages()

			if err != nil {
				// make retry
				log.Printf("Failed to upload post with image: %v", err)
				time.Sleep(10 * time.Second)
			} else {
				vp.PostingService.Db.SetPostedPostDir(postDir, vp.Scope)
				break
			}
		}
		// if post has text then upload post with text
	} else if utils.HasText(postDir) {
		err := api_client.PostWithText(postDir, api_client)
		if err != nil {
			log.Fatalf("Failed to upload post with text: %v", err)
		} else {
			vp.PostingService.Db.SetPostedPostDir(postDir, vp.Scope)
		}
		// if post has video then upload post with video
	} else if utils.HasVideo(postDir) {
		vk_api_video_uploader := vk_api.NewVideoUploader(postDir, api_client, &postTime)
		err := vk_api_video_uploader.UploadPostWithMedia()

		if err != nil {
			log.Fatalf("Failed to upload post with video: %v", err)
		} else {
			vp.PostingService.Db.SetPostedPostDir(postDir, vp.Scope)
		}
		// if only image without text then upload only image
	} else if utils.HasOnlyImage(postDir) {
		vk_api_image_uploader := vk_api.NewImageUploader(postDir, api_client, &postTime)
		err := vk_api_image_uploader.UploadOnlyImageOnWall()

		if err != nil {
			log.Fatalf("Failed to upload only image: %v", err)
		} else {
			vp.PostingService.Db.SetPostedPostDir(postDir, vp.Scope)
		}
	} else {
		log.Printf("Post directory %s does not persist for posting", postDir)
	}
}
