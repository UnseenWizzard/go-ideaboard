package ideas

import "fmt"

type inMemoryPersistence struct {
	ideas map[int]Idea
	votes map[string]map[int]struct{}
}

func (p inMemoryPersistence) GetAll() []Idea {
	var list []Idea
	for _, v := range p.ideas {
		list = append(list, v)
	}
	return list
}

func (p inMemoryPersistence) StoreIdea(idea Idea) error {
	p.ideas[idea.Id] = idea
	return nil
}

func (p inMemoryPersistence) StoreVote(userId string, ideaId int) (votes int, err error) {
	usrVotes, exists := p.votes[userId]
	if !exists {
		usrVotes = make(map[int]struct{})
	}
	if _, exists := usrVotes[ideaId]; exists {
		return 0, fmt.Errorf("user %q already voted for idea %d", userId, ideaId)
	}
	usrVotes[ideaId] = struct{}{}
	p.votes[userId] = usrVotes

	idea, ok := p.ideas[ideaId]
	if !ok {
		return 0, fmt.Errorf("idea %d does not exist", ideaId)
	}

	idea.Votes++
	p.ideas[ideaId] = idea

	return idea.Votes, nil
}
