package upload

import (
	"context"
	"log"
	"os"
	"vk_poster/internal/database"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// VideoMetadata содержит информацию о загружаемом видео.
type VideoMetadata struct {
	Title       string
	Description string
	CategoryID  string
	Privacy     string
}

type UploadVideoService struct {
	Db      *database.DBConnection
	Scope   string
	PostDir string
}

// UploadVideo загружает видео на YouTube.
func (up *UploadVideoService) UploadVideo(ctx context.Context, tokenSource oauth2.TokenSource, videoPath string, metadata VideoMetadata) {
	client := oauth2.NewClient(ctx, tokenSource)

	// Создание YouTube сервиса
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	// Настройка метаданных видео
	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       metadata.Title,
			Description: metadata.Description,
			CategoryId:  metadata.CategoryID,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: metadata.Privacy,
		},
	}

	// Открытие файла видео
	file, err := os.Open(videoPath)
	if err != nil {
		log.Fatalf("Error opening video file: %v", err)
	}
	defer file.Close()

	// Загрузка видео
	call := service.Videos.Insert([]string{"snippet", "status"}, video)
	response, err := call.Media(file).Do()
	if err != nil {
		log.Fatalf("Error uploading video: %v", err)
	}

	log.Printf("Video uploaded successfully! ID: %s\n", response.Id)

	up.Db.SetPostedPostDir(up.PostDir, up.Scope)
}
