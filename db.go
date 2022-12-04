package main

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	wan "github.com/duo-labs/webauthn/webauthn"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id    uint32    `db:"id"`
	Name  string    `db:"name"`
	Ctime time.Time `db:"ctime"`
}

func (u User) WebAuthnID() []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, u.Id)
	return bs
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	// return cases.Title(u.Name)
	return strings.Title(u.Name)
}

func (u User) WebAuthnIcon() string {
	return fmt.Sprintf("%s/s/avatar.jpeg", site)
}

// Credentials owned by the user
func (u User) WebAuthnCredentials() []wan.Credential {
	return []wan.Credential{}
}

// TODO:
type Store interface {
	GetUser() wan.User
}

type user_store struct {
	db *sqlx.DB
}

func NewStore() *user_store {
	return &user_store{db: sqlx.MustConnect("sqlite3", ":memory:")}
}

func (*user_store) GetUser() wan.User {
	return User{}
}

type Credential struct {
	Id              uint64    `db:"id"`
	Uid             uint32    `db:"uid"`
	PublicKey       string    `db:"public_key"`
	AttestationType string    `db:"attestation_type"`
	Ctime           time.Time `db:"ctime"`
}
