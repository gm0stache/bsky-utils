package dload_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/gm0stache/bsky-utils/pkg/dload"
	"github.com/stretchr/testify/require"
)

const (
	atID         string = "did:plc:5molcdxko5rtwkjivyiviss6"
	carExtension string = "car"
)

func TestCarToDir(t *testing.T) {
	// arrange
	atID := &identity.Identity{
		DID: syntax.DID(atID),
	}

	carFilePath := fmt.Sprintf("../../testdata/%s.%s", atID.DID.String(), carExtension)

	// act
	err := dload.ConvertCarToDir(context.TODO(), carFilePath, atID) // todo: fix tests. there seems to be an issue with the lex-to-record-type conversion. init lexutil?

	// assert
	require.NoError(t, err)
}
