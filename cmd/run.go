package cmd

import (
	"context"
	"fmt"
	"mime"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"wabot/internal/helpers"

	"github.com/gosimple/slug"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
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

var logLevel = "INFO"

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		wLog = waLog.Stdout("Main", logLevel, true)
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
	dbLog := waLog.Stdout("Database", logLevel, true)
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

	wLog.Infof("Device: %s", helpers.PrettyPrint(deviceStore.ID))
	wLog.Infof("Pushname: %s", helpers.PrettyPrint(deviceStore.PushName))
	wLog.Infof("Platform: %s", helpers.PrettyPrint(deviceStore.Platform))

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

/*
- Type:

	{
	"Info": {
	"Chat": "6281805931588-1492190836@g.us",
	"Sender": "15185546427:5@s.whatsapp.net",
	"IsFromMe": true,
	"IsGroup": true,
	"BroadcastListOwner": "",
	"ID": "3EB022B467F7CBCD31345F",
	"ServerID": 0,
	"Type": "text",
	"PushName": "Jeflon Zuckergates",
	"Timestamp": "2024-02-26T09:00:46+07:00",
	"Category": "",
	"Multicast": false,
	"MediaType": "",
	"Edit": "",
	"VerifiedName": null,
	"DeviceSentMeta": null
	},
	"Message": {
	"conversation": "KIONG HI GAEZZZZZZZZZZZZZZ"
	},
	"IsEphemeral": false,
	"IsViewOnce": false,
	"IsViewOnceV2": false,
	"IsDocumentWithCaption": false,
	"IsEdit": false,
	"SourceWebMsg": null,
	"UnavailableRequestID": "",
	"RetryCount": 0,
	"NewsletterMeta": null,
	"RawMessage": {
	"conversation": "KIONG HI GAEZZZZZZZZZZZZZZ"
	}
	}
*/
func HandleMessage(evt *events.Message) {
	// log.Println(helpers.PrettyPrint(evt))
	// msg := evt.Message.GetConversation()
	// log.Println(msg)
	// if evt.Info.IsFromMe {
	// 	return
	// }

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
		exts, _ := mime.ExtensionsByType(img.GetMimetype())
		path := fmt.Sprintf("%s%s", evt.Info.ID, exts[0])
		userSlug := slug.Make(evt.Info.PushName)
		mediaDest := fmt.Sprintf("%s/%s", mediaDir, userSlug)
		err = os.MkdirAll(mediaDest, 0755)
		if err != nil {
			wLog.Errorf("Failed to create media directory: %v", err)
			return
		}
		err = os.WriteFile(mediaDest+"/"+path, data, 0755)
		if err != nil {
			wLog.Errorf("Failed to save image: %v", err)
			return
		}
		wLog.Infof("Saved image in message to %s", path)
	}

	vid := evt.Message.GetVideoMessage()
	if vid != nil {
		data, err := cli.Download(vid)
		if err != nil {
			wLog.Errorf("Failed to download video: %v", err)
			return
		}
		exts, _ := mime.ExtensionsByType(vid.GetMimetype())
		path := fmt.Sprintf("%s%s", evt.Info.ID, exts[0])
		userSlug := slug.Make(evt.Info.PushName)
		mediaDest := fmt.Sprintf("%s/%s", mediaDir, userSlug)
		err = os.MkdirAll(mediaDest, 0644)
		if err != nil {
			wLog.Errorf("Failed to create media directory: %v", err)
			return
		}
		err = os.WriteFile(mediaDest+"/"+path, data, 0644)
		if err != nil {
			wLog.Errorf("Failed to save video: %v", err)
			return
		}
		wLog.Infof("Saved video in message to %s", path)
	}
}

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			wLog.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			wLog.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}
