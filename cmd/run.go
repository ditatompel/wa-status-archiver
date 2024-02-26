package cmd

import (
	"context"
	"fmt"
	"mime"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"wabot/internal/database"
	"wabot/internal/helpers"

	"github.com/gosimple/slug"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/spf13/cobra"
)

var (
	cli  *whatsmeow.Client
	wLog waLog.Logger
)

type waRepo struct {
	db      *database.DB
	account *store.Device
}

func newWaRepo(db *database.DB, account *store.Device) *waRepo {
	return &waRepo{
		db:      db,
		account: account,
	}
}

var repo *waRepo

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Long:  `Run the bot by listening to WhatsApp websocket events`,
	Run: func(_ *cobra.Command, _ []string) {
		wLog = waLog.Stdout("Main", LogLevel, true)
		WaClient := CreateClient()
		ConnectClient(WaClient)

		WaClient.AddEventHandler(HandleEvent)

		WaClient.Connect()

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		WaClient.Disconnect()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func CreateClient() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", LogLevel, true)
	sql, err := sqlstore.New("sqlite3", "file:data/db/accounts.db?_foreign_keys=on", dbLog)
	if err != nil {
		dbLog.Errorf("Error connecting to database: %v", err)
		os.Exit(1)
	}

	deviceStore, err := sql.GetFirstDevice()
	if err != nil {
		wLog.Errorf("Error getting device: %v", err)
		os.Exit(1)
	}

	repo = newWaRepo(database.GetDB(), deviceStore)

	wLog.Infof("Device: %s", helpers.PrettyPrint(deviceStore.ID))
	wLog.Infof("Pushname: %s", helpers.PrettyPrint(deviceStore.PushName))
	wLog.Infof("Platform: %s", helpers.PrettyPrint(deviceStore.Platform))
	wLog.Infof("Account: %s", helpers.PrettyPrint(deviceStore.Account))

	cli = whatsmeow.NewClient(deviceStore, wLog)
	return cli
}

