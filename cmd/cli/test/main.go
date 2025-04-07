package main

import (
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/internal/utils"
)

func main() {
	fmt.Println(utils.GenerateUUIDToken(123))
}
