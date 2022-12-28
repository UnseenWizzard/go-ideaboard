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

var templ = template.Must(template.New("list").Parse(templateStr))

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

const templateStr = `
<html>
<head>
<title>Go Idea Board</title>
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-GLhlTQ8iRABdZLl6O3oVMWSktQOp6b7In1Zl3/Jr59b6EGGoI1aFkw7cmDA6j6gD" crossorigin="anonymous">
</head>
<body>
<section class="vh-100" style="background-color: #e2d5de;">
  <div class="container py-5 h-100">
    <div class="row d-flex justify-content-center align-items-center h-100">
      <div class="col col-xl-10">

        <div class="card" style="border-radius: 15px;">
          <div class="card-body p-5">

            <h6 class="mb-3">Go Idea Board</h6>

            <form action="/" method="POST" class="d-flex justify-content-center align-items-center mb-4">
              <div class="form-outline flex-fill">
                <input type="text" id="idea" name="idea" value="" title="Idea to add" class="form-control form-control-lg" />
              </div>
			  <input type="hidden" name="type" value="input" />
              <button type="submit" class="btn btn-primary btn-lg ms-2">Add</button>
            </form>

            <ul class="list-group mb-0">
			{{ range . }}
              <li
                class="list-group-item d-flex justify-content-between align-items-start border-start-0 border-top-0 border-end-0 border-bottom rounded-0 mb-2">
				<div class="ms-2 me-auto">
					<div class="fw-bold">{{.Text}}</div>
					 {{.Creator}} 
			  	</div>
			  	<span class="badge bg-primary rounded-pill">{{.Votes}}</span>
				<form action="/" method="POST">
					<input type="hidden" name="id" value="{{.Id}}" />
					<input type="hidden" name="type" value="vote" />
				  	<button type="sumbit" class="btn btn-primary btn-sm ms-2">Vote</button>
				</form>
              </li>
			{{ end }}
            </ul>

          </div>
        </div>

      </div>
    </div>
  </div>
</section>
</body>
</html>
`