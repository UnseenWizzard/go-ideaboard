package main

import (
	"flag"
	"fmt"
	"github.com/unseenwizzard/go-ideaboard/internal/ideas"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

var templ = template.Must(template.ParseFiles("web/index.html"))

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

var idealist = ideas.New()

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(display))
	http.Handle("/static/custom.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/custom.css")
	}))
	// TODO: TLS!
	err := http.ListenAndServe(*addr, nil) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func display(w http.ResponseWriter, req *http.Request) {

	var c *http.Cookie
	if existingCookie, err := req.Cookie("uid"); err == nil {
		c = existingCookie
	} else {
		c = &http.Cookie{
			Name:     "uid",
			Value:    fmt.Sprintf("u_%d", randGen.Intn(64000)),
			HttpOnly: true,
			Secure:   true,
		}
	}

	switch req.Method {
	case http.MethodPost:
		switch req.FormValue("type") {
		case "input":
			addIdea(req, c.Value)
		case "vote":
			countVote(req)
			// TODO display double votes and errors to user
		}
		// TODO allow deletion of my own ideas
	}
	http.SetCookie(w, c)

	list := idealist.GetAll()
	sort.Slice(list, func(i, j int) bool {
		return list[i].Votes > list[j].Votes //sort descending
	})
	err := templ.Execute(w, list)
	if err != nil {
		log.Println("failed to execute template: ", err)
	}

}

func addIdea(req *http.Request, uid string) {
	// TODO: allow more user input (description, present/idea, who) - in struct and html template
	id := randGen.Intn(2560)
	i := ideas.Idea{
		Id:      id,
		Text:    req.FormValue("idea"),
		Creator: uid,
	}
	err := idealist.StoreIdea(i)
	if err != nil {
		log.Println("failed to store idea:", err)
	}
}

func countVote(req *http.Request) {
	idStr := req.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("got id that was not an int: ", idStr)
		return
	}

	c, err := req.Cookie("uid")
	if err != nil {
		log.Println("failed to get uid cookie: ", err)
	}

	uid := c.Value

	votes, err := idealist.StoreVote(uid, id)
	if err != nil {
		log.Println("failed to store vote:", err)
		return
	}

	log.Printf("counted %qs vote for idea %d (%d total votes)", uid, id, votes)

}
