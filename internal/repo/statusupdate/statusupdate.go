package statusupdate

import (
	"wabot/internal/database"
)

type StatusUpdate struct {
	Id           int    `json:"id" db:"id"`
	MessageId    string `json:"message_id" db:"message_id"`
	Account      string `json:"account" db:"account"`
	SenderJid    string `json:"sender_jid" db:"sender_jid"`
	SenderName   string `json:"sender_name" db:"sender_name"`
	Caption      string `json:"caption" db:"caption"`
	MediaType    string `json:"media_type" db:"media_type"`
	Mimetype     string `json:"mimetype" db:"mimetype"`
	Filesize     int    `json:"filesize" db:"filesize"`
	Height       int    `json:"height" db:"height"`
	Width        int    `json:"width" db:"width"`
	FileLocation string `json:"file_location" db:"file_location"`
	MsgDate      string `json:"msg_date" db:"msg_date"`
}

type StatusUpdateRepo struct {
	db *database.DB
}

type StatusUpdateRepository interface {
	StatusUpdates() ([]*StatusUpdate, error)
}

func NewStatusUpdateRepo(db *database.DB) StatusUpdateRepository {
	return &StatusUpdateRepo{db}
}

func (repo *StatusUpdateRepo) StatusUpdates() ([]*StatusUpdate, error) {
	query := `SELECT id, message_id, our_jid, sender_jid, sender_name, caption, media_type, mimetype, filesize, height, width, file_location, msg_date FROM tbl_status_updates`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusUpdates := []*StatusUpdate{}
	for rows.Next() {
		statusUpdate := &StatusUpdate{}
		err := rows.Scan(
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
			return nil, err
		}
		statusUpdates = append(statusUpdates, statusUpdate)

	}
	return statusUpdates, nil
}
