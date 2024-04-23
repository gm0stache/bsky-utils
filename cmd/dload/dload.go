package main

import (
	"context"
	"fmt"
	"log"
	"os"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/xrpc"
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
	if err := downloadRepo(ctx, path, atID); err != nil {
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

func downloadRepo(ctx context.Context, path string, id *identity.Identity) error {
	client := xrpc.Client{
		Host: id.PDSEndpoint(),
	}
	repoByts, err := comatproto.SyncGetRepo(ctx, &client, id.DID.String(), "")
	if err != nil {
		return err
	}
	return os.WriteFile(path, repoByts, 0666)
}
