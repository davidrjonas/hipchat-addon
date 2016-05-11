package addon

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Notification struct {
	MessageFormat string           `json:"message_format,omitempty"` // text, html
	Message       string           `json:"message"`
	Notify        bool             `json:"notify,omitempty"`
	Color         string           `json:"color,omitempty"` // yellow, green, red, purple, gray, random
	Card          *json.RawMessage `json:"card,omitempty"`
}

func (a *HipchatAddon) SendNotification(installation *Installation, notification *Notification) error {
	roomId := installation.RoomId.String()

	if roomId == "" {
		return errors.New("no room id for this installation")
	}

	notificationUrl := installation.ApiUrl + "room/" + roomId + "/notification"

	data, _ := json.Marshal(notification)

	resp, err := a.postWithToken(installation, notificationUrl, bytes.NewReader(data))

	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}
