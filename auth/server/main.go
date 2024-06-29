package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/iamstep4ik/quick-meet/auth/lib"
	pb "github.com/iamstep4ik/quick-meet/auth/pb"
	"github.com/iamstep4ik/quick-meet/db"
	"github.com/iamstep4ik/quick-meet/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedRegisterUserServer
	pb.UnimplementedLoginUserServer
	db *sqlx.DB
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	username := req.GetUsername()
	email := req.GetEmail()
	password := req.GetPassword()

	if err := lib.ValidateInput(username, email, password); err != nil {
		return nil, err
	}
	hash_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := models.User{
		Username:     username,
		Email:        email,
		HashPassword: string(hash_password),
		InsertedAt:   time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `INSERT INTO users (username, email, password_hash, inserted_at, updated_at) VALUES (:username, :email, :password_hash, :inserted_at, :updated_at) RETURNING id`
	rows, err := s.db.NamedQuery(query, user)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}
	if rows.Next() {
		err := rows.Scan(&user.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user ID: %v", err)
		}
	}
	rows.Close()

	return &pb.RegisterResponse{
		Message: "User registered successfully",
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User
	username := req.GetUsername()
	password := req.GetPassword()
	query := `SELECT id, password_hash FROM users WHERE username = $1`
	err := s.db.Get(&user, query, username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}
	token, err := lib.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize the database: %v", err)
	}
	defer db.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterRegisterUserServer(grpcServer, &server{db: db})
	pb.RegisterLoginUserServer(grpcServer, &server{db: db})
	reflection.Register(grpcServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
