package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

var templ = template.Must(template.ParseFiles("resources/index.html"))

type idea struct {
	Id int
	Text string
	Creator string
	Votes int
}

var inputs = make(map[int]idea)
var votes = make(map[string]map[int]struct{})

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
    flag.Parse()
    http.Handle("/", http.HandlerFunc(display))
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func display(w http.ResponseWriter, req *http.Request) {

	c, err := req.Cookie("uid")
	if err != nil {
		c = &http.Cookie{
			Name: "uid",
			Value: fmt.Sprintf("u_%d",randGen.Intn(64000)),
		}
	}

	switch req.Method {
	case http.MethodPost:
		switch req.FormValue("type") {
		case "input":
			addIdea(req)
		case "vote":
			countVote(req)
			//TODO display double votes and errors to user
		}
	}
	http.SetCookie(w, c)
	templ.Execute(w, inputs)
    
}

func addIdea(req *http.Request) {
	id := randGen.Intn(2560)
	i := idea {
		Id: id,
		Text: req.FormValue("idea"),
		Creator: req.UserAgent(),
	}
	inputs[id] = i
}

func countVote(req *http.Request) {
	idStr := req.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("got id that was not an int: ", idStr)
		return
	}

	c, err := req.Cookie("uid")
	if err != nil {
		println("failed to get uid cookie: ", err)
	}

	uid := c.Value

	usrVotes, exists := votes[uid]
	if !exists {
		usrVotes = make(map[int]struct{})
	}
	if _, exists := usrVotes[id]; exists {
		fmt.Println("user", uid, "already voted for item", id)
		return
	}
	usrVotes[id] = struct{}{}
	votes[uid] = usrVotes

	if idea, ok := inputs[id]; ok {
		idea.Votes++
		inputs[id] = idea
	}
}
