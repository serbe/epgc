package epgc

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

// Training - struct for training
type Training struct {
	ID        int64     `sql:"id" json:"id" `
	StartDate time.Time `sql:"start_date" json:"start-date"`
	EndDate   time.Time `sql:"end_date" json:"end-date"`
	StartStr  string    `sql:"-" json:"start-str"`
	EndStr    string    `sql:"-" json:"end-str"`
	Note      string    `sql:"note, null" json:"note"`
	CreatedAt time.Time `sql:"created_at" json:"created_at"`
	UpdatedAt time.Time `sql:"updated_at" json:"updated_at"`
}

func scanTraining(row *sql.Row) (Training, error) {
	var (
		sID        sql.NullInt64
		sStartDate pq.NullTime
		sEndDate   pq.NullTime
		sNote      sql.NullString
		training   Training
	)
	err := row.Scan(&sID, &sStartDate, &sEndDate, &sNote)
	if err != nil {
		log.Println("scanTraining row.Scan ", err)
		return training, err
	}
	training.ID = n2i(sID)
	training.StartDate = n2d(sStartDate)
	training.EndDate = n2d(sEndDate)
	training.Note = n2s(sNote)
	return training, nil
}

func scanTrainings(rows *sql.Rows, opt string) ([]Training, error) {
	var trainings []Training
	for rows.Next() {
		var (
			sID        sql.NullInt64
			sStartDate pq.NullTime
			sEndDate   pq.NullTime
			sNote      sql.NullString
			training   Training
		)
		switch opt {
		case "list":
			err := rows.Scan(&sID, &sStartDate, &sEndDate, &sNote)
			if err != nil {
				log.Println("scanTrainings rows.Scan list ", err)
				return trainings, err
			}
		}
		training.ID = n2i(sID)
		training.StartDate = n2d(sStartDate)
		training.EndDate = n2d(sEndDate)
		training.Note = n2s(sNote)
		trainings = append(trainings, training)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanTrainings rows.Err ", err)
	}
	return trainings, err
}

// GetTraining - get training by id
func (e *Edb) GetTraining(id int64) (Training, error) {
	if id == 0 {
		return Training{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		id,
		start_date
		end_date,
		note
	FROM
		trainings
	ORDER BY
		start_date
	WHERE id = $1`)
	if err != nil {
		log.Println("GetTraining e.db.Prepare ", err)
		return Training{}, err
	}
	row := stmt.QueryRow(id)
	training, err := scanTraining(row)
	return training, err
}

// GetTrainingList - get all training for list
func (e *Edb) GetTrainingList() ([]Training, error) {
	rows, err := e.db.Query(`SELECT
		id,
		start_date
		end_date,
		note
	FROM
		trainings
	ORDER BY
		start_date`)
	if err != nil {
		log.Println("GetTrainingList e.db.Query ", err)
		return []Training{}, err
	}
	trainings, err := scanTrainings(rows, "list")
	if err != nil {
		log.Println("GetTrainingList scanTrainings ", err)
		return []Training{}, err
	}
	for i := range trainings {
		trainings[i].StartStr = setStrMonth(trainings[i].StartDate)
		trainings[i].EndStr = setStrMonth(trainings[i].EndDate)
	}
	return trainings, err
}

// CreateTraining - create new training
func (e *Edb) CreateTraining(training Training) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO trainings(start_date, end_date, note, created_at) VALUES($1, $2, $3, now()) RETURNING id`)
	if err != nil {
		log.Println("CreateTraining e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(d2n(training.StartDate), d2n(training.EndDate), s2n(training.Note)).Scan(&training.ID)
	return training.ID, err
}

// UpdateTraining - save changes to training
func (e *Edb) UpdateTraining(training Training) error {
	stmt, err := e.db.Prepare(`UPDATE trainings SET start_date = $2, end_date = $3, note = $4, updated_at = now() WHERE id = $1`)
	if err != nil {
		log.Println("UpdateTraining e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(training.ID, d2n(training.StartDate), d2n(training.EndDate), s2n(training.Note))
	return err
}

// DeleteTraining - delete training by id
func (e *Edb) DeleteTraining(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE * FROM trainings WHERE id = $1`, id)
	if err != nil {
		log.Println("DeleteTraining ", err)
	}
	return err
}

func (e *Edb) trainingCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS trainings (id bigserial primary key, start_date date, end_date date, note text, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("trainingCreateTable ", err)
	}
	return err
}
