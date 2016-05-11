package addon

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Installation struct {
	OauthId         string      `json:"oauthId"`
	CapabilitiesUrl string      `json:"capabilitiesUrl"`
	RoomId          json.Number `json:"roomId"`
	GroupId         json.Number `json:"groupId"`
	OauthSecret     string      `json:"oauthSecret"`
	TokenUrl        string      `json:"tokenUrl"`
	ApiUrl          string      `json:"apiUrl"`

	token *AccessToken
}

func (i *Installation) GetAccessToken(client HttpDoer) (*AccessToken, error) {

	if i.token == nil || i.token.Valid() {

		newToken, err := i.getFreshAccessToken(client)

		if err != nil {
			return nil, err
		}

		i.token = newToken
	}

	return i.token, nil
}

func (i *Installation) getFreshAccessToken(client HttpDoer) (*AccessToken, error) {

	params := url.Values{"grant_type": {"client_credentials"}, "scope": {"send_notification"}}

	req, err := http.NewRequest("POST", i.TokenUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(i.OauthId, i.OauthSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return NewAccessTokenFromJson(resp.Body)
}
