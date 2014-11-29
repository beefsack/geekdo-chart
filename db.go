package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/beefsack/go-geekdo"
	r "github.com/dancannon/gorethink"
)

const (
	DbName         = "geekdorc"
	TableName      = "things"
	RatingsPerPage = 100
)

var loadMutexAccess = &sync.Mutex{}
var loadMutexes = map[string]*sync.Mutex{}
var client = geekdo.NewClient()

func migrate(sess *r.Session) error {
	var name string
	// Create database if it doesn't exist.
	res, err := r.DbList().Run(sess)
	if err != nil {
		return err
	}
	found := false
	for res.Next(&name) {
		if name == DbName {
			found = true
			break
		}
	}
	if err := res.Close(); err != nil {
		return err
	}
	if !found {
		log.Printf("Creating database %s", DbName)
		if _, err := r.DbCreate(DbName).RunWrite(sess); err != nil {
			return err
		}
	}
	// Create table if it doesn't exist.
	res, err = r.Db(DbName).TableList().Run(sess)
	if err != nil {
		return err
	}
	found = false
	for res.Next(&name) {
		if name == TableName {
			found = true
			break
		}
	}
	if err := res.Close(); err != nil {
		return err
	}
	if !found {
		log.Printf("Creating table %s", TableName)
		if _, err := r.Db(DbName).TableCreate(TableName).RunWrite(sess); err != nil {
			return err
		}
	}
	return nil
}

func loadMutex(key string) *sync.Mutex {
	loadMutexAccess.Lock()
	defer loadMutexAccess.Unlock()
	if loadMutexes[key] == nil {
		loadMutexes[key] = &sync.Mutex{}
	}
	return loadMutexes[key]
}

func Key(kind string, id int) string {
	return fmt.Sprintf("%s:%d", kind, id)
}

func LoadThing(kind string, id int, sess *r.Session) (Thing, error) {
	key := Key(kind, id)
	mutex := loadMutex(key)
	mutex.Lock()
	defer mutex.Unlock()
	// See if we have it in the DB at all, if not, fetch it.
	thing := NewThing()
	thing.Id = key
	res, err := r.Db(DbName).Table(TableName).Get(key).Run(sess)
	if err != nil {
		return thing, err
	}
	if !res.IsNil() {
		if err := res.One(&thing); err != nil {
			return thing, err
		}
	}
	if l := len(thing.Ranks); l == 0 || thing.Ranks[l-1].Date.Before(
		time.Now().AddDate(0, 0, -2)) {
		// We don't have ranks or the ranks are old, fetch more.
		page := 1 + (l / RatingsPerPage)
		for {
			log.Printf("Loading page %d for %s %d", page, kind, id)
			things, err := client.Thing(kind, id, geekdo.ThingOptions{
				Historical: true,
				Page:       page,
			})
			if err != nil {
				return thing, err
			}
			if things.Items == nil || len(things.Items) == 0 {
				return thing, errors.New("could not find items by that kind and id")
			}
			item := things.Items[0]
			if thing.Name == "" && item.Names != nil && len(item.Names) > 0 {
				thing.Name = things.Items[0].Names[0].Value
			}
			ratings := item.Statistics.Ratings
			if ratings == nil {
				return thing, errors.New("no ratings fetched")
			}
			for _, rating := range ratings {
				if rating.Ranks == nil || len(rating.Ranks) == 0 {
					continue
				}
				d, err := rating.Time()
				if err != nil {
					continue
				}
				if l := len(thing.Ranks); l > 0 && !d.After(thing.Ranks[l-1].Date) {
					continue
				}
				thing.Ranks = append(thing.Ranks, Rank{
					Date: d,
					Rank: rating.Ranks[0].Value,
				})
			}
			if len(ratings) < RatingsPerPage {
				// We've hit the end, stop.
				break
			}
			page += 1
		}
		// Up to date, save.
		if _, err := r.Db(DbName).Table(TableName).Insert(thing).RunWrite(sess); err != nil {
			return thing, err
		}
	}
	return thing, nil
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
