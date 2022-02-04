package eth2http

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skillz-blockchain/go-utils/eth2/client"
)

func TestClientImplementsInterface(t *testing.T) {
	iClient := (*client.Client)(nil)
	client := new(Client)
	assert.Implements(t, iClient, client)
}
