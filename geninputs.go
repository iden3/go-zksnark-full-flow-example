package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/genesis"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
	"github.com/stretchr/testify/assert"
)

func IdStateInputs(t *testing.T) string {

	nLevels := 3

	privKHex := "28156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
	// Create new claim
	var k babyjub.PrivateKey
	if _, err := hex.Decode(k[:], []byte(privKHex)); err != nil {
		panic(err)
	}
	// fmt.Println("sk", skToBigInt(&k))
	pk := k.Public()

	claimKOp := claims.NewClaimKeyBabyJub(pk, 1)

	clt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), nLevels)
	assert.Nil(t, err)
	rot, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), nLevels)
	assert.Nil(t, err)

	id, err := genesis.CalculateIdGenesisMT(clt, rot, claimKOp, []merkletree.Entrier{})
	assert.Nil(t, err)

	// get claimproof
	hi, err := claimKOp.Entry().HIndex()
	assert.Nil(t, err)
	// generate merkle proof
	proof, err := clt.GenerateProof(hi, nil)
	assert.Nil(t, err)
	siblings := merkletree.SiblingsFromProof(proof)
	for i := len(siblings); i < clt.MaxLevels(); i++ { // add the rest of empty levels to the siblings
		siblings = append(siblings, &merkletree.HashZero)
	}
	siblings = append(siblings, &merkletree.HashZero) // add extra level for circom compatibility
	var siblingsStr []string
	for i := 0; i < len(siblings); i++ {
		siblingsStr = append(siblingsStr, new(big.Int).SetBytes(common3.SwapEndianness(siblings[i].Bytes())).String())
	}
	jsonSiblings, err := json.Marshal(siblingsStr)
	assert.Nil(t, err)

	// newIdState
	newIdState := new(big.Int).SetBytes(common3.SwapEndianness(id.Bytes()))

	var out string
	out += fmt.Sprintln("{")
	out += fmt.Sprintf(`"id": "%s",`+"\n", new(big.Int).SetBytes(common3.SwapEndianness(id.Bytes())))
	out += fmt.Sprintf(`"oldIdState": "%s",`+"\n", "0")
	out += fmt.Sprintf(`"userPrivateKey": "%s",`+"\n", skToBigInt(&k))
	out += fmt.Sprintf(`"siblings": %s,`+"\n", jsonSiblings)
	out += fmt.Sprintf(`"claimsTreeRoot": "%s",`+"\n", new(big.Int).SetBytes(common3.SwapEndianness(clt.RootKey().Bytes())))
	out += fmt.Sprintf(`"newIdState": "%s"`+"\n", newIdState) // TMP
	out += fmt.Sprintf("}")
	return out
}

func pruneBuffer(buf *[32]byte) *[32]byte {
	buf[0] = buf[0] & 0xF8
	buf[31] = buf[31] & 0x7F
	buf[31] = buf[31] | 0x40
	return buf
}

func skToBigInt(k *babyjub.PrivateKey) *big.Int {
	sBuf := babyjub.Blake512(k[:])
	sBuf32 := [32]byte{}
	copy(sBuf32[:], sBuf[:32])
	pruneBuffer(&sBuf32)
	s := new(big.Int)
	cryptoUtils.SetBigIntFromLEBytes(s, sBuf32[:])
	s.Rsh(s, 3)
	return s
}
