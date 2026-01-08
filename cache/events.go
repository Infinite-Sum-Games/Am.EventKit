package cache

import "github.com/google/uuid"

func CacheAllEvents() (bool, error) {
	return true, nil
}

func FetchAllEventsCache() {
}

func EvictEventCache(eventId uuid.UUID) (bool, error) {
	return false, nil
}
