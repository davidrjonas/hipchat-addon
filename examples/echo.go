package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/davidrjonas/hipchat-addon"
)

var urlBase string
var host string
var port int
var stateFilename string

func init() {
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)

	flag.StringVar(&urlBase, "url", "", "The base URL, defaults to the host and port")
	flag.StringVar(&host, "host", "127.0.0.1", "The IP address on which to listen")
	flag.IntVar(&port, "port", 3000, "The port on which to listen")
	flag.StringVar(&stateFilename, "state-file", "state", "The file to read/store state information")
}

func onEchoWebHook(a *addon.HipchatAddon, installation *addon.Installation, webhook *addon.WebHook, event map[string]interface{}) error {

	logrus.Info("Received webhook callback")

	var msg string

	// Encode the entire event in json and send it back to the chat room. This
	// isn't really "echoing" it is more like var dumping.
	if data, err := json.Marshal(event); err == nil {
		msg = string(data)
	} else {
		msg = fmt.Sprintf("Error: %v", err)
	}

	if err := a.SendNotification(installation, &addon.Notification{Message: msg}); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func url(resource string) string {
	return urlBase + resource
}

func main() {
	flag.Parse()

	if urlBase == "" {
		urlBase = fmt.Sprintf("http://%s:%d", host, port)
	}

	a := addon.NewWithStateFile(
		&addon.CapabilitiesDescriptor{
			Name:        "Echo HipChat AddOn",
			Description: "Example HipChat Addon",
			Key:         "echo-addon",
			Vendor: &addon.Vendor{
				Name: "davidrjonas",
				Url:  "https://github.com/davidrjonas/hipchat-addon",
			},
			Links: &addon.Links{
				Homepage: "https://github.com/davidrjonas/hipchat-addon",
				Self:     url("/capabilities.json"),
			},
			Capabilities: &addon.Capabilities{
				HipchatApiConsumer: &addon.HipchatApiConsumer{
					Scopes: []string{"send_notification"},
				},
				Installable: &addon.Installable{
					AllowGlobal: false,
					AllowRoom:   true,
					CallbackUrl: url("/install"),
				},
				WebHook: []*addon.WebHook{&addon.WebHook{
					Event:          "room_message",
					Pattern:        "^/echo\b.+",
					Authentication: "jwt",
					Name:           "Echo",
					Url:            url("/webhook/0"),

					Callback: onEchoWebHook,
				}},
			},
		},
		stateFilename,
		addon.Logger(logrus.StandardLogger()),
	)

	logrus.Infof("Saving state to file '%s'", stateFilename)
	logrus.Infof("Starting server on %s:%d for url %s", host, port, urlBase)

	a.Serve(fmt.Sprintf("%s:%d", host, port))
}
