package service

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/interceptor"
	"grpc-auth-service/internal/model"
	"grpc-auth-service/internal/repository"
	"grpc-auth-service/internal/utils"
	pb "grpc-auth-service/proto/pb"
)

// UserService implements the pb.UserServiceServer gRPC interface
type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUser fetches a single user by their ID
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.UserResponse{
		User: user.ToProto(),
	}, nil
}

// GetCurrentUser returns the profile of the caller extracted from the validated JWT in gRPC metadata
func (s *UserService) GetCurrentUser(ctx context.Context, _ *pb.GetCurrentUserRequest) (*pb.UserResponse, error) {
	// Extract the authenticated user's ID injected by the AuthInterceptor
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "authenticated user account no longer exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve profile: %v", err)
	}

	return &pb.UserResponse{
		User: user.ToProto(),
	}, nil
}

// ListUsers returns a paginated list of all registered users
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, total, err := s.userRepo.List(ctx, req.GetPage(), req.GetLimit())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	return &pb.ListUsersResponse{
		Users: model.ToProtoList(users),
		Total: total,
	}, nil
}

// UpdateUser applies partial updates to a user record
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	updateInput := model.UpdateUserInput{}

	// Check optional fields using protobuf presence getters
	if req.Username != nil {
		username := strings.TrimSpace(req.GetUsername())
		updateInput.Username = &username
	}

	if req.Email != nil {
		email := strings.TrimSpace(req.GetEmail())
		updateInput.Email = &email
	}

	if req.Password != nil {
		password := req.GetPassword()
		if len(password) < 6 {
			return nil, status.Errorf(codes.InvalidArgument, "new password must be at least 6 characters")
		}
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash new password: %v", err)
		}
		updateInput.Password = &hashedPassword
	}

	updatedUser, err := s.userRepo.Update(ctx, id, updateInput)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "email already in use by another account")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		User: updatedUser.ToProto(),
	}, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}
