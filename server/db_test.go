package main

import (
	"encoding/base64"
	"encoding/binary"
	"testing"

	"github.com/duo-labs/webauthn/protocol"
	wan "github.com/duo-labs/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	is := assert.New(t)

	// db.Uid        =>    webauthn.User.WebAuthnID()
	// 1 => 0000 0001
	// json encode: AgAAAA==
	// wrong: stringToBuffer('AgAAAA==')
	// base64URLStringToBuffer('AgAAAA==') => ArrayBuffer[2,0,0,0]

	bs, erb := base64.RawURLEncoding.DecodeString("QVFBQUFBPT0")
	is.Nil(erb)
	i1 := binary.LittleEndian.Uint32(bs)
	is.Equal(uint32(0x41415141), i1)
}

func TestStore(t *testing.T) {
	is := assert.New(t)

	store := NewStore()
	is.NotNil(store)

	u := store.GetUser("bar")
	is.NotNil(u)
	is.Equal("bar", u.Name)
	is.Greater(u.Uid, uint32(0))
	is.False(u.Status)
	is.False(u.Ctime.IsZero())

	sd := wan.SessionData{Challenge: "UILZsLXDM92UVGiKXrPT4-LRm0M6w4MCSEoDuy57hjg",
		UserID:               uint32_to_bytes(u.Uid),
		AllowedCredentialIDs: [][]uint8(nil),
		UserVerification:     "required", Extensions: protocol.AuthenticationExtensions(nil)}

	sid, err := store.SaveSession(&sd)
	is.Nil(err)
	is.NotEqual(sid, 0)

	u2, sd2 := store.GetSession(sid)
	is.NotEmpty(sd2)
	is.Equal(sd.Challenge, sd2.Challenge)

	is.NotEmpty(u2)
	is.Equal(u.Name, u2.Name)
	is.Equal(u.Uid, u2.Uid)

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
