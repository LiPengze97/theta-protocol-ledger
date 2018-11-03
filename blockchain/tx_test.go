package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thetatoken/ukulele/common"
	"github.com/thetatoken/ukulele/core"
	"github.com/thetatoken/ukulele/crypto"
)

func TestTxIndex(t *testing.T) {
	assert := assert.New(t)

	tx1 := common.Bytes("tx1")
	tx2 := common.Bytes("tx2")
	tx3 := common.Bytes("tx3")
	tx4 := common.Bytes("tx4")
	block1 := core.NewBlock()
	block1.ChainID = "testchain"
	hash := common.BytesToHash(common.Bytes("block1"))
	block1.Hash = hash[:]
	block1.Height = 10
	block1.Txs = []common.Bytes{tx1, tx2, tx3}

	chain := CreateTestChain()
	chain.AddBlock(block1)

	for _, t := range block1.Txs {
		tx, block, found := chain.FindTxByHash(crypto.Keccak256Hash(t))
		assert.True(found)
		assert.NotNil(tx)
		assert.Equal(t, tx)
		assert.NotNil(block)
		assert.Equal(block.Hash, block1.Hash)
	}

	tx, block, found := chain.FindTxByHash(crypto.Keccak256Hash(tx4))
	assert.False(found)
	assert.Nil(tx)
	assert.Nil(block)
}

func TestTxIndexDuplicateTx(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tx1 := common.Bytes("tx1")
	tx2 := common.Bytes("tx2")
	tx3 := common.Bytes("tx3")

	block1 := core.NewBlock()
	block1.ChainID = "testchain"
	hash1 := common.BytesToHash(common.Bytes("block1"))
	block1.Hash = hash1[:]
	block1.Height = 10
	block1.Txs = []common.Bytes{tx1, tx2}

	block2 := core.NewBlock()
	block2.ChainID = "testchain"
	hash2 := common.BytesToHash(common.Bytes("block2"))
	block2.Hash = hash2[:]
	block2.Height = 20
	block2.Txs = []common.Bytes{tx2, tx3}

	chain := CreateTestChain()
	_, err := chain.AddBlock(block1)
	require.Nil(err)

	_, err = chain.AddBlock(block2)
	require.Nil(err)

	tx, block, found := chain.FindTxByHash(crypto.Keccak256Hash(tx1))
	assert.True(found)
	assert.NotNil(tx)
	assert.Equal(tx1, tx)
	assert.NotNil(block)
	assert.Equal(block.Hash, block1.Hash)

	// Tx2 should be linked with block1 instead of block2.
	tx, block, found = chain.FindTxByHash(crypto.Keccak256Hash(tx2))
	assert.True(found)
	assert.NotNil(tx)
	assert.Equal(tx2, tx)
	assert.NotNil(block)
	assert.Equal(block.Hash, block1.Hash)

	tx, block, found = chain.FindTxByHash(crypto.Keccak256Hash(tx3))
	assert.True(found)
	assert.NotNil(tx)
	assert.Equal(tx3, tx)
	assert.NotNil(block)
	assert.Equal(block.Hash, block2.Hash)

	// Tx2 should be linked with block2 after force insert.
	eb := &core.ExtendedBlock{Block: block2}
	chain.AddTxsToIndex(eb, true)
	tx, block, found = chain.FindTxByHash(crypto.Keccak256Hash(tx2))
	assert.True(found)
	assert.NotNil(tx)
	assert.Equal(tx2, tx)
	assert.NotNil(block)
	assert.Equal(block.Hash, block2.Hash)
}