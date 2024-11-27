package notification

import (
	"time"

	"wraith.me/message_server/pkg/util"
)

var (
	//The default TTL of a notification (default: 1 week).
	DefaultNotifTTL = "1w"
)

// Represents a single notification message that is sent to a user when certain actions occur.
type Notification struct {
	//The ID of this notification.
	ID util.UUID `json:"id"`

	//The ID of the user to whom this notification belongs.
	Recipient util.UUID `json:"recipient"`

	//The body of the notification.
	Content string `json:"content"`

	//The type of notification this is.
	Type Type `json:"type"`

	//Extra information needed for the notification to function depending on the `Type`.
	Context string `json:"context"`

	//Whether the notification was read by the user.
	Read bool `json:"read"`

	//The time at which the notification will expire and be auto-purged.
	Expires time.Time `json:"expires"`
}

// Gets an expiration time from the current time.
func resolveExpiryTime(now time.Time) time.Time {
	delta, _ := time.ParseDuration(DefaultNotifTTL)
	return now.Add(delta)
}
