package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Environments struct {
	ServerAddr       string
	DatabaseAddr     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	SecretKey        []byte
	NatsAddr         string
}

var Env *Environments

func New() (*Environments, error) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}

	serverAddr, exists := os.LookupEnv("SERVER_ADDR")
	if exists == false {
		return nil, fmt.Errorf("SERVER_ADDR not found.")
	}

	databaseAddr, exists := os.LookupEnv("DATABASE_ADDR")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_ADDR not found.")
	}

	databaseUser, exists := os.LookupEnv("DATABASE_USER")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_USER not found.")
	}

	databasePassword, exists := os.LookupEnv("DATABASE_PASSWORD")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_PASSWORD not found.")
	}

	databaseName, exists := os.LookupEnv("DATABASE_NAME")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_NAME not found.")
	}

	secretKey, exists := os.LookupEnv("SECRET_KEY")
	if exists == false {
		return nil, fmt.Errorf("SECRET_KEY not found")
	}

	natsAddr, exists := os.LookupEnv("NATS_ADDR")
	if exists == false {
		return nil, fmt.Errorf("NATS_ADDR not found")
	}

	env := &Environments{
		ServerAddr:       serverAddr,
		DatabaseAddr:     databaseAddr,
		DatabaseUser:     databaseUser,
		DatabasePassword: databasePassword,
		DatabaseName:     databaseName,
		SecretKey:        []byte(secretKey),
		NatsAddr:         natsAddr,
	}

	Env = env
	return env, nil
}
