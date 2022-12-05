package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	wan "github.com/duo-labs/webauthn/webauthn"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

type User struct {
	Uid    uint32    `db:"uid"`
	Name   string    `db:"name"`
	Status bool      `db:"status"`
	Ctime  time.Time `db:"ctime"`
}

func (u User) CreateTableSql() string {
	return `CREATE TABLE IF NOT EXISTS user(
uid INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL UNIQUE,
status INTEGER NOT NULL DEFAULT 0,
ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`
}

func (u User) WebAuthnID() []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, u.Uid)
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
	cs := _store.GetCredential(u.Uid)
	return lo.Map(cs, func(v Credential, i int) wan.Credential {
		wc := wan.Credential{
			AttestationType: v.AttestationType,
		}
		wc.ID, _ = base64.RawURLEncoding.DecodeString(v.Id)
		wc.PublicKey, _ = base64.RawURLEncoding.DecodeString(v.PublicKey)
		return wc
	})
}

type Session struct {
	Id        uint32    `db:"id"`
	Uid       uint32    `db:"uid"`
	Challenge string    `db:"challenge"`
	Ctime     time.Time `db:"ctime"`
}

func (s Session) CreateTableSql() string {
	return `CREATE TABLE IF NOT EXISTS session(
id INTEGER PRIMARY KEY AUTOINCREMENT,
uid INTEGER NOT NULL,
challenge string NOT NULL,
ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`
}

type Credential struct {
	Id              string    `db:"id"`
	PublicKey       string    `db:"public_key"`
	AttestationType string    `db:"attestation_type"`
	Uid             uint32    `db:"uid"`
	Ctime           time.Time `db:"ctime"`
}

func (c Credential) CreateTableSql() string {
	return `CREATE TABLE IF NOT EXISTS credential(
id TEXT PRIMARY KEY,
public_key string NOT NULL,
attestation_type string NOT NULL,
uid INTEGER NOT NULL,
ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`
}

// TODO:
type Store interface {
	GetUser(name string) *User
	SaveSession(s *wan.SessionData) error
	GetSession(name string) []wan.SessionData
	SaveCredential(u *User, c *wan.Credential) error
}

type dummy_store struct {
	db *sqlx.DB
}

var _store *dummy_store

func NewStore() Store {
	// db := sqlx.MustConnect("sqlite3", ":memory:")
	db := sqlx.MustConnect("sqlite3", "s.db")

	db.MustExec(User{}.CreateTableSql())
	db.MustExec(Session{}.CreateTableSql())
	db.MustExec(Credential{}.CreateTableSql())

	// HACK:
	_store = &dummy_store{db: db}
	return _store
}

// create or get user
func (d *dummy_store) GetUser(name string) *User {
	var u User
	if err := d.db.Get(&u, "SELECT * FROM user WHERE name=?", name); err == nil {
		return &u
	}

	if res, err := d.db.Exec("INSERT INTO user(name) VALUES(?)", name); err == nil {
		id, _ := res.LastInsertId()
		if err := d.db.Get(&u, "SELECT * FROM user WHERE uid=?", id); err == nil {
			return &u
		}
	}
	return nil
}

func (d *dummy_store) SaveSession(s *wan.SessionData) error {
	sess := Session{
		Uid:       binary.LittleEndian.Uint32(s.UserID),
		Challenge: s.Challenge,
	}
	res, err := d.db.NamedExec(
		"INSERT INTO session(uid,challenge) VALUES(:uid,:challenge)",
		sess,
	)
	_ = res
	return err
}

func uint32_to_bytes(i uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, i)
	return bs
}

func (d *dummy_store) GetSession(name string) []wan.SessionData {
	var ss []Session
	err := d.db.Select(
		&ss,
		`SELECT session.* FROM session JOIN user ON user.uid = session.uid WHERE user.name=? ORDER BY id desc`,
		name,
	)
	if err != nil {
		return nil
	}

	return lo.Map(ss, func(v Session, i int) wan.SessionData {
		return wan.SessionData{
			Challenge:            v.Challenge,
			UserID:               uint32_to_bytes(v.Uid),
			AllowedCredentialIDs: nil,
			UserVerification:     protocol.VerificationRequired,
			Extensions:           nil,
		}
	})
}

func (d *dummy_store) SaveCredential(u *User, c *wan.Credential) error {
	lc := Credential{
		Id:              base64.RawURLEncoding.EncodeToString(c.ID),
		PublicKey:       base64.RawURLEncoding.EncodeToString(c.PublicKey),
		AttestationType: c.AttestationType,
		Uid:             u.Uid,
	}
	res, err := d.db.NamedExec(
		"INSERT INTO credential(id, public_key, attestation_type, uid) VALUES(:id,:public_key,:attestation_type,:uid)",
		lc,
	)
	_ = res
	return err
}

func (d *dummy_store) GetCredential(uid uint32) []Credential {
	var cs []Credential
	err := d.db.Select(
		&cs,
		`SELECT credential.* FROM credential JOIN user ON user.uid = credential.uid WHERE user.uid=? ORDER BY id desc`,
		uid,
	)
	if err != nil {
		return nil
	}

	return cs
}
