package contact

import (
	"fmt"
	"slices"
	"wabot/internal/database"
)

type Contact struct {
	OurJID       string              `db:"our_jid" json:"our_jid"`
	TheirJID     string              `db:"their_jid" json:"their_jid"`
	FirstName    database.NullString `db:"first_name" json:"first_name"`
	FullName     database.NullString `db:"full_name" json:"full_name"`
	PushName     database.NullString `db:"push_name" json:"push_name"`
	BusinessName database.NullString `db:"business_name" json:"business_name"`
}

type Contacts struct {
	CurrentPage int
	RowsPerPage int
	NextPage    int
	Contacts    []*Contact
}

type ContactQueryParams struct {
	Search      string `json:"search"`
	Sort        string `json:"sort"`
	Dir         string `json:"dir"`
	RowsPerPage int    `json:"rows_per_page"`
	Page        int    `json:"page"`
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

	allowedSorts := []string{"their_jid", "push_name"}
	if !slices.Contains(allowedSorts, q.Sort) {
		q.Sort = "push_name"
	}

	sortDir := "DESC"
	if q.Dir == "asc" {
		sortDir = q.Dir
	}

	query := fmt.Sprintf(`SELECT
		our_jid, their_jid, first_name, full_name, push_name, business_name
		FROM whatsmeow_contacts %s ORDER BY %s %s LIMIT $1 OFFSET $2`, where, q.Sort, sortDir)

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
