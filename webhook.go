package addon

import (
	"encoding/json"
	"net/http"
)

// https://www.hipchat.com/docs/apiv2/webhooks
// Events: room_archived, room_created, room_deleted, room_enter, room_exit,
// room_file_upload, room_message, room_notification, room_topic_change,
// room_unarchived
type WebHook struct {
	Authentication string `json:"authentication,omitempty"`
	Event          string `json:"event"`
	Key            string `json:"key,omitempty"`
	Name           string `json:"name,omitempty"`
	Pattern        string `json:"pattern,omitempty"`
	Url            string `json:"url"`

	Callback WebHookCallback `json:"-"`
}

type WebHookCallback func(a *HipchatAddon, installation *Installation, webhook *WebHook, event map[string]interface{}) error

func (a *HipchatAddon) newWebHookHandler(webhook *WebHook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We have a valid jwt token verified by upstream middleware. Since
		// http doesn't support contexts yet and I don't want to rely on
		// gorilla lets just reparse the jwt for now.
		token, err := jwtParseFromHipChatRequest(r, a.jwtKeyLookup)

		if !token.Valid {
			// Someone didn't do their middleware.
			panic(err)
		}

		oauthId := token.Claims["iss"].(string)

		installation := a.installations.Get(oauthId)

		if installation == nil {
			http.Error(w, "404 No installation found for iss", http.StatusNotFound)
			return
		}

		data := map[string]interface{}{}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "400 Unable to parse json", http.StatusBadRequest)
			return
		}

		if err := webhook.Callback(a, installation, webhook, data); err != nil {
			http.Error(w, "500 Internal Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
	}
}
