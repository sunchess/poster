package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func ProcessedVideoPath(postDir string) string {
	return filepath.Join(MediaDir(postDir), "processed.mp4")
}

func MediaDir(postDir string) string {
	return filepath.Join(postDir, "media")
}

func MessagePath(postDir string) string {
	return filepath.Join(postDir, "message.txt")
}

func hasTextFile(postDir string) bool {
	messageFilePath := filepath.Join(postDir, "message.txt")
	if _, err := os.Stat(messageFilePath); os.IsNotExist(err) {
		if err != nil {
			return false
		}
	}
	return true
}
func HasImages(postDir string) bool {
	files, err := os.ReadDir(MediaDir(postDir))
	if err != nil {
		return false
	}

	for _, file := range files {
		if hasTextFile(postDir) && slices.Contains([]string{".jpg", ".jpeg", ".png"}, strings.ToLower(filepath.Ext(file.Name()))) {
			return true
		}
	}
	return false
}

// check if media directory has files with extensions .mp4
func HasVideo(postDir string) bool {
	files, err := os.ReadDir(MediaDir(postDir))
	if err != nil {
		return false
	}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		fmt.Printf("file: %v\n", ext)
		if ext == ".mp4" {
			return true
		}
	}
	return false
}

// not exist media directory
func HasText(postDir string) bool {
	_, err := os.Stat(MediaDir(postDir))
	return os.IsNotExist(err)
}

// check if media directory has only one image file
func HasOnlyImage(postDir string) bool {
	files, err := os.ReadDir(MediaDir(postDir))
	if err != nil {
		return false
	}

	if !hasTextFile(postDir) && len(files) == 1 && slices.Contains([]string{".jpg", ".jpeg", ".png"}, strings.ToLower(filepath.Ext(files[0].Name()))) {
		return true
	}
	return false
}

func IsLastLineContainsString(postDir, searchString string) bool {
	file, err := os.ReadFile(MessagePath(postDir))

	if err != nil {
		return false
	}

	// Разделить содержимое файла на строки
	lines := strings.Split(string(file), "\n")

	// Проверить последнюю строку
	if len(lines) == 0 {
		return false
	}
	lastLine := lines[len(lines)-1]

	return strings.Contains(lastLine, searchString)
}

func CleanText(path string) (string, bool) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read file: %v", err)
		return "", false
	}
	return strings.ReplaceAll(string(file), "**", ""), true
}
