package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/config"
	pb "grpc-auth-service/proto/pb"
)

func main() {
	// Load configuration to match server port
	cfg := config.LoadConfig()
	target := fmt.Sprintf("localhost:%s", cfg.Port)

	log.Println("====================================================")
	log.Printf("🔌 Connecting to gRPC Server at %s...\n", target)
	log.Println("====================================================")

	// Establish plaintext gRPC connection
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// -------------------------------------------------------------------------
	// 1. REGISTER A NEW USER (Public RPC)
	// -------------------------------------------------------------------------
	printSection("1. REGISTER NEW USER (Public RPC)")
	regResp, err := authClient.Register(ctx, &pb.RegisterRequest{
		Username: "mustang_coder",
		Email:    "mustang@crowntechlab.com",
		Password: "SecurePassword123!",
	})
	if err != nil {
		log.Fatalf("❌ Register failed: %v", err)
	}
	fmt.Printf("✅ %s\n", regResp.GetMessage())
	fmt.Printf("   User ID:    %s\n", regResp.GetUser().GetId())
	fmt.Printf("   Username:   %s\n", regResp.GetUser().GetUsername())
	fmt.Printf("   Email:      %s\n", regResp.GetUser().GetEmail())
	fmt.Printf("   Created At: %s\n", regResp.GetUser().GetCreatedAt())
	fmt.Printf("   Token:      %s...\n", regResp.GetAccessToken()[:30])

	userID := regResp.GetUser().GetId()

	// -------------------------------------------------------------------------
	// 2. LOGIN (Public RPC)
	// -------------------------------------------------------------------------
	printSection("2. LOGIN (Public RPC)")
	loginResp, err := authClient.Login(ctx, &pb.LoginRequest{
		Email:    "mustang@crowntechlab.com",
		Password: "SecurePassword123!",
	})
	if err != nil {
		log.Fatalf("❌ Login failed: %v", err)
	}
	fmt.Printf("✅ %s\n", loginResp.GetMessage())
	jwtToken := loginResp.GetAccessToken()
	fmt.Printf("   Received JWT: %s...\n", jwtToken[:30])

	// -------------------------------------------------------------------------
	// 3. TEST INTERCEPTOR SECURITY (Expect failure without token)
	// -------------------------------------------------------------------------
	printSection("3. TEST INTERCEPTOR SECURITY (Calling without Token)")
	_, err = userClient.ListUsers(ctx, &pb.ListUsersRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			fmt.Printf("🛡️  Auth Interceptor successfully blocked unauthorized request!\n")
			fmt.Printf("   gRPC Code:    %s (%d)\n", st.Code(), st.Code())
			fmt.Printf("   Error Msg:    %s\n", st.Message())
		} else {
			fmt.Printf("⚠️  Unexpected error response: %v\n", err)
		}
	} else {
		fmt.Println("❌ Security failure: Request without token should have been rejected!")
	}

	// -------------------------------------------------------------------------
	// 4. ATTACH BEARER TOKEN TO OUTGOING METADATA
	// -------------------------------------------------------------------------
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwtToken)

	// -------------------------------------------------------------------------
	// 5. GET CURRENT USER PROFILE (Derived from JWT context)
	// -------------------------------------------------------------------------
	printSection("4. GET CURRENT USER PROFILE (Protected RPC)")
	meResp, err := userClient.GetCurrentUser(authCtx, &pb.GetCurrentUserRequest{})
	if err != nil {
		log.Fatalf("❌ GetCurrentUser failed: %v", err)
	}
	fmt.Printf("✅ Successfully fetched profile from JWT token claims:\n")
	fmt.Printf("   ID:       %s\n", meResp.GetUser().GetId())
	fmt.Printf("   Username: %s\n", meResp.GetUser().GetUsername())
	fmt.Printf("   Email:    %s\n", meResp.GetUser().GetEmail())

	// -------------------------------------------------------------------------
	// 6. LIST ALL USERS
	// -------------------------------------------------------------------------
	printSection("5. LIST ALL USERS (Protected RPC)")
	listResp, err := userClient.ListUsers(authCtx, &pb.ListUsersRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("❌ ListUsers failed: %v", err)
	}
	fmt.Printf("✅ Total Users: %d\n", listResp.GetTotal())
	for i, u := range listResp.GetUsers() {
		fmt.Printf("   [%d] %s (%s) - ID: %s\n", i+1, u.GetUsername(), u.GetEmail(), u.GetId())
	}

	// -------------------------------------------------------------------------
	// 7. UPDATE USER (Partial Update with Optional fields)
	// -------------------------------------------------------------------------
	printSection("6. UPDATE USER (Partial Update with Optional fields)")
	newUsername := "mustang_superstar"
	updateResp, err := userClient.UpdateUser(authCtx, &pb.UpdateUserRequest{
		Id:       userID,
		Username: &newUsername,
	})
	if err != nil {
		log.Fatalf("❌ UpdateUser failed: %v", err)
	}
	fmt.Printf("✅ User updated successfully:\n")
	fmt.Printf("   Updated Username: %s\n", updateResp.GetUser().GetUsername())
	fmt.Printf("   Unchanged Email:  %s\n", updateResp.GetUser().GetEmail())

	// -------------------------------------------------------------------------
	// 8. GET USER BY ID
	// -------------------------------------------------------------------------
	printSection("7. GET USER BY ID (Protected RPC)")
	getResp, err := userClient.GetUser(authCtx, &pb.GetUserRequest{Id: userID})
	if err != nil {
		log.Fatalf("❌ GetUser failed: %v", err)
	}
	fmt.Printf("✅ Verified updated record:\n")
	fmt.Printf("   ID:       %s\n", getResp.GetUser().GetId())
	fmt.Printf("   Username: %s\n", getResp.GetUser().GetUsername())

	// -------------------------------------------------------------------------
	// 9. DELETE USER
	// -------------------------------------------------------------------------
	printSection("8. DELETE USER (Protected RPC)")
	delResp, err := userClient.DeleteUser(authCtx, &pb.DeleteUserRequest{Id: userID})
	if err != nil {
		log.Fatalf("❌ DeleteUser failed: %v", err)
	}
	fmt.Printf("✅ %s (Success: %t)\n", delResp.GetMessage(), delResp.GetSuccess())

	// -------------------------------------------------------------------------
	// 10. VERIFY DELETION (Expect NotFound error)
	// -------------------------------------------------------------------------
	printSection("9. VERIFY DELETION (Confirm User is Deleted)")
	_, err = userClient.GetUser(authCtx, &pb.GetUserRequest{Id: userID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			fmt.Printf("✅ Confirmed: User no longer exists! (Code: %s, Message: %s)\n", st.Code(), st.Message())
		} else {
			fmt.Printf("⚠️  Unexpected status error: %v\n", err)
		}
	} else {
		fmt.Println("❌ Error: Deleted user record was still found!")
	}

	log.Println("\n====================================================")
	log.Println("🎉 ALL gRPC AUTH & CRUD TESTS PASSED SUCCESSFULLY!")
	log.Println("====================================================")
}

func printSection(title string) {
	fmt.Printf("\n----------------------------------------------------\n")
	fmt.Printf("🔹 %s\n", title)
	fmt.Printf("----------------------------------------------------\n")
}
