package addon

// https://www.hipchat.com/docs/apiv2/actions
type Action struct {
	Key      string     `json:"key"`
	Location string     `json:"location"`
	Name     *I18nValue `json:"name"`
	Target   string     `json:"target"`
	Weight   int16      `json:"weight,omitempty"`
}

// TODO: implement actions
