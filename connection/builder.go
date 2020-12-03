package connection

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/przebro/couchdb/context"
)

type ConnectionScheme string

const (
	httpScheme  ConnectionScheme = "http"
	httpsScheme ConnectionScheme = "https"
)

type builder struct {
	cert      bool
	host      string
	port      int
	auth      context.AuthType
	authData  string
	sesToken  string
	jwtToken  string
	usernmame string
	password  string
}

type ConnectionBuilder interface {
	WithAddress(address string, port int) ConnectionBuilder
	WithCertificate() ConnectionBuilder
	WithToken(token string) ConnectionBuilder
	WithAuthentication(atype context.AuthType, username, password string) ConnectionBuilder
	Build() (*Connection, error)
}

func NewBuilder() ConnectionBuilder {

	b := &builder{auth: context.None}
	return b
}

func (b *builder) WithAddress(host string, port int) ConnectionBuilder {

	b.host = host
	b.port = port
	return b
}
func (b *builder) WithAuthentication(atype context.AuthType, username, password string) ConnectionBuilder {

	b.auth = atype

	if atype == context.Basic {
		data := fmt.Sprintf("%s:%s", username, password)
		b.authData = base64.StdEncoding.EncodeToString([]byte(data))
	} else {
		b.authData = fmt.Sprintf("%s:%s", username, password)
	}

	return b
}
func (b *builder) WithCertificate() ConnectionBuilder {

	return b
}
func (b *builder) WithToken(token string) ConnectionBuilder {
	b.jwtToken = token
	b.auth = context.JwtToken
	return b
}

func (b *builder) Build() (*Connection, error) {

	scheme := func() ConnectionScheme {
		if b.cert {
			return httpsScheme
		}
		return httpScheme

	}()

	addr := func() string {
		if b.auth == context.None {
			return fmt.Sprintf(`%s://%s@%s:%d`, scheme, b.authData, b.host, b.port)
		}
		return fmt.Sprintf(`%s://%s:%d`, scheme, b.host, b.port)

	}()

	if b.auth == context.Cookie && b.authData == "" {
		return nil, errors.New("no authentication data")
	}

	authData := ""
	if b.auth != context.Cookie {
		authData = b.authData
	}

	ctx := &context.CouchContext{Client: &http.Client{},
		BaseAddr:       addr,
		Authentication: b.auth,
		AuthData:       authData,
	}
	conn := &Connection{ctx: ctx}

	var err error

	if b.auth == context.Cookie {

		userpass := strings.SplitN(b.authData, ":", 2)
		_, err = conn.Session(userpass[0], userpass[1])
	} else {
		_, err = conn.Up()
	}

	return conn, err
}
