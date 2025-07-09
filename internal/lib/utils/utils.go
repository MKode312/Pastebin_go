package utils

import (
	"fmt"
	"strings"
)

type FileDataType struct {
	FileName string
	Data     []byte
}

type OperationError struct {
	ObjectID string
	Error    error
}

// Вытаскиваем имя файла из URL
func ExtractFilenameFromURL(url string) string {
	// Можно использовать стандартную библиотеку или просто взять последний сегмент URL
	lastSlash := -1
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			lastSlash = i
			break
		}
	}
	filename := url[lastSlash+1:]

	// Можно дополнительно очистить имя файла или оставить как есть
	return filename
}

func ExtractFilenameFromString(input string) (string, error) {
    questionMarkIdx := strings.Index(input, "?")
    if questionMarkIdx == -1 {
        return "", fmt.Errorf("char '?' was not found in string input filename")
    }
    filename := input[:questionMarkIdx]
    return filename, nil
}
