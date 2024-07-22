package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccessLevel int

const (
	AccessViewer AccessLevel = iota // iota starts at 0
	AccessModerator
	AccessAdmin
)

type User struct {
	ID       string      `bson:"_id,omitempty"`
	Username string      `bson:"username"`
	Email    string      `bson:"email"`
	Password string      `bson:"password"`
	Access   AccessLevel `bson:"access"`
}

type UserStorage struct {
	collection *mongo.Collection
}

func NewUserStorage(db *mongo.Database) *UserStorage {
	return &UserStorage{
		collection: db.Collection("users"),
	}
}

func (s *UserStorage) CreateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.collection.InsertOne(ctx, user)
	return err
}

func (s *UserStorage) GetUser(id string) (*User, error) {
	var user User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) UpdateUser(id string, user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.collection.UpdateByID(ctx, id, bson.M{"$set": user})
	return err
}

func (s *UserStorage) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *UserStorage) GetAllUsers() ([]User, error) {
	var users []User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
