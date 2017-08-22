package main

import (
	"fmt"
	"github.com/jhamon/uaa/info"
)

func main() {
	client := info.UaaClient{"https://login.run.pivotal.io"}
	status := info.Health(client)
	fmt.Println("Status was:" + status)

	client2 := info.UaaClient{"https://safdsfaslogin.run.pivotal.io"}
	status2 := info.Health(client2)
	fmt.Println("Status was:" + status2)
}