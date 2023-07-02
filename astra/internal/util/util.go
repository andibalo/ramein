package util

import "github.com/gocql/gocql"

func ParseUUIDToGoCQLUUID(s string) (gocql.UUID, error) {
	u, err := gocql.ParseUUID(s)
	if err != nil {
		return u, err
	}
	return u, nil
}
