package addon

// https://www.hipchat.com/docs/apiv2/webpanels
type WebPanel struct {
	Icon     *Image     `json:"icon,omitempty"`
	Key      string     `json:"key"`
	Location string     `json:"location"`
	Name     *I18nValue `json:"name"`
	Url      string     `json:"url"`
	Weight   int16      `json:"weight,omitempty"`
}

// TODO: implement webpanels
