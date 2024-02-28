package contact

import (
	"wabot/internal/database"
)

type Contact struct {
	Phone  string `db:"phone" json:"phone"`
	JID    string `db:"jid" json:"jid"`
	Server string `db:"server" json:"server"`
	Name   string `db:"name" json:"name"`
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
	query := `SELECT phone, jid, server, name FROM tbl_users`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := []*Contact{}
	for rows.Next() {
		contact := &Contact{}
		err = rows.Scan(&contact.Phone, &contact.JID, &contact.Server, &contact.Name)
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
