package utils

import (
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"
)

// import "strings"

func GetToken(header string) string {
	// split header by space
	token := strings.Split(header, " ")
	if len(token) < 2 {
		return ""
	}
	// return token
	return token[1]

}

func GennarateUserName(email string) string {
	// split email by @ last
	nameEmail := email[:strings.LastIndex(email, "@")]
	// ramdom 6 char in email
	var id [6]int
	for i := 0; i < 6; i++ {
		id[i] = rand.IntN(len(nameEmail))
	}
	sort.Ints(id[:])
	var userName string
	for i := 0; i < len(id); i++ {
		userName += string(nameEmail[id[i]])
	}
	return userName
}

// convert string to uint64
func StringToUint64(s string) uint64 {
	n, _ := strconv.ParseUint(s, 10, 64)
	return n
}

// convert string to int64
func StringToInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
