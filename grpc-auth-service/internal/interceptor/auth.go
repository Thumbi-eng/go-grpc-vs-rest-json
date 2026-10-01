package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/utils"
)

type contextKey string

const userClaimsContextKey contextKey = "user_claims"

// UserClaimsContext represents the authenticated user data stored in the request context
type UserClaimsContext struct {
	UserID   string
	Email    string
	Username string
}

// AuthInterceptor enforces JWT authentication on protected gRPC endpoints
type AuthInterceptor struct {
	jwtSecret     string
	publicMethods map[string]bool
}

// NewAuthInterceptor creates a new AuthInterceptor with whitelist for public methods
func NewAuthInterceptor(jwtSecret string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSecret: jwtSecret,
		publicMethods: map[string]bool{
			"/pb.AuthService/Register": true,
			"/pb.AuthService/Login":    true,
		},
	}
}

// Unary returns a gRPC UnaryServerInterceptor function
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// If the method is public, allow access without token validation
		if i.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract metadata from the incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata from request")
		}

		// Retrieve the authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is required")
		}

		// Expect format: 'Bearer <token>'
		authHeader := authHeaders[0]
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format; expected 'Bearer <token>'")
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			return nil, status.Errorf(codes.Unauthenticated, "token cannot be empty")
		}

		// Validate the token and parse claims
		claims, err := utils.ValidateToken(tokenString, i.jwtSecret)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token: %v", err)
		}

		// Inject the authenticated user identity into context
		userContext := &UserClaimsContext{
			UserID:   claims.UserID,
			Email:    claims.Email,
			Username: claims.Username,
		}

		newCtx := context.WithValue(ctx, userClaimsContextKey, userContext)

		// Proceed to handler with enriched context
		return handler(newCtx, req)
	}
}

// GetUserClaims extracts the authenticated user's claims from the gRPC context
func GetUserClaims(ctx context.Context) (*UserClaimsContext, error) {
	claims, ok := ctx.Value(userClaimsContextKey).(*UserClaimsContext)
	if !ok || claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: user claims not found in context")
	}
	return claims, nil
}

// GetUserID is a convenience helper to extract only the UserID from the context
func GetUserID(ctx context.Context) (string, error) {
	claims, err := GetUserClaims(ctx)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
