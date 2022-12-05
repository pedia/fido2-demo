package main

import (
	"testing"

	"github.com/duo-labs/webauthn/protocol"
	wan "github.com/duo-labs/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	is := assert.New(t)

	store := NewStore()
	is.NotNil(store)

	u := store.GetUser("foo")
	is.NotNil(u)
	is.Equal("foo", u.Name)
	is.Greater(u.Uid, uint32(0))
	is.False(u.Status)
	is.False(u.Ctime.IsZero())

	sd := wan.SessionData{Challenge: "UILZsLXDM92UVGiKXrPT4-LRm0M6w4MCSEoDuy57hjg",
		UserID:               []uint8{0x69, 0x64, 0x2d, 0x66},
		AllowedCredentialIDs: [][]uint8(nil),
		UserVerification:     "required", Extensions: protocol.AuthenticationExtensions(nil)}

	err := store.SaveSession(&sd)
	is.Nil(err)

	sds := store.GetSession(u.Name)
	is.NotEmpty(sds)

	c := &wan.Credential{
		ID:              []uint8{0xc1, 0x9c, 0xd8, 0xd1, 0x68, 0xc, 0xb6, 0x30, 0xa0, 0x3a, 0xa1, 0x7c, 0x3c, 0x6c, 0x59, 0xad},
		PublicKey:       []uint8{0xa5, 0x1, 0x2, 0x3, 0x26, 0x20, 0x1, 0x21, 0x58, 0x20, 0x44, 0xef, 0xc2, 0x64, 0x33, 0xb2, 0x57, 0x31, 0x95, 0xbd, 0xaf, 0xd0, 0x5a, 0x32, 0x0, 0x8f, 0x0, 0x52, 0x7, 0x5a, 0xe1, 0xcc, 0xc7, 0xa3, 0x19, 0x4f, 0xf, 0xab, 0xc6, 0x7c, 0xb4, 0x2e, 0x22, 0x58, 0x20, 0x86, 0x55, 0x60, 0x34, 0xd1, 0x67, 0x22, 0x9a, 0x25, 0xdd, 0x24, 0x93, 0x23, 0x61, 0x4, 0x2c, 0x6b, 0xad, 0x28, 0xae, 0x88, 0x75, 0x3, 0xc3, 0xea, 0xff, 0xa3, 0x65, 0x71, 0x47, 0x74, 0xd4},
		AttestationType: "none",
		Authenticator: wan.Authenticator{
			AAGUID:       []uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			SignCount:    0x0,
			CloneWarning: false,
		},
	}
	erc := store.SaveCredential(u, c)
	is.Nil(erc)

	// store.GetCredential
}