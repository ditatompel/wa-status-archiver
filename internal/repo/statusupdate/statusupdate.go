package statusupdate

import (
	"wabot/internal/database"
)

type StatusUpdate struct {
	Id          int    `json:"id" db:"id"`
	MsgId       string `json:"msg_id" db:"msg_id"`
	Account     string `json:"account" db:"account"`
	SenderPhone string `json:"sender_phone" db:"sender_phone"`
	SenderJid   string `json:"sender_jid" db:"sender_jid"`
	SenderName  string `json:"sender_name" db:"sender_name"`
	Caption     string `json:"caption" db:"caption"`
	MediaType   string `json:"media_type" db:"media_type"`
	Mimetype    string `json:"mimetype" db:"mimetype"`
	Filesize    int    `json:"filesize" db:"filesize"`
	Height      int    `json:"height" db:"height"`
	Width       int    `json:"width" db:"width"`
	FileLoc     string `json:"file_loc" db:"file_loc"`
	MsgDate     string `json:"msg_date" db:"msg_date"`
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
	query := `SELECT id, msg_id, account, sender_phone, sender_jid, sender_name, caption, media_type, mimetype, filesize, height, width, file_loc, msg_date FROM tbl_status_updates`

	rows, err := repo.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusUpdates := []*StatusUpdate{}
	for rows.Next() {
		statusUpdate := &StatusUpdate{}
		err = rows.StructScan(statusUpdate)
		if err != nil {
			return nil, err
		}
		statusUpdates = append(statusUpdates, statusUpdate)
	}
	return statusUpdates, nil
}
