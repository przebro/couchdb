package context

import (
	"net/http"
)

//AuthType - Authentication type used to establish connection
type AuthType string

const (
	//None - Each request will contain login in password in an url
	None AuthType = "NONE"
	//Basic - Authenticate with standard RFC2617 mechanism
	Basic AuthType = "BASIC"
	/*Cookie - Authenticate the user by sending the POST request to /_session endpoint with a name and password in the body of the request.
	If successful then the returned cookie will be attached to the context.
	*/
	Cookie AuthType = "COOKIE"
	/*JwtToken - As refrecne says: "Enables CouchDB to use externally generated tokens instead of defining users or roles..."
		Make sure that this kind of authentication is enabled in the CouchDB instance.
		[chttpd]
	authentication_handlers = {chttpd_auth, cookie_authentication_handler}, {chttpd_auth, jwt_authentication_handler}, {chttpd_auth, default_authentication_handler}
	*/
	JwtToken AuthType = "JWT"
)

//CouchContext - Holds connection
type CouchContext struct {
	BaseAddr       string
	isSecure       bool
	Authentication AuthType
	AuthData       string
	Client         *http.Client
}
