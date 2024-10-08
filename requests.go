package googlemessaging

import (
	"errors"
)

var (
	ErrInvalidToken                = errors.New("invalid token")
	ErrDeviceNotFound              = errors.New("device not found")
	ErrInvalidRequest              = errors.New("invalid request")
	ErrInalidFCMServiceAccountFile = errors.New("invalid fcm service account file")
)

type FcmMessageBody struct {
	Message *FcmHttpMessage `json:"message"`
}

type FcmHttpMessage struct {
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
	CollapseKey           string                  `json:"collapse_key"`
	Priority              string                  `json:"priority,omitempty"`
	TimeToLive            *string                 `json:"ttl,omitempty"`
	RestrictedPackageName string                  `json:"restricted_package_name,omitempty"`
	Notification          *FcmAndroidNotification `json:"notification"`
}

type FcmAndroidNotification struct {
	Icon              string   `json:"icon,omitempty"`
	Color             string   `json:"color,omitempty"`
	Sound             string   `json:"sound,omitempty"`
	Tag               string   `json:"tag,omitempty"`
	ClickAction       string   `json:"click_action,omitempty"`
	BodyLocKey        string   `json:"body_loc_key,omitempty"`
	BodyLocArgs       []string `json:"body_loc_args,omitempty"`
	TitleLocKey       string   `json:"title_loc_key,omitempty"`
	TitleLocArgs      []string `json:"title_loc_args,omitempty"`
	Ticker            string   `json:"ticker,omitempty"`
	Sticky            bool     `json:"sticky,omitempty"`
	LocalOnly         bool     `json:"local_only,omitempty"`
	NotificationCount int      `json:"notification_count,omitempty"`
	ChannelId         string   `json:"channel_id,omitempty"`
}

type FcmData map[string]string

type FcmSendHttpResponse struct {
	Status int    `json:"-"`
	Name   string `json:"name"`
}

type InstanceInformationResponse struct {
	AuthorizedEntity string `json:"authorizedEntity"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
