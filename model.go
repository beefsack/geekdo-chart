package main

import "time"

type Thing struct {
	Id    string `gorethink:"id,omitempty"`
	Name  string
	Ranks []Rank
}

type Rank struct {
	Date time.Time
	Rank int
}

func NewThing() Thing {
	return Thing{
		Ranks: []Rank{},
	}
}
