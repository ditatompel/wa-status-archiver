package database

// NOTE:
// This is not whatsmeow's database schema. This is our own schema for this app.
func CreateSchema(db *DB) error {
	// status updates table schema
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tbl_statuses
	(
		id SERIAL,
		message_id character varying(255) NOT NULL,
		our_jid character varying(255) NOT NULL,
		sender_jid character varying(255) NOT NULL,
		sender_name character varying(255) NOT NULL,
		caption text NOT NULL,
		media_type character varying(50) NOT NULL,
		mimetype character varying(50) NOT NULL,
		filesize bigint NOT NULL DEFAULT 0,
		height integer NOT NULL DEFAULT 0,
		width integer NOT NULL DEFAULT 0,
		file_location text NOT NULL,
		msg_date timestamp with time zone NOT NULL,
		CONSTRAINT tbl_status_updates_pkey PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS statuses_sender_jid_idx ON tbl_statuses (sender_jid)`)
	if err != nil {
		return err
	}

	// admin table schema
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tbl_admin
	(
		id SERIAL,
		username character varying(255) NOT NULL,
		password character varying(255) NOT NULL,
		lastactive_ts timestamp with time zone,
		created_ts timestamp with time zone NOT NULL,
		CONSTRAINT tbl_admin_pkey PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS admin_username_idx ON tbl_admin (username)`)
	if err != nil {
		return err
	}

	// chats table schema
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tbl_chats
	(
		message_id character varying(255) NOT NULL,
		room_id character varying(255) NOT NULL,
		our_jid character varying(255) NOT NULL,
		sender_jid character varying(255) NOT NULL,
		sender_name character varying(255) NOT NULL,
		is_group boolean NOT NULL,
		is_from_me boolean NOT NULL,
		msg_type character varying(50) NOT NULL,
		media_type character varying(50) NOT NULL,
		msg_conversation text NOT NULL,
		category character varying(50) NOT NULL,
		msg_date timestamp with time zone NOT NULL,
		CONSTRAINT tbl_chats_pkey PRIMARY KEY (message_id, room_id, our_jid, sender_jid)
	)`)
	if err != nil {
		return err
	}

	return nil
}
