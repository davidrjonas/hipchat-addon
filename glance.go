package addon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// https://www.hipchat.com/docs/apiv2/glances
type Glance struct {
	Icon     *Image     `json:"icon"`
	Key      string     `json:"key"`
	Name     *I18nValue `json:"name"`
	QueryUrl string     `json:"queryUrl,omitempty"`
	Target   string     `json:"target,omitempty"`
	Weight   int16      `json:"weight,omitempty"`

	DataCallback GlanceDataCallbackFunc `json:"-"`
}

type GlanceDataCallbackFunc func(a *HipchatAddon, installation *Installation, g *Glance) *GlanceData

type GlanceUpdates struct {
	Updates []GlanceUpdate `json:"glance"`
}

type GlanceUpdate struct {
	Content *GlanceData `json:"content"`
	Key     string      `json:"key"`
}

type GlanceData struct {
	Label    *GlanceLabel    `json:"label,omitempty"`
	Status   *GlanceStatus   `json:"status,omitempty"`
	Metadata *GlanceMetadata `json:"metadata,omitempty"`
}

type GlanceLabel struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type GlanceStatus struct {
	Type  string             `json:"type"`
	Value *GlanceStatusValue `json:"value"`
}

type GlanceStatusValue struct {
	Label string `json:"label,omitempty"`
	Type  string `json:"type,omitempty"`
	Url   string `json:"url,omitempty"`
	Url2x string `json:"url@2x,omitempty"`
}

type GlanceMetadata struct {
	CustomData map[string]string `json:"customData"`
}

func (a *HipchatAddon) newGlanceHandler(glance *Glance) http.HandlerFunc {
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

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		data, _ := json.Marshal(glance.DataCallback(a, installation, glance))
		w.Write(data)
	}
}

func (a *HipchatAddon) UpdateGlances(updates *GlanceUpdates) {
	for _, installation := range a.installations.GetAll() {
		if err := a.UpdateGlanceData(installation, updates); err != nil {
			a.logger.Error(err)
		}
	}
}

func (a *HipchatAddon) UpdateGlanceData(installation *Installation, updates *GlanceUpdates) error {

	// FIXME: This is room specific. Use the installation to figure out what kind of url we need.
	// See https://developer.atlassian.com/hipchat/guide/glances#Glances-UpdatingtheGlancedata
	roomId := installation.RoomId.String()

	if roomId == "" {
		return fmt.Errorf("no room id for installation")
	}

	url := installation.ApiUrl + "addon/ui/room/" + roomId

	payload, err := json.Marshal(updates)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	resp, err := a.postWithToken(installation, url, bytes.NewReader(payload))

	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}
