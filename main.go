package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/GeertJohan/go.rice"
	"github.com/beefsack/go-geekdo"
	r "github.com/dancannon/gorethink"
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
	router.Handle("/assets/{rest:.*}", http.StripPrefix("/assets/",
		http.FileServer(rice.MustFindBox("assets").HTTPBox())))
	router.HandleFunc("/search", SearchHandler)
	router.HandleFunc("/{kind}/{ids}", Handler(sess, ChartHandler))
	router.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":3000", router)
}

func Handler(
	sess *r.Session,
	handler func(
		wr http.ResponseWriter,
		req *http.Request,
		sess *r.Session,
	),
) func(wr http.ResponseWriter, req *http.Request) {
	return func(wr http.ResponseWriter, req *http.Request) {
		handler(wr, req, sess)
	}
}

func HomeHandler(wr http.ResponseWriter, req *http.Request) {
	http.Redirect(wr, req, "/boardgame/154203,150376,147020,148228,157354,148949", 302)
}

func ChartHandler(
	wr http.ResponseWriter,
	req *http.Request,
	sess *r.Session,
) {
	var overallErr error
	vars := mux.Vars(req)
	kind := vars["kind"]
	if kind == "" {
		wr.WriteHeader(400)
		wr.Write([]byte("You must specify a kind, such as boardgame"))
	}
	ids := []int{}
	idMap := map[int]bool{} // Store which IDs have been parsed for uniqueness
	for _, idStr := range strings.Split(vars["ids"], ",") {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			wr.WriteHeader(400)
			wr.Write([]byte("Each ID must be a number"))
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
		log.Printf("Error loading thing, %v", overallErr)
		wr.WriteHeader(500)
		wr.Write([]byte("Unable to load"))
		return
	}
	graphs, dataProvider, err := ChartJson(things)
	if err != nil {
		log.Printf("Error getting chart JSON, %v", err)
		wr.WriteHeader(500)
		wr.Write([]byte("Unable to generate chart data, possible because no ranks were available"))
		return
	}
	ExecuteTemplate(wr, "chart.tmpl", struct {
		Graphs, DataProvider interface{}
		Kind                 string
	}{
		Graphs:       template.JS(graphs),
		DataProvider: template.JS(dataProvider),
		Kind:         kind,
	})
}

func SearchHandler(wr http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")
	if query == "" {
		wr.WriteHeader(400)
		wr.Write([]byte("You must provide a query"))
		return
	}
	kinds := req.URL.Query().Get("kinds")
	kindStrs := []string{}
	if kinds != "" {
		kindStrs = strings.Split(kinds, ",")
	}
	things, err := client.Search(query, geekdo.SearchOptions{
		Kinds: kindStrs,
	})
	if err != nil {
		log.Printf("Error searching, %v", err)
		wr.WriteHeader(500)
		wr.Write([]byte("Unable to search via the Geekdo API right now"))
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	result := []map[string]interface{}{}
	for _, t := range things.Items {
		result = append(result, map[string]interface{}{
			"type": t.Type,
			"id":   t.Id,
			"name": t.Names[0].Value,
		})
	}
	enc := json.NewEncoder(wr)
	if err := enc.Encode(result); err != nil {
		log.Printf("Error encoding JSON, %v", err)
	}
}
