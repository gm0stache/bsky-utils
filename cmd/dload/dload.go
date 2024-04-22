package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/joho/godotenv"
)

const (
	envBskyHandle string = "BSKY_HANDLE"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	handle := os.Getenv(envBskyHandle)
	if handle == "" {
		log.Fatalf("env var %q must be set", envBskyHandle)
	}

	ctx := context.Background()
	atID, err := getID(ctx, handle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("handle: %q, PDS URL: %q", handle, atID.PDSEndpoint())

	path := atID.DID.String() + ".car"
	if err := downloadRepo(path, atID); err != nil {
		log.Fatal(err)
	}

}

func getID(ctx context.Context, handle string) (*identity.Identity, error) {
	atID, err := syntax.ParseAtIdentifier(handle)
	if err != nil {
		return nil, err
	}
	idDir := identity.DefaultDirectory()
	ident, err := idDir.Lookup(ctx, *atID)
	if err != nil {
		return nil, err
	}
	return ident, nil
}

func downloadRepo(path string, id *identity.Identity) error {
	return nil
}
