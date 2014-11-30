package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

func ChartJson(things []Thing) (graphs []byte, dataProvider []byte, err error) {
	var earliest, latest time.Time
	g := []map[string]interface{}{}
	d := []map[string]interface{}{}
	for _, t := range things {
		var (
			kind  string
			id, n int
		)
		n, err = fmt.Sscanf(
			strings.Replace(t.Id, ":", " ", -1),
			"%s %d",
			&kind,
			&id,
		)
		if err != nil || n == 0 {
			err = errors.New("could not get ID")
			return
		}
		g = append(g, map[string]interface{}{
			"title":      t.Name,
			"valueField": t.Id,
			"id":         id,
		})
		for _, r := range t.Ranks {
			if r.Rank == 0 {
				continue
			}
			if earliest.IsZero() || r.Date.Before(earliest) {
				earliest = r.Date
			}
			if latest.IsZero() || r.Date.After(latest) {
				latest = r.Date
			}
		}
	}
	if earliest.IsZero() {
		err = errors.New("no data found")
		return
	}
	thingPtrs := map[int]int{}
	for !earliest.After(latest) {
		val := map[string]interface{}{
			"date": earliest.Format("2006-01-02"),
		}
		for i, t := range things {
			// If this thing is earlier than the earliest, increment up to the value not before.
			for thingPtrs[i] < len(t.Ranks) &&
				t.Ranks[thingPtrs[i]].Date.Before(earliest) {
				thingPtrs[i] += 1
			}
			if thingPtrs[i] >= len(t.Ranks) {
				continue
			}
			rank := t.Ranks[thingPtrs[i]]
			if rank.Date.Equal(earliest) {
				if rank.Rank > 0 {
					val[t.Id] = rank.Rank
				}
				thingPtrs[i] += 1
			}
		}
		d = append(d, val)
		earliest = earliest.AddDate(0, 0, 1)
	}
	if graphs, err = json.Marshal(g); err != nil {
		return
	}
	dataProvider, err = json.Marshal(d)
	return
}
