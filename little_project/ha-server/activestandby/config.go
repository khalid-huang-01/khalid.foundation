package activestandby

import "sync"

var Config Configure
var once sync.Once

type Configure struct {
	Checker *ReadyzAdaptor
}

func InitConfigure() {
	once.Do(func() {
		Config = Configure{}
	})
}