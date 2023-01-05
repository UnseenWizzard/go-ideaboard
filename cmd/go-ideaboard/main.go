package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/unseenwizzard/go-ideaboard/internal/ideas"
)

var addr = flag.String("addr", ":8080", "http service address")
var basepath = flag.String("basepath", "/", "base path of service")

var templ = template.Must(template.ParseFiles("web/index.html"))

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

var idealist ideas.MongoDBPersistence

type templateArgs struct {
	BasePath   string
	CreatePath string
	VotePath   string
	Ideas      []ideas.Idea
}

func main() {
	flag.Parse()

	uri := os.Getenv("MONGODB_CONNECTION_URI")
	var err error
	idealist, err = ideas.NewMongoDBPersistence(uri)
	if err != nil {
		log.Fatal("MongoDB connection err: ", err)
	}

	http.Handle("/", http.HandlerFunc(display))
	http.Handle("/static/custom.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/custom.css")
	}))

	err = http.ListenAndServe(*addr, nil) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
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

	list, err := idealist.GetAll()
	if err != nil {
		log.Println("failed to get existing ideas: ", err)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Votes > list[j].Votes //sort descending
	})

	args := templateArgs{
		BasePath:   *basepath,
		CreatePath: *basepath,
		VotePath:   *basepath,
		Ideas:      list,
	}

	err = templ.Execute(w, args)
	if err != nil {
		log.Println("failed to execute template: ", err)
	}

}

func addIdea(req *http.Request, uid string) {
	id := randGen.Intn(2560)

	hasSpeaker := false
	if req.FormValue("hasSpeaker") == "on" {
		hasSpeaker = true
	}

	i := ideas.Idea{
		Id:          id,
		Text:        req.FormValue("idea"),
		Description: req.FormValue("details"),
		HasSpeaker:  hasSpeaker,
		Creator:     uid,
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
