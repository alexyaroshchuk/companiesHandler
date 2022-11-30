package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"companiesHandler/auth"
	"companiesHandler/middlewares"
	proto "companiesHandler/proto/pb"
	"companiesHandler/services"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const tokenDuration = 120 * time.Minute

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(fmt.Println("Connection couldn't be opened", err))
	}
	defer db.Close()

	fmt.Printf("Starting producer...")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	if err != nil {
		log.Fatal(fmt.Println("Can't start kafka producer", err))
	}
	defer p.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(fmt.Println("Connection not established, ping didn't work", err))
	}
	userService := services.NewUserService(db)
	companyService := services.NewCompanyService(db, p)
	jwtManager := auth.NewJWTManager(os.Getenv("SECRET_KEY"), tokenDuration)
	authServer := services.NewAuthServer(userService, jwtManager)
	interceptor := auth.NewAuthInterceptor(jwtManager, middlewares.AccessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
	}
	s := grpc.NewServer(serverOptions...)
	proto.RegisterAuthServiceServer(s, authServer)
	proto.RegisterUserServiceServer(s, userService)
	proto.RegisterCompanyServiceServer(s, companyService)

	tl, err := net.Listen("tcp", os.Getenv("TCP_ADD"))
	if err != nil {
		log.Fatal(fmt.Println("Error starting tcp listener", err))
	}

	if err := s.Serve(tl); err != nil {
		log.Fatal(fmt.Println("Failed to server the listener", err))
	}
}
