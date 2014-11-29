package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	r "github.com/dancannon/gorethink"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

func main() {
	sess, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "geekdorankchart",
	})
	if err != nil {
		log.Fatalf("Error connecting to database, %v", err)
	}
	if err := migrate(sess); err != nil {
		log.Fatalf("Error migrating database, %v", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/{kind}/{ids}", HandlerWithDB(sess, ChartHandler))
	router.HandleFunc("/", HandlerWithDB(sess, HomeHandler))
	http.ListenAndServe(":3000", router)
}

func HandlerWithDB(
	sess *r.Session,
	handler func(wr http.ResponseWriter, req *http.Request, sess *r.Session),
) func(wr http.ResponseWriter, req *http.Request) {
	return func(wr http.ResponseWriter, req *http.Request) {
		handler(wr, req, sess)
	}
}

func HomeHandler(wr http.ResponseWriter, req *http.Request, sess *r.Session) {
	wr.Write([]byte("hello"))
}

func ChartHandler(wr http.ResponseWriter, req *http.Request, sess *r.Session) {
	var overallErr error
	vars := mux.Vars(req)
	kind := vars["kind"]
	ids := []int{}
	idMap := map[int]bool{} // Store which IDs have been parsed for uniqueness
	for _, idStr := range strings.Split(vars["ids"], ",") {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			wr.WriteHeader(400)
			wr.Write([]byte("All ids must be integers"))
			return
		}
		if idMap[id] {
			continue
		}
		idMap[id] = true
		ids = append(ids, id)
	}
	// Load from DB
	things := []Thing{}
	wg := sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			thing, err := LoadThing(kind, id, sess)
			if err != nil {
				overallErr = err
			}
			things = append(things, thing)
			wg.Done()
		}(id)
	}
	wg.Wait()
	if overallErr != nil {
		log.Fatalf("Error loading thing, %v", overallErr)
	}
	spew.Dump(things)
}
