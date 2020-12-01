package util

import (
	"github.com/rs/zerolog"
	"os"
	"path"
)

type AudienceLog struct {
	zerolog.Logger
	FileHandler	*os.File
}

func NewLog() (*AudienceLog, error) {
	filepath, err := GetHomePath()
	if err != nil {
		return nil, err
	}

	filepath = path.Join(filepath, "audience.log")
	file, err := os.OpenFile(filepath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(file).With().Timestamp().Logger()

	return &AudienceLog{
		Logger:      logger,
		FileHandler: file,
	}, nil
}
