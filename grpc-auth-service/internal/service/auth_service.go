package service

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/config"
	"grpc-auth-service/internal/model"
	"grpc-auth-service/internal/repository"
	"grpc-auth-service/internal/utils"
	pb "grpc-auth-service/proto/pb"
)

// AuthService implements the pb.AuthServiceServer gRPC interface
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register registers a new user with hashed password and returns a JWT access token
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	username := strings.TrimSpace(req.GetUsername())
	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	// Validate inputs
	if username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is required")
	}
	if email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if len(password) < 6 {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 6 characters")
	}

	// Hash password with bcrypt
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Create user in the database
	createdUser, err := s.userRepo.Create(ctx, model.CreateUserInput{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "email is already registered")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(
		createdUser.ID,
		createdUser.Email,
		createdUser.Username,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate authentication token: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken: token,
		User:        createdUser.ToProto(),
		Message:     "User registered successfully",
	}, nil
}

// Login authenticates a user by email and password, returning a JWT access token
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist")
		}
		return nil, status.Errorf(codes.Internal, "database lookup failed: %v", err)
	}

	// Verify password hash
	if !utils.CheckPassword(password, user.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Username,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate authentication token: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken: token,
		User:        user.ToProto(),
		Message:     "Login successful",
	}, nil
}
