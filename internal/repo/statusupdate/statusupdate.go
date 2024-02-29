package statusupdate

import (
	"fmt"
	"slices"
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

type StatusUpdates struct {
	CurrentPage int
	RowsPerPage int
	NextPage    int
	Statuses    []*StatusUpdate
}

type StatusUpdateQueryParams struct {
	Search      string `json:"search"`
	Sort        string `json:"sort"`
	Dir         string `json:"dir"`
	RowsPerPage int    `json:"rows_per_page"`
	Page        int    `json:"page"`
}

type StatusUpdateRepo struct {
	db *database.DB
}

type StatusUpdateRepository interface {
	StatusUpdates(q StatusUpdateQueryParams) (StatusUpdates, error)
}

func NewStatusUpdateRepo(db *database.DB) StatusUpdateRepository {
	return &StatusUpdateRepo{db}
}

func (repo *StatusUpdateRepo) StatusUpdates(q StatusUpdateQueryParams) (StatusUpdates, error) {
	queryParams := []interface{}{}
	where := ""

	queryParams = append(queryParams, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

	if q.Search != "" {
		where = "WHERE sender_name ILIKE $3 OR caption ILIKE $4"
		queryParams = append(queryParams, "%"+q.Search+"%", "%"+q.Search+"%")
	}

	allowedSorts := []string{"their_jid", "push_name"}
	if !slices.Contains(allowedSorts, q.Sort) {
		q.Sort = "msg_date"
	}

	sortDir := "DESC"
	if q.Dir == "asc" {
		sortDir = q.Dir
	}

	statusUpdates := StatusUpdates{}
	query := fmt.Sprintf(`SELECT id, message_id, our_jid, sender_jid, sender_name, caption, media_type, mimetype, filesize, height, width, file_location, msg_date FROM tbl_status_updates %s ORDER BY %s %s LIMIT $1 OFFSET $2`, where, q.Sort, sortDir)

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
