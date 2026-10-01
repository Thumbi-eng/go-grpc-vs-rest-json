package model

import (
	"time"

	pb "grpc-auth-service/proto/pb"
)

// User represents the domain user entity
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hashed password, never exposed in JSON or protobuf responses
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserInput represents data required to register or create a user
type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

// UpdateUserInput represents optional fields for updating an existing user
type UpdateUserInput struct {
	Username *string
	Email    *string
	Password *string
}

// ToProto maps an internal User model to the Protobuf User message
func (u *User) ToProto() *pb.User {
	if u == nil {
		return nil
	}
	return &pb.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// ToProtoList converts a slice of domain users into a slice of Protobuf users
func ToProtoList(users []*User) []*pb.User {
	pbUsers := make([]*pb.User, len(users))
	for i, u := range users {
		pbUsers[i] = u.ToProto()
	}
	return pbUsers
}
