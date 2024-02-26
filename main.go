package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wabot/internal/helpers"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var WaClient *whatsmeow.Client

func main() {
	WaClient := CreateClient()
	ConnectClient(WaClient)

	WaClient.AddEventHandler(HandleEvent)

	WaClient.Connect()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	WaClient.Disconnect()
}

func CreateClient() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "INFO", true)
	sql, err := sqlstore.New("sqlite3", "file:data/accounts.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalln(err)
	}

	deviceStore, err := sql.GetFirstDevice()
	if err != nil {
		log.Fatalln(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)

	client := whatsmeow.NewClient(deviceStore, clientLog)
	return client
}

func ConnectClient(client *whatsmeow.Client) {
	if client.Store.ID == nil {
		// No ID stored, new login, show a qr code
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			log.Fatalln(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				log.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := client.Connect()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func HandleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go HandleMessage(v)
	default:
		log.Printf("Unknown event: %T\n", v)
		log.Println(helpers.PrettyPrint(v))
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
func HandleMessage(messageEvent *events.Message) {
	log.Println(helpers.PrettyPrint(messageEvent))
	msg := messageEvent.Message.GetConversation()
	log.Println(msg)
	if messageEvent.Info.IsFromMe {
		return
	}
}
