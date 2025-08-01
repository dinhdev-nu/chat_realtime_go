package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GenerateUUIDToken(userId int64) string {
	newUUID := uuid.New()
	cvNewUUID := strings.ReplaceAll(newUUID.String(), "-", "")
	return strconv.Itoa(int(userId)) + "token" + cvNewUUID
}

func GetUserIdFromUUIDToken(uuidToken string) int64 {
	id := strings.Split(uuidToken, "token")[0]
	userId, _ := strconv.Atoi(id)
	return int64(userId)
}
