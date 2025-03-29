package random

import (
	"fmt"
	"math/rand"
	"time"
)

// func genarateOp 6 number

func CreateOtp() string {
	src := rand.NewSource(time.Now().UnixNano()) // tao ra 1 so ngau nhien dua tren thoi gian
	rand := rand.New(src)
	otp := fmt.Sprintf("%06v", rand.Intn(1000000))
	return otp
}
