package app

import (
	"encoding/gob"
	"os"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	Store *sessions.FilesystemStore
	DB    *bolt.DB
)

type User struct {
	Email    string
	Username string
	Role     RoleType
}

type RoleType int

const (
	Admin = iota
	ReadOnly
)

// Init sets up the user database and session store. the Store provides a map
// between session IDs and users; from the user ID (email) contained in the
// session (and validated by an upstream identity provider), we can get the
// stored user and their attributes from the database
func Init() error {
	log.SetLevel(log.DebugLevel)
	db, err := bolt.Open("/tmp/bolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("error initializing database. %s", err)
	}
	// 1st key - auth key - should be 32 or 64 bytes
	// 2nd key - encryption key - 16, 24, 32 -> AES-128, AES-192, AES-256 (optional)
	Store = sessions.NewFilesystemStore(
		os.Getenv("APP_SESSION_STORE_PATH"),
		[]byte(os.Getenv("APP_SESSION_KEY")),
		[]byte(os.Getenv("APP_SESSION_ENCRYPTION_KEY")),
	)
	DB = db
	gob.Register(map[string]interface{}{})

	return nil
}
