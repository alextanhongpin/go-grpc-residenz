package database

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

// Config represents the database config
type Config struct {
	Addrs    []string
	Timeout  time.Duration
	Database string
	Username string
	Password string
}

// Database holds the database session
type Database struct {
	Session  *mgo.Session
	Database string
}

// CopySession creates a new session with the same parameters as the original session,
// but preserves the exact authentication information from the original session
func (db *Database) CopySession() *mgo.Session {
	return db.Session.Copy()
}

// Collection returns a collection by name
func (db *Database) Collection(session *mgo.Session, collection string) *mgo.Collection {
	return session.Copy().DB(db.Database).C(collection)
}

// New creates a new database connection
func New(cfg Config) (*Database, error) {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    cfg.Addrs,
		Timeout:  60 * time.Second,
		Database: cfg.Database,
		Username: cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return &Database{
		Session:  session,
		Database: cfg.Database,
	}, nil
}
