package main

import (
	"fmt"
	"sync"

	r "github.com/dancannon/gorethink"
)

var loadMutexAccess = &sync.Mutex{}
var loadMutexes = map[string]*sync.Mutex{}

func loadMutex(key string) *sync.Mutex {
	loadMutexAccess.Lock()
	defer loadMutexAccess.Unlock()
	if loadMutexes[key] == nil {
		loadMutexes[key] = &sync.Mutex{}
	}
	return loadMutexes[key]
}

func LoadThing(kind string, id int) {
	mutex := loadMutex(fmt.Sprintf("%s:%d", kind, id))
	mutex.Lock()
	defer mutex.Unlock()
}

var sess *r.Session

func Connect() error {
	var err error
	sess, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "geekdorankchart",
	})
	return err
}
