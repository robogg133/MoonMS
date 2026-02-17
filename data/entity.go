package data

import (
	"sync"
)

type EntityID uint16

const InvalidEntityID EntityID = 0

var entityMu sync.RWMutex

// string → ID
var entityByName = map[string]EntityID{}

// ID → string (debug / save / protocol)
var entityByID = []string{
	"", // 0 index reserved as invalid
}

func IdentifierToEntityID(id string) EntityID {
	entityMu.RLock()
	defer entityMu.RUnlock()

	if eid, ok := entityByName[id]; ok {
		return eid
	}
	return InvalidEntityID
}

func RegisterNewEntity(id string) EntityID {
	entityMu.Lock()
	defer entityMu.Unlock()

	if eid, exists := entityByName[id]; exists {
		return eid
	}

	eid := EntityID(len(entityByID))
	entityByName[id] = eid
	entityByID = append(entityByID, id)

	return eid
}
