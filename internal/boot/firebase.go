package boot

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"google.golang.org/api/option"
)

func getFirebaseCredentialsForProd() string {
	secretName := os.Getenv("FIREBASE_SERVICE_ACCOUNT_SECRET")
	region := os.Getenv("AWS_REGION")

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	return *result.SecretString
}

func NewFirebaseAuth() *auth.Client {
	// Initialize Firebase with service account
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	var opt option.ClientOption
	if os.Getenv("STAGE") == "prod" {
		opt = option.WithCredentialsJSON([]byte(getFirebaseCredentialsForProd()))
	} else {
		opt = option.WithCredentialsFile(serviceAccountPath)
	}
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	authClient, err := app.Auth(context.Background())

	if err != nil {
		log.Fatalf("error initializing auth client: %v", err)
	}

	return authClient
}
