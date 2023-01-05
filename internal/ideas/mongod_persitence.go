package ideas

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Vote struct {
	UsrID  string `bson:"usr_id"`
	IdeaID int    `bson:"idea_id"`
}

type MongoDBPersistence struct {
	client *mongo.Client
}

func (m *MongoDBPersistence) Close() {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func NewMongoDBPersistence(connectionstring string) (MongoDBPersistence, error) {
	if connectionstring == "" {
		return MongoDBPersistence{}, fmt.Errorf("connection string needs to be defined")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionstring))
	if err != nil {
		return MongoDBPersistence{}, fmt.Errorf("Failed to create mongodb client: %w", err)
	}

	return MongoDBPersistence{
		client: client,
	}, nil
}

func (m *MongoDBPersistence) GetAll() ([]Idea, error) {
	coll := m.client.Database("idea_board").Collection("ideas")
	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return []Idea{}, fmt.Errorf("failed query ideas: %w", err)
	}

	var ideas []Idea
	if err := cur.All(context.TODO(), &ideas); err != nil {
		return []Idea{}, fmt.Errorf("failed to parse existing ideas: %w", err)
	}

	return ideas, nil
}

func (m *MongoDBPersistence) StoreIdea(idea Idea) error {
	coll := m.client.Database("idea_board").Collection("ideas")
	result, err := coll.InsertOne(context.TODO(), idea)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted document with id: %v\n", result.InsertedID)
	return nil
}

func (m *MongoDBPersistence) StoreVote(userId string, ideaId int) (votes int, err error) {
	ideas := m.client.Database("idea_board").Collection("ideas")
	idFilter := bson.D{{"id", ideaId}}
	idea := ideas.FindOne(context.TODO(), idFilter)
	if idea == nil {
		return 0, fmt.Errorf("idea %d does not exist", ideaId)
	}

	vote := Vote{UsrID: userId, IdeaID: ideaId}

	votesColl := m.client.Database("idea_board").Collection("votes")
	existingVote := votesColl.FindOne(context.TODO(), vote)
	if existingVote.Err() != mongo.ErrNoDocuments {
		var v bson.D
		existingVote.Decode(&v)
		return 0, fmt.Errorf("user %q already voted for idea %d: %v", userId, ideaId, v)
	}

	voteRes, err := votesColl.InsertOne(context.TODO(), vote)
	if err != nil {
		return 0, fmt.Errorf("failed to store vote %v: %w", vote, err)
	}
	fmt.Printf("Stored vote %v with document id: %v\n", vote, voteRes.InsertedID)

	incrementVotes := bson.D{{"$inc", bson.D{{"votes", 1}}}}
	_, err = ideas.UpdateOne(context.TODO(), idFilter, incrementVotes)
	if err != nil {
		return 0, fmt.Errorf("failed to count vote %v: %w", vote, err)
	}

	return 0, nil
}
