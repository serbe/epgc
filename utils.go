package epgc

import "database/sql"

func ns(s string) sql.NullString {
	var ns sql.NullString
	if s == "" {
		ns.Valid = false
	} else {
		ns.Valid = true
		ns.String = s
	}
	return s
}
