package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
)

func main() {
	sess, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "geekdorankchart",
	})
	if err != nil {
		log.Fatal("Error connecting to database, ", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/{type}/{ids}", HandlerWithDB(sess, ChartHandler))
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
	vars := mux.Vars(req)
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
	wr.Write([]byte(vars["type"]))
	wr.Write([]byte(fmt.Sprintf("%v", ids)))
}
