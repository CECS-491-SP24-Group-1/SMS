package util

import (
	"fmt"
	"strings"
	"time"
)

//
//-- CLASS: CookieBuilder
//

/*
Utility class for building HTTP cookies. See the following MDN article
for full docs: https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies
This class may be superseded by Go's built-in `net/http/cookie.go`.
*/
type CookieBuilder struct {
	//The cookie's name.
	key string

	//The cookie's value.
	value string

	//The domain/subdomain for which the cookie is valid.
	domain string

	//The time at which the cookie should expire. Generally, `maxAge` should be preferred over this value.
	expiry string

	//Whether the cookie is restricted to HTTP(S) requests and inaccessible to client-side JS.
	httpOnly bool

	//THe maximum time (in seconds) that the cookie should be valid for.
	maxAge uint

	//The path at which the cookie is valid.
	path string

	//The value for the "SameSite" attribute, which controls access by 3p sites.
	sameSite string

	//Whether the cookie should be restricted to HTTPS requests. Requests to `localhost` are an exception.
	secure bool
}

//-- Constructors

// Creates a new CookieBuilder object.
func NewCookieBuilder(key, value string) *CookieBuilder {
	assertNonNullStr(key, "key")
	assertNonNullStr(value, "value")
	return &CookieBuilder{key: key, value: value}
}

//-- Methods

// Builds the cookie constructed by this class.
func (cb *CookieBuilder) Build() string {
	//Create a cookie string and write the key/value
	var cookie strings.Builder
	cookie.WriteString(fmt.Sprintf("%s=%s", cb.key, cb.value))

	//Add the attributes
	addIfNonNull(&cookie, "Domain", cb.domain)
	addIfNonNull(&cookie, "Expires", cb.expiry)
	if cb.httpOnly {
		addIfNonNull(&cookie, "", "HttpOnly", true)
	}
	if cb.maxAge > 0 {
		addIfNonNull(&cookie, "Max-Age", fmt.Sprintf("%d", cb.maxAge))
	}
	addIfNonNull(&cookie, "Path", cb.path)
	addIfNonNull(&cookie, "SameSite", cb.sameSite)
	if cb.secure {
		addIfNonNull(&cookie, "", "Secure", true)
	}

	//Return the full cookie string
	return cookie.String()
}

/*
Adds a `HttpOnly` attribute to the cookie. This restricts its
access and usage to HTTP(S) requests, denying access via JS. This
can help defend against XSS attacks.
*/
func (cb *CookieBuilder) SetHttpOnly() *CookieBuilder {
	cb.httpOnly = true
	return cb
}

/*
Sets the cookie as being "lax" with regards to the `SameSite`
attribute. This causes the cookie to be sent for requests from
the same site and for simple navigation links from other sites.
This is the default value for `SameSite`.
*/
func (cb *CookieBuilder) SetSameSiteLax() *CookieBuilder {
	cb.sameSite = "Lax"
	return cb
}

/*
Sets the cookie as being "lax" with regards to the `SameSite`
attribute. This causes the cookie to be sent for any request,
even if from other sites. This is a very insecure option, and
requires the cookie to be secure to work properly.
*/
func (cb *CookieBuilder) SetSameSiteNone() *CookieBuilder {
	cb.sameSite = "None"
	return cb
}

/*
Sets the cookie as being "strict" with regards to the `SameSite`
attribute. This causes the cookie to only be sent for requests
from the same site.
*/
func (cb *CookieBuilder) SetSameSiteStrict() *CookieBuilder {
	cb.sameSite = "Strict"
	return cb
}

/*
Adds a `Secure` attribute to the cookie. This restricts its
usage to secure protocols such as HTTPS.
*/
func (cb *CookieBuilder) SetSecure() *CookieBuilder {
	cb.secure = true
	return cb
}

/*
Adds a `Domain` attribute to the cookie. This restricts its
usage to the domain and subdomains in the specified value.
*/
func (cb *CookieBuilder) WithDomain(domain string) *CookieBuilder {
	assertNonNullStr(domain, "domain")
	cb.domain = domain
	return cb
}

/*
Adds an `Expires` attribute to the cookie. This defines a specific
date and time at which the cookie will be deleted. Cookies that don't
include this field or `Max-Age` will be deleted upon termination of
the "current session".
*/
func (cb *CookieBuilder) WithExpiry(expiration time.Time) *CookieBuilder {
	cb.expiry = expiration.UTC().Format(time.RFC1123)
	return cb
}

/*
Adds a `Max-Age` attribute to the cookie. This defines how long
(in seconds) the browser should retain the cookie before deleting
it. Cookies that don't include this field or `Expiry` will be
deleted upon termination of the "current session".
*/
func (cb *CookieBuilder) WithMaxAge(maxAge uint) *CookieBuilder {
	cb.maxAge = maxAge
	return cb
}

/*
Adds a `Path` attribute to the cookie. This restricts its
usage to the file path in the specified value.
*/
func (cb *CookieBuilder) WithPath(path string) *CookieBuilder {
	assertNonNullStr(path, "path")
	cb.path = path
	return cb
}

//-- Utilities

// Asserts that a variable is a non-null, non-empty string.
func assertNonNullStr(val, name string) {
	if val == "" {
		panic(fmt.Sprintf("'%s' must be a non-empty string", name))
	}
}

// Adds a keypair to a cookie string only if the value is non-null.
func addIfNonNull(cstr *strings.Builder, keyName, val string, noKey ...bool) {
	if val != "" {
		if len(noKey) > 0 && noKey[0] {
			cstr.WriteString("; " + val)
		} else {
			cstr.WriteString(fmt.Sprintf("; %s=%s", keyName, val))
		}
	}
}
