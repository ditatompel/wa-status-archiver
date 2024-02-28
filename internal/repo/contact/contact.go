package contact

import (
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

type ContactRepo struct {
	db *database.DB
}

type ContactRepository interface {
	Contacts() ([]*Contact, error)
}

func NewContactRepo(db *database.DB) ContactRepository {
	return &ContactRepo{db}
}

func (repo *ContactRepo) Contacts() ([]*Contact, error) {
	query := `SELECT our_jid, their_jid, first_name, full_name,push_name, business_name FROM whatsmeow_contacts`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := []*Contact{}
	for rows.Next() {
		contact := &Contact{}
		err = rows.Scan(&contact.OurJID, &contact.TheirJID, &contact.FirstName, &contact.FullName, &contact.PushName, &contact.BusinessName)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
