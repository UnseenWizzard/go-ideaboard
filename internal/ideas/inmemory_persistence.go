package ideas

import "fmt"

type InMemoryPersistence struct {
	Ideas map[int]Idea
	Votes map[string]map[int]struct{}
}

func NewInMemoryPersistence() InMemoryPersistence {
	return InMemoryPersistence{
		Ideas: make(map[int]Idea),
		Votes: make(map[string]map[int]struct{}),
	}
}

func (p InMemoryPersistence) GetAll() []Idea {
	var list []Idea
	for _, v := range p.Ideas {
		list = append(list, v)
	}
	return list
}

func (p InMemoryPersistence) StoreIdea(idea Idea) error {
	p.Ideas[idea.Id] = idea
	return nil
}

func (p InMemoryPersistence) StoreVote(userId string, ideaId int) (votes int, err error) {
	usrVotes, exists := p.Votes[userId]
	if !exists {
		usrVotes = make(map[int]struct{})
	}
	if _, exists := usrVotes[ideaId]; exists {
		return 0, fmt.Errorf("user %q already voted for idea %d", userId, ideaId)
	}
	usrVotes[ideaId] = struct{}{}
	p.Votes[userId] = usrVotes

	idea, ok := p.Ideas[ideaId]
	if !ok {
		return 0, fmt.Errorf("idea %d does not exist", ideaId)
	}

	idea.Votes++
	p.Ideas[ideaId] = idea

	return idea.Votes, nil
}