func ConnectClient(client *whatsmeow.Client) {
	if client.Store.ID == nil {
		// No ID stored, new login, show a qr code
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			wLog.Errorf("Error connecting: %v", err)
			os.Exit(1)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				wLog.Infof("Login event: %s", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := client.Connect()
		if err != nil {
			wLog.Errorf("Error connecting: %v", err)
			os.Exit(1)
		}
	}
}

func HandleEvent(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		wLog.Infof("Marked self as available")
		// if len(cli.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
		// 	err := cli.SendPresence(types.PresenceAvailable)
		// 	if err != nil {
		// 		wLog.Warnf("Failed to send available presence: %v", err)
		// 	} else {
		// 		wLog.Infof("Marked self as available")
		// 	}
		// }
	case *events.Connected, *events.PushNameSetting:
		if len(cli.Store.PushName) == 0 {
			return
		}
		wLog.Infof("Marked self as available")
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		// err := cli.SendPresence(types.PresenceAvailable)
		// if err != nil {
		// 	wLog.Warnf("Failed to send available presence: %v", err)
		// } else {
		// 	wLog.Infof("Marked self as available")
		// }
	case *events.StreamReplaced:
		os.Exit(0)
	case *events.Message:
		wLog.Infof("TYPE MESSAGE:", helpers.PrettyPrint(evt))
		go HandleMessage(evt)
	case *events.Receipt:
		if evt.Type == types.ReceiptTypeRead || evt.Type == types.ReceiptTypeReadSelf {
			wLog.Infof("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp)
		} else if evt.Type == types.ReceiptTypeDelivered {
			wLog.Infof("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp)
		}
	case *events.Presence:
		if evt.Unavailable {
			if evt.LastSeen.IsZero() {
				wLog.Infof("%s is now offline", evt.From)
			} else {
				wLog.Infof("%s is now offline (last seen: %s)", evt.From, evt.LastSeen)
			}
		} else {
			wLog.Infof("%s is now online", evt.From)
		}
	// case *events.HistorySync:
	// 	id := atomic.AddInt32(&historySyncID, 1)
	// 	fileName := fmt.Sprintf("history-%d-%d.json", startupTime, id)
	// 	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	// 	if err != nil {
	// 		log.Errorf("Failed to open file to write history sync: %v", err)
	// 		return
	// 	}
	// 	enc := json.NewEncoder(file)
	// 	enc.SetIndent("", "  ")
	// 	err = enc.Encode(evt.Data)
	// 	if err != nil {
	// 		log.Errorf("Failed to write history sync: %v", err)
	// 		return
	// 	}
	// 	log.Infof("Wrote history sync to %s", fileName)
	// 	_ = file.Close()
	case *events.AppState:
		wLog.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)
	case *events.KeepAliveTimeout:
		wLog.Debugf("Keepalive timeout event: %+v", evt)
	case *events.KeepAliveRestored:
		wLog.Debugf("Keepalive restored")
	case *events.Blocklist:
		wLog.Infof("Blocklist event: %+v", evt)
	default:
		wLog.Infof("Unknown event: %T\n", evt)
		wLog.Infof("EVENT: %s\n", helpers.PrettyPrint(evt))
	}
}

func HandleMessage(evt *events.Message) {
	// log.Println(helpers.PrettyPrint(evt))
	// msg := evt.Message.GetConversation()
	// log.Println(msg)
	if evt.Info.IsFromMe {
		return
	}

	repo.storeUser(evt.Info.Sender, evt.Info.PushName)

	mediaDir := "data/media"
	isStatusBroadcast := false

	metaParts := []string{fmt.Sprintf("pushname: %s", evt.Info.PushName), fmt.Sprintf("timestamp: %s", evt.Info.Timestamp)}
	if evt.Info.Type != "" {
		metaParts = append(metaParts, fmt.Sprintf("type: %s", evt.Info.Type))
	}
	if evt.Info.Category != "" {
		metaParts = append(metaParts, fmt.Sprintf("category: %s", evt.Info.Category))
	}
	if evt.IsViewOnce {
		metaParts = append(metaParts, "view once")
	}
	if evt.IsViewOnce {
		metaParts = append(metaParts, "ephemeral")
	}
	if evt.IsViewOnceV2 {
		metaParts = append(metaParts, "ephemeral (v2)")
	}
	if evt.IsDocumentWithCaption {
		metaParts = append(metaParts, "document with caption")
	}
	if evt.IsEdit {
		metaParts = append(metaParts, "edit")
	}
	if evt.Info.Chat.String() == "status@broadcast" {
		mediaDir = "data/media/_broadcast"
		isStatusBroadcast = true
	}

	if !isStatusBroadcast {
		repo.storeConversation(evt)
	}

	wLog.Infof("Received message %s from %s (%s): %+v", evt.Info.ID, evt.Info.SourceString(), strings.Join(metaParts, ", "), evt.Message)

	if isStatusBroadcast {
		wLog.Infof("Status broadcast received: %v", evt.Info.Chat)
	}

	if evt.Message.GetPollUpdateMessage() != nil {
		decrypted, err := cli.DecryptPollVote(evt)
		if err != nil {
			wLog.Errorf("Failed to decrypt vote: %v", err)
		} else {
			wLog.Infof("Selected options in decrypted vote:")
			for _, option := range decrypted.SelectedOptions {
				wLog.Infof("- %X", option)
			}
		}
	} else if evt.Message.GetEncReactionMessage() != nil {
		decrypted, err := cli.DecryptReaction(evt)
		if err != nil {
			wLog.Errorf("Failed to decrypt encrypted reaction: %v", err)
		} else {
			wLog.Infof("Decrypted reaction: %+v", decrypted)
		}
	}

	img := evt.Message.GetImageMessage()
	if img != nil {
		data, err := cli.Download(img)
		if err != nil {
			wLog.Errorf("Failed to download image: %v", err)
			return
		}

		savedFile, err := storeMedia(mediaDir, data, evt, evt.Info.PushName)
		if err != nil {
			wLog.Errorf("Failed to save video: %v", err)
			return
		}
		wLog.Infof("Saved video in message to %s", savedFile)
	}

	vid := evt.Message.GetVideoMessage()
	if vid != nil {
		data, err := cli.Download(vid)
		if err != nil {
			wLog.Errorf("Failed to download video: %v", err)
			return
		}

		savedFile, err := storeMedia(mediaDir, data, evt, evt.Info.PushName)
		if err != nil {
			wLog.Errorf("Failed to save video: %v", err)
			return
		}

		wLog.Infof("Saved video in message to %s", savedFile)
	}
}

func isStatusUpdate(evt *events.Message) bool {
	return evt.Info.Chat.String() == "status@broadcast"
}

type mediaInfo struct {
	caption   string
	mediaType string
	filesize  uint64
	height    uint32
	width     uint32
	duration  uint32
	fileLoc   string
}

func storeMedia(path string, data []byte, evt *events.Message, username string) (string, error) {
	filename := evt.Info.ID
	fileInfo := mediaInfo{}
	if evt.Info.MediaType == "image" {
		ext, _ := mime.ExtensionsByType(evt.Message.ImageMessage.GetMimetype())
		filename = fmt.Sprintf("%s%s", evt.Info.ID, ext[0])
		fileInfo.caption = evt.Message.ImageMessage.GetCaption()
		fileInfo.mediaType = evt.Message.ImageMessage.GetMimetype()
		fileInfo.filesize = evt.Message.ImageMessage.GetFileLength()
		fileInfo.height = evt.Message.ImageMessage.GetHeight()
		fileInfo.width = evt.Message.ImageMessage.GetWidth()

	}
	if evt.Info.MediaType == "video" {
		ext, _ := mime.ExtensionsByType(evt.Message.VideoMessage.GetMimetype())
		filename = fmt.Sprintf("%s%s", evt.Info.ID, ext[0])
		fileInfo.caption = evt.Message.VideoMessage.GetCaption()
		fileInfo.mediaType = evt.Message.VideoMessage.GetMimetype()
		fileInfo.filesize = evt.Message.VideoMessage.GetFileLength()
		fileInfo.height = evt.Message.VideoMessage.GetHeight()
		fileInfo.width = evt.Message.VideoMessage.GetWidth()
		fileInfo.duration = evt.Message.VideoMessage.GetSeconds()
	}

	userSlug := slug.Make(username)
	mediaDest := fmt.Sprintf("%s/%s_%s", path, evt.Info.Sender.User, userSlug)
	err := os.MkdirAll(mediaDest, 0755)
	if err != nil {
		return "", err
	}

	fileInfo.fileLoc = mediaDest + "/" + filename

	err = os.WriteFile(fileInfo.fileLoc, data, 0644)
	if err != nil {
		return "", err
	}

	if isStatusUpdate(evt) {
		err = repo.recordStatusUpdates(evt, fileInfo)
		if err != nil {
			wLog.Errorf("Failed to record status update: %v", err)
		}
	}

	return fileInfo.fileLoc, nil
}

func (r *waRepo) recordStatusUpdates(msg *events.Message, mediaInfo mediaInfo) error {
	sql := `INSERT INTO tbl_status_updates
	(msg_id, account, sender_phone, sender_jid, sender_name, caption, media_type, mimetype, filesize, height, width, file_loc, msg_date)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(sql,
		msg.Info.ID,
		r.account.ID,
		msg.Info.Sender.User,
		msg.Info.Sender.String(),
		msg.Info.PushName,
		mediaInfo.caption,
		mediaInfo.mediaType,
		mediaInfo.mediaType,
		mediaInfo.filesize,
		mediaInfo.height,
		mediaInfo.width,
		mediaInfo.fileLoc,
		msg.Info.Timestamp,
	)

	if err != nil {
		wLog.Errorf("Failed to store status update: %v", err)
	}
	return nil
}

func (r *waRepo) storeUser(jid types.JID, name string) error {
	sql := `INSERT INTO tbl_users (phone, jid, server, name) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = ?`
	_, err := r.db.Exec(sql, jid.User, jid.String(), jid.Server, name, name)
	if err != nil {
		wLog.Errorf("Failed to store user: %v", err)
	}
	return err
}

func (r *waRepo) storeConversation(msg *events.Message) error {
	sql := `INSERT INTO tbl_chats
	(room, msg_id, account, sender, name, is_group, msg_type, media_type, category, message, msg_date)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	isGroup := 0
	if msg.Info.IsGroup {
		isGroup = 1
	}

	_, err := r.db.Exec(sql,
		msg.Info.Chat.String(),
		msg.Info.ID,
		r.account.ID,
		msg.Info.Sender,
		msg.Info.PushName,
		isGroup,
		msg.Info.Type,
		msg.Info.MediaType,
		msg.Info.Category,
		msg.Message.GetConversation(),
		msg.Info.Timestamp)
	if err != nil {
		wLog.Errorf("Failed to store conversation: %v", err)
	}
	return err
}
