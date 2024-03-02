package repo

import (
	"fmt"

	"github.com/ditatompel/wa-status-archiver/internal/database"
)

type StatusUpdate struct {
	Id           int    `db:"id"`
	MessageId    string `db:"message_id"`
	Account      string `db:"account"`
	SenderJid    string `db:"sender_jid"`
	SenderName   string `db:"sender_name"`
	Caption      string `db:"caption"`
	MediaType    string `db:"media_type"`
	Mimetype     string `db:"mimetype"`
	Filesize     int    `db:"filesize"`
	Height       int    `db:"height"`
	Width        int    `db:"width"`
	FileLocation string `db:"file_location"`
	MsgDate      string `db:"msg_date"`
}

type StatusUpdates struct {
	CurrentPage int
	RowsPerPage int
	NextPage    int
	Statuses    []*StatusUpdate
}

type StatusUpdateQueryParams struct {
	JID         string
	RowsPerPage int
	Page        int
}

type StatusUpdateRepo struct {
	db *database.DB
}

type StatusUpdateRepository interface {
	Contacts() ([]contacts, error)
	StatusUpdates(q StatusUpdateQueryParams) (StatusUpdates, error)
}

func NewStatusUpdateRepo(db *database.DB) StatusUpdateRepository {
	return &StatusUpdateRepo{db}
}

type contacts struct {
	JID      string `db:"jid"`
	PushName string `db:"push_name"`
}

func (repo *StatusUpdateRepo) Contacts() ([]contacts, error) {
	query := `SELECT status.sender_jid AS jid, contact.push_name
	FROM tbl_status_updates AS status
	LEFT JOIN whatsmeow_contacts AS contact ON status.sender_jid = contact.their_jid
	GROUP BY jid, push_name
	ORDER BY push_name ASC`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []contacts{}
	for rows.Next() {
		contact := contacts{}
		if err := rows.Scan(&contact.JID, &contact.PushName); err != nil {
			return nil, err
		}
		res = append(res, contact)
	}
	return res, nil
}

func (repo *StatusUpdateRepo) StatusUpdates(q StatusUpdateQueryParams) (StatusUpdates, error) {
	queryParams := []interface{}{}
	where := ""

	queryParams = append(queryParams, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

	if q.JID != "" {
		where = "WHERE sender_jid = $3"
		queryParams = append(queryParams, q.JID)
	}

	statusUpdates := StatusUpdates{}
	query := fmt.Sprintf(`SELECT
		id, message_id, our_jid, sender_jid, sender_name, caption, media_type,
		mimetype, filesize, height, width, file_location, msg_date
	FROM tbl_status_updates %s ORDER BY msg_date DESC LIMIT $1 OFFSET $2`, where)

	rows, err := repo.db.Query(query, queryParams...)
	if err != nil {
		return statusUpdates, err
	}
	defer rows.Close()

	statusUpdates.CurrentPage = q.Page
	statusUpdates.RowsPerPage = q.RowsPerPage
	statusUpdates.NextPage = q.Page + 1

	for rows.Next() {
		statusUpdate := StatusUpdate{}
		err = rows.Scan(
			&statusUpdate.Id,
			&statusUpdate.MessageId,
			&statusUpdate.Account,
			&statusUpdate.SenderJid,
			&statusUpdate.SenderName,
			&statusUpdate.Caption,
			&statusUpdate.MediaType,
			&statusUpdate.Mimetype,
			&statusUpdate.Filesize,
			&statusUpdate.Height,
			&statusUpdate.Width,
			&statusUpdate.FileLocation,
			&statusUpdate.MsgDate,
		)
		if err != nil {
			return statusUpdates, err
		}

		statusUpdates.Statuses = append(statusUpdates.Statuses, &statusUpdate)
	}
	return statusUpdates, nil
}
