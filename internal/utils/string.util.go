package utils

import (
	"math/rand/v2"
	"sort"
	"strings"
)

// import "strings"

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
