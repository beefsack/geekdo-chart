package main

import (
	"html/template"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
)

func parseTemplates() (*template.Template, error) {
	box, err := rice.FindBox("templates")
	if err != nil {
		return nil, err
	}
	t := template.New("")
	err = box.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		str, err := box.String(path)
		if err != nil {
			return err
		}
		t, err = t.New(path).Parse(str)
		return err
	})
	return t, err
}
