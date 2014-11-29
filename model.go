package main

import "time"

type Thing struct {
	Kind  string
	Id    int
	Name  string
	Ranks []Rank
}

type Rank struct {
	Date time.Time
	Rank int
}
