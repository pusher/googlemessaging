package googlemessaging

type FcmMessageBody struct {
	Message *FcmHttpMessage `json:"message"`
}

type FcmHttpMessage struct {
	Topic        string            `json:"topic,omitempty"`
	Token        string            `json:"token,omitempty"`
	Notification *FcmNotification  `json:"notification,omitempty"`
	Data         FcmData           `json:"data,omitempty"`
	Android      *FcmAndroidConfig `json:"android"`
}

type FcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image"`
}

type FcmAndroidConfig struct {
	Notification *FcmAndroidNotification `json:"notification"`
}

type FcmAndroidNotification struct {
	Icon              string   `json:"icon,omitempty"`
	Sound             string   `json:"sound,omitempty"`
	Badge             string   `json:"badge,omitempty"`
	Tag               string   `json:"tag,omitempty"`
	Color             string   `json:"color,omitempty"`
	ClickAction       string   `json:"click_action,omitempty"`
	ChannelId         string   `json:"channel_id,omitempty"`
	BodyLocKey        string   `json:"body_loc_key,omitempty"`
	BodyLocArgs       []string `json:"body_loc_args,omitempty"`
	TitleLocArgs      []string `json:"title_loc_args,omitempty"`
	TitleLocKey       string   `json:"title_loc_key,omitempty"`
	Ticker            string   `json:"ticker,omitempty"`
	Sticky            string   `json:"sticky,omitempty"`
	LocalOnly         bool     `json:"local_only,omitempty"`
	NotificationCount int      `json:"notification_count,omitempty"`
}

type FcmData map[string]interface{}

type FcmSendHttpResponse struct {
	Status        int               `json:"-"`
	Name          string            `json:"name"`
	AndroidConfig *FcmAndroidConfig `json:"android"`
	Token         string            `json:"token"`
	Topic         string            `json:"topic"`
}
