package eth2http

import (
	"testing"

	"github.com/stretchr/testify/assert"

	eth2client "github.com/skillz-blockchain/go-utils/ethereum/consensus/client"
)

func TestClientImplementsInterface(t *testing.T) {
	iClient := (*eth2client.Client)(nil)
	client := new(Client)
	assert.Implements(t, iClient, client)
}
