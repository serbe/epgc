package epgc

import (
	"database/sql"
	"strings"
)

func s2n(s string) sql.NullString {
	var s2n sql.NullString
	if s == "" {
		s2n.Valid = false
	} else {
		s2n.Valid = true
		s2n.String = s
	}
	return s
}

func i2n(i int64) sql.NullInt64 {
	var i2n sql.NullInt64
	if i == 0 {
		i2n.Valid = false
	} else {
		i2n.Valid = true
		i2n.Int64 = i
	}
	return i2n
}

func n2s(n2s sql.NullString) string {
	return n2s.String
}

func n2i(n2i sql.NullInt64) int64 {
	return n2i.Int64
}

func n2emails(emails sql.NullString) []Email {
	var (
		e  string
		es []string
		ee []Email
	)
	e = n2s(emails)
	if e == "" {
		return ee
	}
	es = strings.Split(e, ",")
	for _, e = range es {
		var email Email
		email.Email = e
		append(ee, email)
	}
	return ee
}

func n2phones(phones sql.NullString) []Phone {
	var (
		p  string
		ps []string
		pp []Phone
	)
	p = n2s(phones)
	if p == "" {
		return pp
	}
	ps = strings.Split(p, ",")
	for _, p = range ps {
		var phone Phone
		phone.Phone = p
		phone.Fax = false
		append(pp, phone)
	}
	return pp
}

func n2faxes(faxes sql.NullString) []Phone {
	var (
		f  string
		fs []string
		ff []Phone
	)
	f = n2s(faxes)
	if f == "" {
		return ff
	}
	fs = strings.Split(f, ",")
	for _, f = range fs {
		var fax Phone
		fax.Phone = f
		fax.Fax = true
		append(ff, fax)
	}
	return ff
}

func n2practices(practices sql.NullString) []Practice {
	var (
		p  string
		ps []string
		pp []Practice
	)
	p = n2s(practices)
	if p == "" {
		return pp
	}
	ps = strings.Split(p, ",")
	for _, p = range ps {
		var practice Practice
		practice.Topic = p
		append(pp, practice)
	}
	return pp
}
