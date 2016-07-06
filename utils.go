package epgc

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

func s2n(val string) sql.NullString {
	var s2n sql.NullString
	if val == "" {
		s2n.Valid = false
	} else {
		s2n.Valid = true
		s2n.String = val
	}
	return s2n
}

func d2n(val time.Time) pq.NullTime {
	var d2n pq.NullTime
	if val.Format("02.01.2006") == "01.01.0001" {
		d2n.Valid = false
	} else {
		d2n.Valid = true
		d2n.Time = val
	}
	return d2n
}

func i2n(val int64) sql.NullInt64 {
	var i2n sql.NullInt64
	if val == 0 {
		i2n.Valid = false
	} else {
		i2n.Valid = true
		i2n.Int64 = val
	}
	return i2n
}

func b2n(val bool) sql.NullBool {
	var b2n sql.NullBool
	if val == false {
		b2n.Valid = false
	} else {
		b2n.Valid = true
		b2n.Bool = val
	}
	return b2n
}

func n2s(val sql.NullString) string {
	return val.String
}

func n2i(val sql.NullInt64) int64 {
	return val.Int64
}

func n2b(val sql.NullBool) bool {
	return val.Bool
}

func n2d(val pq.NullTime) time.Time {
	return val.Time
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
		ee = append(ee, email)
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
		i, err := strconv.ParseInt(p, 10, 64)
		if err == nil {
			var phone Phone
			phone.Phone = i
			phone.Fax = false
			pp = append(pp, phone)
		}
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
		i, err := strconv.ParseInt(f, 10, 64)
		if err == nil {
			var fax Phone
			fax.Phone = i
			fax.Fax = true
			ff = append(ff, fax)
		}
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
		pp = append(pp, practice)
	}
	return pp
}

// func toInt(b bool) int {
// 	if b {
// 		return 1
// 	}
// 	return 0
// }

func int64InSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func setStrMonth(d time.Time) (output string) {
	str := d.Format("02.01.2006")
	spl := strings.Split(str, ".")
	month := map[string]string{"01": "января", "02": "февраля", "03": "марта", "04": "апреля", "05": "мая", "06": "июня", "07": "июля", "08": "августа", "09": "сентября", "10": "октября", "11": "ноября", "12": "декабря "}
	output = spl[0] + " " + month[spl[1]] + " " + spl[2] + " года"
	return
}
