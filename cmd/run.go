package cmd

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"os"
	"os/signal"
	"syscall"

	// "github.com/ditatompel/wa-status-archiver/internal/botcmd"
	"github.com/ditatompel/wa-status-archiver/internal/database"
	"github.com/ditatompel/wa-status-archiver/internal/helpers"

	"github.com/gosimple/slug"
	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
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

var wa *waRepo

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Long:  `Run the bot by listening to WA websocket events`,
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
	// sql, err := sqlstore.New("sqlite3", "file:data/db/accounts.db?_foreign_keys=on", dbLog)
	sql := sqlstore.NewWithDB(database.GetDB().DB, "postgres", dbLog)
	err := sql.Upgrade()
	if err != nil {
		wLog.Errorf("Failed to upgrade db: %v", err)
		os.Exit(1)
	}
	err = database.CreateSchema(database.GetDB())
	if err != nil {
		wLog.Errorf("Error creating app schema: %v", err)
		os.Exit(1)
	}

	deviceStore, err := sql.GetFirstDevice()
	if err != nil {
		wLog.Errorf("Error getting device: %v", err)
		os.Exit(1)
	}

	wa = newWaRepo(database.GetDB(), deviceStore)

	wLog.Infof("Device: %s", deviceStore.ID)
	wLog.Infof("Pushname: %s", deviceStore.PushName)
	wLog.Infof("Platform: %s", deviceStore.Platform)

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

// You can implement any event handler here.
// see https://github.com/tulir/whatsmeow/blob/main/mdtest/main.go
func HandleEvent(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(cli.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := cli.SendPresence(types.PresenceAvailable)
			if err != nil {
				wLog.Warnf("Failed to send available presence: %v", err)
			} else {
				wLog.Infof("Marked self as available")
			}
		}
	// case *events.Connected:
	// 	if len(cli.Store.PushName) == 0 {
	// 		return
	// 	}
	// 	err := cli.SendPresence(types.PresenceAvailable)
	// 	if err != nil {
	// 		wLog.Warnf("Failed to send available presence: %v", err)
	// 	} else {
	// 		wLog.Infof("Marked self as available")
	// 	}
	case *events.PushNameSetting:
		if len(cli.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := cli.SendPresence(types.PresenceAvailable)
		if err != nil {
			wLog.Warnf("Failed to send available presence: %v", err)
		} else {
			wLog.Infof("Marked self as available")
		}
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
	case *events.AppState:
		wLog.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)
	case *events.KeepAliveTimeout:
		wLog.Debugf("Keepalive timeout event: %+v", evt)
	case *events.KeepAliveRestored:
		wLog.Debugf("Keepalive restored")
	case *events.Blocklist:
		wLog.Infof("Blocklist event: %+v", evt)
	default:
		wLog.Infof("Untracked event: %T", evt)
		wLog.Debugf("EVENT: %s", helpers.PrettyPrint(evt))
	}
}

// this is message handler example
func HandleMessage(evt *events.Message) {
	mediaDir := "data/media"
	isStatusBroadcast := false

	if evt.Info.Chat.String() == "status@broadcast" {
		mediaDir = "data/media/_broadcast"
		isStatusBroadcast = true
	}

	if !isStatusBroadcast {
		wa.storeConversation(evt)

		// this example how the bot response to messages
		// this not enabled by default, you can uncomment it and implement your bot logic
		// if !evt.Info.IsFromMe {
		// 	botresp := botcmd.ParseCmd(evt.Message.GetConversation())
		// 	if botresp != "" {
		// 		if err := botSendMsg(evt, botresp); err != nil {
		// 			wLog.Errorf("Failed to send message: %v", err)
		// 		}
		// 	}
		// }
	} else {
		wLog.Infof("Status broadcast received: %v", evt.Info.Chat)
	}

	wLog.Infof("Received message %s from %s", evt.Info.ID, evt.Info.SourceString())
	wLog.Debugf("%s DATA: %+v", evt.Info.ID, evt)

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
	caption      string
	mediaType    string
	filesize     uint64
	height       uint32
	width        uint32
	duration     uint32
	fileLocation string
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
	if err := os.MkdirAll(mediaDest, 0755); err != nil {
		return "", err
	}

	fileInfo.fileLocation = mediaDest + "/" + filename

	if err := os.WriteFile(fileInfo.fileLocation, data, 0644); err != nil {
		return "", err
	}

	if isStatusUpdate(evt) {
		if err := wa.recordStatusUpdates(evt, fileInfo); err != nil {
			wLog.Errorf("Failed to record status update: %v", err)
		}
	}

	return fileInfo.fileLocation, nil
}

func (r *waRepo) recordStatusUpdates(msg *events.Message, mediaInfo mediaInfo) error {
	sql := `INSERT INTO tbl_statuses
	(message_id, our_jid, sender_jid, sender_name, caption, media_type, mimetype, filesize, height, width, file_location, msg_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Exec(sql,
		msg.Info.ID,
		r.account.ID,
		msg.Info.Sender.String(),
		msg.Info.PushName,
		mediaInfo.caption,
		msg.Info.MediaType,
		mediaInfo.mediaType,
		mediaInfo.filesize,
		mediaInfo.height,
		mediaInfo.width,
		mediaInfo.fileLocation,
		msg.Info.Timestamp,
	)
	if err != nil {
		wLog.Errorf("Failed to store status update: %v", err)
	}
	return nil
}

func (r *waRepo) storeConversation(msg *events.Message) error {
	sql := `INSERT INTO tbl_chats
	(message_id, room_id, our_jid, sender_jid, sender_name, is_group, is_from_me, msg_type, media_type, msg_conversation, category, msg_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Exec(sql,
		msg.Info.ID,
		msg.Info.Chat.String(),
		r.account.ID,
		msg.Info.Sender,
		msg.Info.PushName,
		msg.Info.IsGroup,
		msg.Info.IsFromMe,
		msg.Info.Type,
		msg.Info.MediaType,
		msg.Message.GetConversation(),
		msg.Info.Category,
		msg.Info.Timestamp,
	)
	if err != nil {
		wLog.Errorf("Failed to store conversation: %v", err)
	}
	return err
}

func botSendMsg(evt *events.Message, message string) error {
	roomId := evt.Info.Chat.ToNonAD()
	jid := evt.Info.Sender.String()
	recipient, ok := helpers.ParseJID(jid)
	if !ok {
		wLog.Errorf("Invalid recipient: %s", jid)
		return errors.New("invalid recipient")
	}
	msg := &waProto.Message{
		// Conversation: proto.String(message),
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(message),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:    proto.String(evt.Info.ID),
				Participant: proto.String(recipient.ADString()),
				QuotedMessage: &waProto.Message{
					Conversation: proto.String(evt.Message.GetConversation()),
				},
			},
		},
	}

	messageId := cli.GenerateMessageID()
	// resp, err := cli.SendMessage(context.Background(), recipient.ToNonAD(), msg)
	resp, err := cli.SendMessage(
		context.Background(),
		roomId,
		msg,
		whatsmeow.SendRequestExtra{
			ID: messageId,
		})
	if err != nil {
		return err
	}
	wLog.Infof("Sent message: %s", resp)

	return nil
}
