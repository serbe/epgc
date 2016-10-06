package epgc

import (
	"database/sql"
	"log"
)

// Siren - struct for siren
type Siren struct {
	ID        int64     `sql:"id" json:"id"`
	NumID     int64     `sql:"num_id, null" json:"num_id"`
	NumPass   string    `sql:"num_pass, null" json:"num_pass"`
	TypeID    int64     `sql:"type_id" json:"type_id"`
	Type      SirenType `sql:"-"`
	Address   string    `sql:"address, null" json:"address"`
	Radio     string    `sql:"radio, null" json:"radio"`
	Desk      string    `sql:"desk, null" json:"desk"`
	ContactID int64     `sql:"contact_id, null" json:"contact_id"`
	Contact   Contact   `sql:"-"`
	CompanyID int64     `sql:"company_id, null" json:"company_id"`
	Company   Company   `sql:"-"`
	Latitude  string    `sql:"latitude, null" json:"latitude"`
	Longitude string    `sql:"longitude, null" json:"longitude"`
	Stage     int64     `sql:"stage, null" json:"stage"`
	Own       string    `sql:"own, null" json:"own"`
	Note      string    `sql:"note, null" json:"note"`
	CreatedAt string    `sql:"created_at" json:"created_at"`
	UpdatedAt string    `sql:"updated_at" json:"updated_at"`
}

func scanSiren(row *sql.Row) (Siren, error) {
	var (
		sID        sql.NullInt64
		sNumID     sql.NullInt64
		sNumPass   sql.NullString
		sTypeID    sql.NullInt64
		sAddress   sql.NullString
		sRadio     sql.NullString
		sDesk      sql.NullString
		sContactID sql.NullInt64
		sCompanyID sql.NullInt64
		sLatitude  sql.NullString
		sLongitude sql.NullString
		sStage     sql.NullInt64
		sOwn       sql.NullString
		sNote      sql.NullString
		siren      Siren
	)
	err := row.Scan(&sID, &sNumID, &sNumPass, &sTypeID, &sAddress, &sRadio, &sDesk, &sContactID, &sCompanyID, &sLatitude, &sLongitude, &sStage, &sOwn, &sNote)
	if err != nil {
		log.Println("scanSiren row.Scan ", err)
		return siren, err
	}
	siren.ID = n2i(sID)
	siren.NumID = n2i(sNumID)
	siren.NumPass = n2s(sNumPass)
	siren.TypeID = n2i(sTypeID)
	siren.Address = n2s(sAddress)
	siren.Radio = n2s(sRadio)
	siren.Desk = n2s(sDesk)
	siren.CompanyID = n2i(sContactID)
	siren.CompanyID = n2i(sCompanyID)
	siren.Latitude = n2s(sLatitude)
	siren.Longitude = n2s(sLongitude)
	siren.Stage = n2i(sStage)
	siren.Own = n2s(sOwn)
	siren.Note = n2s(sNote)
	return siren, nil
}

func scanSirensList(rows *sql.Rows) ([]Siren, error) {
	var sirens []Siren
	for rows.Next() {
		var (
			sID        sql.NullInt64
			sNumID     sql.NullInt64
			sNumPass   sql.NullString
			sTypeID    sql.NullInt64
			sAddress   sql.NullString
			sRadio     sql.NullString
			sDesk      sql.NullString
			sContactID sql.NullInt64
			sCompanyID sql.NullInt64
			sLatitude  sql.NullString
			sLongitude sql.NullString
			sStage     sql.NullInt64
			sOwn       sql.NullString
			sNote      sql.NullString
			siren      Siren
		)
		err := rows.Scan(&sID, &sNumID, &sNumPass, &sTypeID, &sAddress, &sRadio, &sDesk, &sContactID, &sCompanyID, &sLatitude, &sLongitude, &sStage, &sOwn, &sNote)
		if err != nil {
			log.Println("scanSirensList rows.Scan ", err)
			return sirens, err
		}
		siren.ID = n2i(sID)
		siren.NumID = n2i(sNumID)
		siren.NumPass = n2s(sNumPass)
		siren.TypeID = n2i(sTypeID)
		siren.Address = n2s(sAddress)
		siren.Radio = n2s(sRadio)
		siren.Desk = n2s(sDesk)
		siren.CompanyID = n2i(sContactID)
		siren.CompanyID = n2i(sCompanyID)
		siren.Latitude = n2s(sLatitude)
		siren.Longitude = n2s(sLongitude)
		siren.Stage = n2i(sStage)
		siren.Own = n2s(sOwn)
		siren.Note = n2s(sNote)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanSirensList rows.Err ", err)
	}
	return sirens, err
}

