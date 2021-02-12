package main

import (
	"log"
	"os"
	"text/template"
)

type Service struct {
    Name   string
    RateLimit int64
    CacheTtl int64
}

func main() {
	tpl, err := template.ParseFiles("kong.yaml")
	if err != nil {
		log.Fatalln(err)
	}

    // build Service struct by looking up the environment variables

	err = tpl.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
