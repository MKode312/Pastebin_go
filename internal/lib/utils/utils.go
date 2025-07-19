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

func ExtractFilenameFromURL(url string) string {
	lastSlash := -1
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			lastSlash = i
			break
		}
	}
	filename := url[lastSlash+1:]

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
