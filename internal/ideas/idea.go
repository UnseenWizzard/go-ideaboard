package ideas

type Idea struct {
	Id          int    `bson:"id"`
	Text        string `bson:"text,omitempty"`
	Description string `bson:"description,omitempty"`
	HasSpeaker  bool   `bson:"has_speaker,omitempty"`
	Creator     string `bson:"creator,omitempty"`
	Votes       int    `bson:"votes,omitempty"`
}
