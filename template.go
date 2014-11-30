package main

import (
	"html/template"
	"io"

	"github.com/GeertJohan/go.rice"
)

var tmpl = template.New("")

func ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	box, err := rice.FindBox("templates")
	if err != nil {
		return err
	}
	if tmpl.Lookup(name) == nil {
		str, err := box.String(name)
		if err != nil {
			return err
		}
		if tmpl, err = tmpl.New(name).Parse(str); err != nil {
			return err
		}
	}
	return tmpl.ExecuteTemplate(wr, name, data)
}
