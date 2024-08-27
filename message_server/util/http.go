package util

import (
	"net/http"
	"strconv"
	"time"

	"wraith.me/message_server/consts"
)

// Gets a cookie value from a request by its key.
func StringFromCookie(r *http.Request, key string) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// Gets a UTL query param from a request by its key.
func StringFromQuery(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

/*
Calculates the offset from UTC based on the value of the `X-Timezone-Offset`
HTTP header.
*/
func Time2OffsetReq(tin time.Time, r *http.Request) time.Time {
	return Time2Offset(tin, TZOffsetFromReq(r))
}

// Gets the timezone offset from a request; returns 0 if not present.
func TZOffsetFromReq(r *http.Request) int {
	//Get the timezone from the HTTP headers
	off := r.Header.Get(consts.TIMEZONE_OFFSET_HEADER)
	ioff, err := strconv.Atoi(off)
	if err != nil {
		return 0
	}

	/*
		Ensure the offset is in the range +/- 720 since `Date.prototype.getTimezoneOffset()`
		returns the offset from UTC in minutes and 60 * 12 = 720. See the following webpage
		for more info: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/getTimezoneOffset
	*/
	if ioff > 720 {
		ioff = 720
	}
	if ioff < -720 {
		ioff = -720
	}

	//Return the truncated timezone offset
	return ioff
}
