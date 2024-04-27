package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gm0stache/bsky-utils/pkg/dload"
	"github.com/joho/godotenv"
)

const (
	envBskyHandle  string = "BSKY_HANDLE"
	envCarFilePath string = "CAR_FILE"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	handle := os.Getenv(envBskyHandle)
	if handle == "" {
		log.Fatalf("env var %q must be set", envBskyHandle)
	}

	carFilePath := os.Getenv(envCarFilePath)
	if handle == "" {
		log.Fatalf("env var %q must be set", envBskyHandle)
	}

	ctx := context.Background()
	atID, err := dload.GetATID(ctx, handle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("handle: %q, PDS URL: %q\n", handle, atID.PDSEndpoint())

	if err := dload.ConvertCarToDir(ctx, carFilePath, atID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("done.")
}
