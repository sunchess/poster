package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
	"vk_poster/internal/database"
	"vk_poster/internal/transport/vk_api"

	"github.com/joho/godotenv"
)

type PostingService struct {
	Scope           string
	Db              *database.DBConnection
	Limit           int
	PostPublishTime *int64
	PostPublishGap  int
}

func NewPostingService(scope string, db *database.DBConnection, limit int, postPublishGap int) *PostingService {
	// time now + 3 minutes
	postTime := time.Now().Unix() + 180
	atomic.StoreInt64(&postTime, postTime)

	return &PostingService{
		Scope:           scope,
		Db:              db,
		Limit:           limit,
		PostPublishTime: &postTime,
		PostPublishGap:  postPublishGap,
	}
}

type PostingServiceInterface interface {
	ProcessingPost(postDir string)
}

func setEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetGateway(postingService *PostingService) (PostingServiceInterface, error) {
	switch postingService.Scope {
	case "vk":
		api_client := vk_api.NewVkApi()
		api_client.Connect()

		return VkPostingService{PostingService: *postingService, Connection: api_client}, nil
	case "youtube":
		return YoutubePostingService{*postingService}, nil
	default:
		return nil, errors.New("unknown scope")
	}
}

func Posting(scope string, limit int, postPublishGap int) {
	setEnv()

	db_config := database.DBcredentials{DbPath: os.Getenv("DB_PATH")}

	// connect to database
	db, err := database.NewDBConnection(&db_config)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	postingService := NewPostingService(scope, db, limit, postPublishGap)

	directories, err := postingService.Db.GetPostDirectories(postingService.Scope, postingService.Limit)

	if err != nil {
		log.Fatalf("Failed to get post directories: %v", err)
	}

	gateway, err := GetGateway(postingService)
	if err != nil {
		log.Fatalf("Failed to get gateway: %v", err)
	}

	for _, dir := range directories {
		fmt.Printf("Processing post: %s\n", dir)
		gateway.ProcessingPost(dir)
		time.Sleep(5 * time.Second)
	}
}
