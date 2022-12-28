package ideas

type Idea struct {
	Id          int
	Text        string
	Description string
	HasSpeaker  bool
	Creator     string
	Votes       int
}
