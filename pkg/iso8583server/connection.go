package iso8583server

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"net"
)

type Connection struct {
	net.Conn

	// ID assigned by the server to this connection.
	Id string
	// Active represents if the server is trying to read from this connection
	Active bool
	// LastRead represents the last time a message read was successfully performed.
	LastRead time.Time
	// LastWrite represents the last time a message write was successfully performed.
	LastWrite time.Time
}

type connectionDB struct {
	m       map[string]*Connection
	rwMutex sync.RWMutex
}

func (db *connectionDB) Get(id string) (Connection, bool) {
	db.rwMutex.RLock()
	c, ok := db.m[id]
	db.rwMutex.RUnlock()
	if !ok {
		return Connection{}, false
	}
	return *c, true
}

// Delete deletes a map entry, it returns true only on success.
// Delete only can delete entries where the connection is no longer active.
func (db *connectionDB) Delete(id string) bool {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	c, ok := db.m[id]
	if !ok {
		return true
	}
	if c.Active {
		return false
	}
	delete(db.m, id)
	return true
}

// delete has no such restrictions such as Delete.
func (db *connectionDB) deleteOldest(limit int) {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()

	var oldestID string
	var oldestTime time.Time
	var deactivatedConnections int

	first := true
	for _, c := range db.m {
		if c.Active {
			continue
		}
		deactivatedConnections++

		if first {
			oldestID = c.Id
			if c.LastWrite.After(c.LastRead) {
				oldestTime = c.LastRead
				first = false
				continue
			}
			oldestTime = c.LastWrite
			first = false
			continue
		}
		if c.LastWrite.Before(oldestTime) {
			oldestID = c.Id
			oldestTime = c.LastWrite
		}
		if c.LastRead.Before(oldestTime) {
			oldestID = c.Id
			oldestTime = c.LastWrite
		}
	}

	if deactivatedConnections < limit {
		return
	}

	delete(db.m, oldestID)
}

// set adds a entry
func (db *connectionDB) set(entry *Connection) {
	db.rwMutex.Lock()
	db.m[entry.Id] = entry
	db.rwMutex.Unlock()
}

func (db *connectionDB) GetAll() []Connection {
	var out []Connection

	db.rwMutex.RLock()
	for _, c := range db.m {
		out = append(out, *c)
	}
	db.rwMutex.RUnlock()

	return out
}

func (server *Server) newConnection(c net.Conn) (*Connection, error) {
	con := &Connection{
		Conn:      c,
		Active:    true,
		LastRead:  time.Now(),
		LastWrite: time.Now(),
	}

	id, err := server.config.ConnIdGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed generating id for connection: %w", err)
	}

	con.Id = id

	server.Connections.set(con)

	return con, nil
}

func defaultIdGenerator() (string, error) {
	n := strconv.Itoa(99999999999999 + rand.Intn(99999999999999-10000000000000))
	return n[:5] + "-" + n[5:10] + "-" + n[10:], nil
}
