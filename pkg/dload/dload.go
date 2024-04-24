package dload

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/ipfs/go-cid"
)

// GetATID resolves a given handle to an ATproto identity.
func GetATID(ctx context.Context, handle string) (*identity.Identity, error) {
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

// DownloadRepo downloads the PDS content for a given user.
func DownloadRepo(ctx context.Context, path string, id *identity.Identity) error {
	client := xrpc.Client{
		Host: id.PDSEndpoint(),
	}
	repoByts, err := atproto.SyncGetRepo(ctx, &client, id.DID.String(), "")
	if err != nil {
		return err
	}
	return os.WriteFile(path, repoByts, 0666)
}

// ConvertCarToDir extracts all resources of a '.car' export into a directory.
// The created files are human-readable (JSON).
func ConvertCarToDir(ctx context.Context, carPath string, atID *identity.Identity) error {
	fi, err := os.Open(carPath)
	if err != nil {
		return err
	}

	r, err := repo.ReadRepoFromCar(ctx, fi)
	if err != nil {
		return err
	}

	sc := r.SignedCommit()
	did, err := syntax.ParseDID(sc.Did)
	if err != nil {
		return err
	}
	fmt.Printf("DID's are identical (handle/repo): %t\n", did.String() == atID.DID.String())

	rootDir := "."
	err = r.ForEach(ctx, "", func(k string, v cid.Cid) error {
		_, rec, err := r.GetRecord(ctx, k)
		if err != nil {
			return err
		}

		recJson, err := json.MarshalIndent(rec, "", "  ")
		if err != nil {
			return err
		}

		recPath := rootDir + "/" + k
		if err := os.MkdirAll(filepath.Dir(recPath), os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(recPath+".json", recJson, 0666); err != nil {
			return err
		}

		return nil
	})

	return err
}
