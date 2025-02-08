package services

import (
	"context"
	"errors"

	"GolangWorld/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines user-related operations (Interface Segregation Principle - ISP)
// ISP: Instead of a bloated interface, only necessary methods are defined.
type UserService interface {
	SignUp(ctx context.Context, user *models.User) error
}

// MongoUserService implements UserService (Dependency Inversion Principle - DIP)
// DIP: The service depends on an abstraction (UserService), not a concrete database implementation.
type MongoUserService struct {
	collection *mongo.Collection
}

// NewMongoUserService creates a new MongoUserService
// DIP: Dependency (MongoDB collection) is injected, making it easy to replace with another implementation.
func NewMongoUserService(client *mongo.Client) *MongoUserService {
	return &MongoUserService{
		collection: client.Database("appdb").Collection("users"),
	}
}

// SignUp registers a new user
// OCP: Can be extended with additional validation without modifying core logic.
func (s *MongoUserService) SignUp(ctx context.Context, user *models.User) error {
	user.ID = uuid.New().String()
	if user.Email == "" {
		return errors.New("invalid email")
	}

	// Check if user already exists
	count, err := s.collection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exists")
	}

	// Hash password before storing (SRP: Password hashing should be handled separately)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	_, err = s.collection.InsertOne(ctx, user)
	return err
}
