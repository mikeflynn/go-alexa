package request

import (
	"errors"
	"time"
)

func (r *Request) TimestampValid() bool {
	reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", r.Body.Timestamp)
	if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
		return true
	}

	return false
}

func (r *Request) ApplicationIDValid(myAppIDs []string) bool {
	// TODO: Probably a better way to do r
	for _, str := range myAppIDs {
		if str == r.Session.Application.ApplicationID {
			return true
		}
	}
	return false
}

func (r *Request) GetSessionID() string {
	return r.Session.SessionID
}

func (r *Request) GetUserID() string {
	return r.Session.User.UserID
}

func (r *Request) GetRequestType() string {
	return r.Body.Type
}

func (r *Request) GetIntentName() string {
	if r.GetRequestType() == "IntentRequest" {
		return r.Body.Intent.Name
	}

	return r.GetRequestType()
}

func (r *Request) GetSlotValue(slotName string) (string, error) {
	slot, ok := r.Body.Intent.Slots[slotName]
	if !ok {
		return "", errors.New("Slot name not found.")
	}

	return slot.Value, nil
}

func (r *Request) AllSlots() map[string]Slot {
	return r.Body.Intent.Slots
}
