package main

import (
	"flag"
	"log"
	"vk_poster/internal/services"
)

func main() {
	log.Println("Start posting...")

	scope := flag.String("scope", "vk", "Platform to post")
	limit := flag.Int("limit", 2, "Number of posts at one run")
	postPublishGap := flag.Int("post_gap", 1800, "Gap between posts in seconds")

	flag.Parse()

	services.Posting(*scope, *limit, *postPublishGap)
}