// GetSiren - get one siren by id
func (e *Edb) GetSiren(id int64) (Siren, error) {
	if id == 0 {
		return Siren{}, nil
	}
	row := e.db.QueryRow(`
		SELECT
			id,
			num_id,
			num_pass,
			type_id,
			address,
			radio,
			desk,
			contact_id,
			company_id,
			latitude,
			longitude,
			stage,
			own,
			note
		FROM
			sirens
		WHERE
			id = $1
	`, id)
	siren, err := scanSiren(row)
	return siren, err
}

// GetSirenList - get all siren for list
func (e *Edb) GetSirenList() ([]Siren, error) {
	rows, err := e.db.Query(`
		SELECT
			id,
			num_id,
			num_pass,
			type_id,
			address,
			radio,
			desk,
			contact_id,
			company_id,
			latitude,
			longitude,
			stage,
			own,
			note
		FROM
			sirens
		ORDER BY
			name ASC`)
	if err != nil {
		log.Println("GetSirenList e.db.Query ", err)
		return []Siren{}, err
	}
	sirens, err := scanSirensList(rows)
	return sirens, err
}

// CreateSiren - create new siren
func (e *Edb) CreateSiren(siren Siren) (int64, error) {
	stmt, err := e.db.Prepare(`
		INSERT INTO
			sirens (
				num_id,
				num_pass,
				type_id,
				address,
				radio,
				desk,
				contact_id,
				company_id,
				latitude,
				longitude,
				stage,
				own,
				note
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12,
				$13,
				now()
			)
		RETURNING
			id
	`)
	if err != nil {
		log.Println("CreateSiren e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(
		i2n(siren.NumID),
		s2n(siren.NumPass),
		i2n(siren.TypeID),
		s2n(siren.Address),
		s2n(siren.Radio),
		s2n(siren.Desk),
		i2n(siren.CompanyID),
		i2n(siren.CompanyID),
		s2n(siren.Latitude),
		s2n(siren.Longitude),
		i2n(siren.Stage),
		s2n(siren.Own),
		s2n(siren.Note)).Scan(&siren.ID)
	if err != nil {
		log.Println("CreateSiren db.QueryRow ", err)
	}
	return siren.ID, err
}

// UpdateSiren - save siren changes
func (e *Edb) UpdateSiren(siren Siren) error {
	stmt, err := e.db.Prepare(`
		UPDATE
			sirens
		SET
			num_id = $2,
			num_pass = $3,
			type_id = $4,
			address = $5,
			radio = $6,
			desk = $7,
			contact_id = $8,
			company_id = $9,
			latitude = $10,
			longitude = $11,
			stage = $12,
			own = $13,
			note = $14,
			updated_at = now()
		WHERE
			id = $1
	`)
	if err != nil {
		log.Println("UpdateSiren e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(
		i2n(siren.ID),
		i2n(siren.NumID),
		s2n(siren.NumPass),
		i2n(siren.TypeID),
		s2n(siren.Address),
		s2n(siren.Radio),
		s2n(siren.Desk),
		i2n(siren.CompanyID),
		i2n(siren.CompanyID),
		s2n(siren.Latitude),
		s2n(siren.Longitude),
		i2n(siren.Stage),
		s2n(siren.Own),
		s2n(siren.Note))
	if err != nil {
		log.Println("UpdateSiren stmt.Exec ", err)
	}
	return err
}

// DeleteSiren - delete siren by id
func (e *Edb) DeleteSiren(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			sirens
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteSiren e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) sirenCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			sirens (
				id         bigserial PRIMARY KEY,
				num_id     bigint,
				num_pass   text,
				type_id    bigint,
				address    text,
				radio      text,
				desk       text,
				contact_id bigint,
				company_id bigint,
				latitude   text,
				longitude  text,
				stage      bigint,
				own        text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone,
				UNIQUE(num_id, num_pass, type_id)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("sirenCreateTable e.db.Exec ", err)
	}
	return err
}
