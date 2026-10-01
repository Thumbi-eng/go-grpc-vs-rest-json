package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc-auth-service/internal/config"
	"grpc-auth-service/internal/db"
	"grpc-auth-service/internal/interceptor"
	"grpc-auth-service/internal/repository"
	"grpc-auth-service/internal/service"
	pb "grpc-auth-service/proto/pb"
)

func main() {
	// 1. Load application configuration
	cfg := config.LoadConfig()

	// 2. Initialize in-memory database and repository
	database := db.NewMemoryDB()
	userRepo := repository.NewUserRepository(database)

	// 3. Initialize gRPC services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)

	// 4. Initialize Auth Interceptor
	authInterceptor := interceptor.NewAuthInterceptor(cfg.JWTSecret)

	// 5. Create gRPC server with the unary auth interceptor attached
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	// 6. Register gRPC service implementations
	pb.RegisterAuthServiceServer(grpcServer, authService)
	pb.RegisterUserServiceServer(grpcServer, userService)

	// 7. Enable gRPC Server Reflection
	// Allows tools like Postman, BloomRPC, and grpcurl to introspect available services without manual proto files
	reflection.Register(grpcServer)

	// 8. Bind TCP listener on the configured port
	addr := fmt.Sprintf(":%s", cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("❌ Failed to bind TCP listener on %s: %v", addr, err)
	}

	// 9. Start gRPC server in a separate goroutine so the main routine can listen for shutdown signals
	go func() {
		log.Println("====================================================")
		log.Printf("🚀 gRPC Server listening on port %s\n", cfg.Port)
		log.Println("📋 Registered Services:")
		log.Println("   • pb.AuthService (Register, Login)")
		log.Println("   • pb.UserService (GetUser, GetCurrentUser, ListUsers, UpdateUser, DeleteUser)")
		log.Println("🛡️  Auth Interceptor: Active (JWT Bearer Token required for UserService)")
		log.Println("🔍 gRPC Reflection:  Enabled (Ready for Postman / grpcurl)")
		log.Println("====================================================")

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("❌ gRPC server encountered an error: %v", err)
		}
	}()

	// 10. Set up graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Block until a shutdown signal is received (e.g. Ctrl+C)
	sig := <-quit
	log.Printf("\n🛑 Received signal (%s). Gracefully stopping gRPC server...\n", sig)

	grpcServer.GracefulStop()
	log.Println("✅ gRPC server stopped cleanly. Goodbye!")
}
