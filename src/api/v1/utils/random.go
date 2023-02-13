package utils

import uuid2 "github.com/google/uuid"

func Uuid() string {
	uuid := uuid2.New()
	return uuid.String()
}
