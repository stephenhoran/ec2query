package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Simple wrapper around the AWS shared credentials provider making it transparent
// to the developer if they have these setup. If the environment variables are set
// this is ignored.
func main() {
	// Checking the credentials file for anything configured
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
			var creds credentials.SharedCredentialsProvider
			if os.Getenv("AWS_PROFILE") != "" {
				creds.Profile = os.Getenv("AWS_PROFILE")
			}
			c, err := creds.Retrieve()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("AWS_ACCESS_KEY_ID=%s\n", c.AccessKeyID)
			fmt.Printf("AWS_SECRET_ACCESS_KEY=%s\n", c.SecretAccessKey)
		}
	}
}
