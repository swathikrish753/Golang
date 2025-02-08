package models

// User defines the structure of a user document in MongoDB
// SRP: This struct only represents user data.
type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password,omitempty" bson:"password"`
}
