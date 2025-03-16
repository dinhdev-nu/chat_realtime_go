package main

import "github.com/dinhdev-nu/realtime_auth_go/internal"



func main() {

	r:= internal.Run()
	
	r.Run("localhost:3010")

}