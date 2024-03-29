package repo

import (
	"database/sql"
	"fmt"

	"github.com/ditatompel/wa-status-archiver/internal/database"
)

type Contact struct {
	OurJID       string         `db:"our_jid"`
	TheirJID     string         `db:"their_jid" `
	FirstName    sql.NullString `db:"first_name"`
	FullName     sql.NullString `db:"full_name"`
	PushName     sql.NullString `db:"push_name"`
	BusinessName sql.NullString `db:"business_name"`
}

type Contacts struct {
	CurrentPage int
	RowsPerPage int
	NextPage    int
	Contacts    []*Contact
}

type ContactQueryParams struct {
	Search      string
	RowsPerPage int
	Page        int
}

type ContactRepo struct {
	db *database.DB
}

type ContactRepository interface {
	Contacts(q ContactQueryParams) (Contacts, error)
}

func NewContactRepo(db *database.DB) ContactRepository {
	return &ContactRepo{db}
}

func (repo *ContactRepo) Contacts(q ContactQueryParams) (Contacts, error) {
	queryParams := []interface{}{}
	where := ""

	queryParams = append(queryParams, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

	if q.Search != "" {
		where = "WHERE first_name ILIKE $3 OR full_name ILIKE $4 OR push_name ILIKE $5 OR business_name ILIKE $6 OR their_jid ILIKE $7"
		queryParams = append(queryParams, "%"+q.Search+"%", "%"+q.Search+"%", "%"+q.Search+"%", "%"+q.Search+"%", "%"+q.Search+"%")
	}

	contacts := Contacts{}

	query := fmt.Sprintf(`SELECT
		our_jid, their_jid, first_name, full_name, push_name, business_name
		FROM whatsmeow_contacts %s ORDER BY push_name DESC LIMIT $1 OFFSET $2`, where)

	rows, err := repo.db.Query(query, queryParams...)
	if err != nil {
		return contacts, err
	}
	defer rows.Close()

	contacts.CurrentPage = q.Page
	contacts.RowsPerPage = q.RowsPerPage
	contacts.NextPage = q.Page + 1

	for rows.Next() {
		contact := Contact{}
		err = rows.Scan(&contact.OurJID, &contact.TheirJID, &contact.FirstName, &contact.FullName, &contact.PushName, &contact.BusinessName)
		if err != nil {
			return contacts, err
		}
		contacts.Contacts = append(contacts.Contacts, &contact)
	}

	if err = rows.Err(); err != nil {
		return contacts, err
	}

	return contacts, nil
}
