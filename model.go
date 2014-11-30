package main

import (
	"fmt"
	"time"
)

type Identifier struct {
	Kind string
	Id   int
}

func (i Identifier) Key() string {
	return fmt.Sprintf("%s:%d", i.Kind, i.Id)
}

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
