package ideas

type Persistence interface {
	StoreIdea(idea Idea) error
	GetAll() []Idea
	StoreVote(userId string, ideaId int) (votes int, err error)
}

// TODO: actually persist
func New() Persistence {
	return &inMemoryPersistence{
		ideas: make(map[int]Idea),
		votes: make(map[string]map[int]struct{}),
	}
}
