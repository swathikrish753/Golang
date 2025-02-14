package services

import (
	"GolangWorld/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ListService interface {
	ListAllUsers(ctx context.Context) ([]models.User, error)
}

type ListUserService struct {
	collection *mongo.Collection
}

func NewListUserService(client *mongo.Client) *ListUserService {
	return &ListUserService{
		collection: client.Database("appdb").Collection("users"),
	}
}

func (s *ListUserService) ListAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	cursor, err := s.collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
