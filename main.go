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

var sess *r.Session

func main() {
	var err error
	sess, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "geekdorankchart",
	})
	if err != nil {
		log.Fatal("Error connecting to database, ", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/{type}/{ids}", ChartHandler)
	router.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":3000", router)
}

func HomeHandler(wr http.ResponseWriter, req *http.Request) {
	wr.Write([]byte("hello"))
}

func ChartHandler(wr http.ResponseWriter, req *http.Request) {
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
