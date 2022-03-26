package utils

import "github.com/timworks/twilio-sdk-go/utils"

func IsNotFoundError(err error) bool {
	if twilioError, ok := err.(*utils.TwilioError); ok {
		return twilioError.IsNotFoundError()
	}
	return false
}
