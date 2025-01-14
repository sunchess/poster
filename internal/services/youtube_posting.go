package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"vk_poster/internal/transport/youtube_api/auth"
	"vk_poster/internal/transport/youtube_api/upload"
	"vk_poster/internal/utils"
)

type YoutubePostingService struct {
	PostingService
}

func (yp YoutubePostingService) ProcessingPost(postDir string) {
	if !utils.HasVideo(postDir) {
		// Загрузка видео на YouTube
		log.Printf("Post dir %s does not have video file YouTube...", postDir)
		yp.Db.SetPostedPostDir(postDir, yp.Scope)
		return
	}

	ctx := context.Background()

	type DBConnection string
	const dbConnect DBConnection = "dbConnect"

	ctx = context.WithValue(ctx, dbConnect, yp.Db)

	// Настройки
	scope := "https://www.googleapis.com/auth/youtube.upload"
	videoPath := utils.ProcessedVideoPath(postDir)

	postText, ok := utils.CleanText(utils.MessagePath(postDir))

	title := "Short от Ани дизайнера"
	description := "Short от Ани дизайнера"

	if ok {
		// get first 3 words from the text
		title = strings.Join(strings.Fields(postText)[:3], " ")
		description = postText
	}

	metadata := upload.VideoMetadata{
		Title:       title,
		Description: description,
		CategoryID:  "22",     // Категория "People & Blogs"
		Privacy:     "public", // Или "private", "unlisted"
	}

	serviceAccountJsonFile := os.Getenv("SERVICE_ACCOUNT_JSON_FILE")

	// Получение источника токенов
	tokenSource, err := auth.GetTokenSource(ctx, serviceAccountJsonFile, scope)
	if err != nil {
		log.Fatalf("Failed to get token source: %v", err)
	}

	fmt.Println("Token source: ", tokenSource)
	fmt.Println("Video path: ", videoPath)
	fmt.Println("Metadata: ", metadata)

	// Загрузка видео
	uploadService := upload.UploadVideoService{Db: yp.Db, Scope: yp.Scope, PostDir: postDir}
	uploadService.UploadVideo(ctx, tokenSource, videoPath, metadata)
}
