package main

import (
	"math/rand"
	"time"

	"khalid.foundation/haserver/activestandby"
)

var (
	LeaderElect = true
)


func main() {
	if LeaderElect {
		rand.Seed(time.Now().Unix())
		go activestandby.StartHelperHTTPServer()
		checker := activestandby.NewLeaderReadyzAdaptor(time.Second * 20)
		activestandby.Config.Checker = checker
		activestandby.Run(checker)
		return
	}
	//activestandby.StartMainHTTPServer()
}