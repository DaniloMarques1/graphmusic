package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Music struct {
	Id     string `json:"id" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Author string `json:"author" bson:"author"`
}

type MusicRepository interface {
	Save(music *Music) error
	FindAll() ([]Music, error)
	FindByName(name string) (*Music, error)
	UpdateByName(name string, music *Music) (*Music, error)
	RemoveByName(name string) (*Music, error)
	RemoveById(id string) (*Music, error)
}

type MusicRepositoryImpl struct {
	client     *mongo.Client
	collection string
	dbName     string
}

func NewMusicRepository(client *mongo.Client, dbName string) *MusicRepositoryImpl {
	return &MusicRepositoryImpl{client: client, collection: "musics", dbName: dbName}
}

func (mr *MusicRepositoryImpl) Save(music *Music) error {
	collection := mr.client.Database(mr.dbName).Collection(mr.collection)
	_, err := collection.InsertOne(context.Background(), music)
	if err != nil {
		log.Printf("Error inserting music %v\n", err)
		return err
	}
	return nil
}

func (mr *MusicRepositoryImpl) FindAll() ([]Music, error) {
	collection := mr.client.Database(mr.dbName).Collection(mr.collection)
	cursor, err := collection.Find(context.Background(), bson.M{}, &options.FindOptions{})
	if err != nil {
		log.Printf("Error finding all musics %v\n", err)
		return nil, err
	}
	var musics []Music
	if err := cursor.All(context.Background(), &musics); err != nil {
		return nil, err
	}
	return musics, nil
}

func (mr *MusicRepositoryImpl) FindByName(name string) (*Music, error) {
	collection := mr.client.Database(mr.dbName).Collection(mr.collection)
	var music Music
	err := collection.FindOne(
		context.Background(),
		bson.D{{"name", name}},
		&options.FindOneOptions{},
	).Decode(&music)
	if err != nil {
		log.Printf("Error finding by name %v\n", err)
		return nil, err
	}
	return &music, nil
}

func (mr *MusicRepositoryImpl) RemoveByName(name string) (*Music, error) {
	collection := mr.client.Database(mr.dbName).Collection(mr.collection)
	var music Music
	err := collection.FindOneAndDelete(
		context.Background(),
		bson.D{{"name", name}},
		&options.FindOneAndDeleteOptions{},
	).Decode(&music)

	if err != nil {
		return nil, err
	}
	return &music, nil
}

func (mr *MusicRepositoryImpl) UpdateByName(name string, music *Music) (*Music, error) {
	collection := mr.client.Database(mr.dbName).Collection(mr.collection)
	var updatedMusic Music
	update := bson.D{{"$set", music}}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(
		context.Background(),
		bson.D{{"name", name}},
		update,
		opt,
	).Decode(&updatedMusic)
	if err != nil {
		return nil, err
	}
	return &updatedMusic, nil
}

func (mr *MusicRepositoryImpl) RemoveById(id string) (*Music, error) {
	collection := mr.client.Database(mr.dbName).Collection(mr.collection)
	var music Music
	err := collection.FindOneAndDelete(
		context.Background(),
		bson.D{{"_id", id}},
		&options.FindOneAndDeleteOptions{},
	).Decode(&music)
	if err != nil {
		return nil, err
	}
	return &music, nil
}
