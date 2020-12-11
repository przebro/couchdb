package connection

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/przebro/couchdb/response"

	"github.com/przebro/couchdb/client"
)

type connectionScheme string

const (
	httpScheme  connectionScheme = "http"
	httpsScheme connectionScheme = "https"
)

type builder struct {
	cert      bool
	host      string
	port      int
	auth      client.AuthType
	authData  string
	sesToken  string
	jwtToken  string
	usernmame string
	password  string

	skipVerify bool
	caPath     string
	keypath    string
	certpath   string
}

//ConnectionBuilder - Builds a connection with database
type ConnectionBuilder interface {
	WithAddress(address string, port int) ConnectionBuilder
	WithCertificate(rootca, clientkey, cert string, skipVerify bool) ConnectionBuilder
	WithToken(token string) ConnectionBuilder
	WithAuthentication(atype client.AuthType, username, password string) ConnectionBuilder
	Build(connect bool) (*Connection, error)
}

//NewBuilder - Creates a new ConnectionBuilder
func NewBuilder() ConnectionBuilder {

	b := &builder{auth: client.None}
	return b
}

func (b *builder) WithAddress(host string, port int) ConnectionBuilder {

	b.host = host
	b.port = port
	return b
}
func (b *builder) WithAuthentication(atype client.AuthType, username, password string) ConnectionBuilder {

	b.auth = atype

	if username == "" && password == "" {
		return b
	}

	if atype == client.Basic {
		data := fmt.Sprintf("%s:%s", username, password)
		b.authData = base64.StdEncoding.EncodeToString([]byte(data))
	} else {
		b.authData = fmt.Sprintf("%s:%s", username, password)
	}

	return b
}
func (b *builder) WithCertificate(rootca, clientkey, clinetCert string, skipVerify bool) ConnectionBuilder {

	b.cert = true
	b.caPath = rootca
	b.keypath = clientkey
	b.certpath = clinetCert
	b.skipVerify = skipVerify
	return b
}
func (b *builder) WithToken(token string) ConnectionBuilder {
	b.jwtToken = token
	b.auth = client.JwtToken
	return b
}

/*Build - Set up and build connections additionally if flag connect is set to true then invoke an authorization method
or simply call  up endpoint to check if connection is set up properly
*/
func (b *builder) Build(connect bool) (*Connection, error) {

	scheme := func() connectionScheme {
		if b.cert {
			return httpsScheme
		}
		return httpScheme

	}()

	addr := func() string {
		if b.auth == client.None {
			return fmt.Sprintf(`%s://%s@%s:%d`, scheme, b.authData, b.host, b.port)
		}
		return fmt.Sprintf(`%s://%s:%d`, scheme, b.host, b.port)

	}()

	if b.auth == client.Cookie && b.authData == "" {
		return nil, errors.New("no authentication data")
	}

	authData := ""
	if b.auth != client.Cookie {
		authData = b.authData
	}

	var tran *http.Transport = nil

	if b.cert {
		var er error
		tran, er = b.buildSecureTransport()
		if er != nil {
			return nil, er
		}
	} else {
		tran = http.DefaultTransport.(*http.Transport)
	}

	cli := &client.CouchClient{Client: &http.Client{Transport: tran},
		BaseAddr:       addr,
		Authentication: b.auth,
		AuthData:       authData,
	}
	conn := &Connection{cli: cli}

	var err error

	if connect {

		var res *response.CouchResult
		var err error
		if b.auth == client.Cookie {

			userpass := strings.SplitN(b.authData, ":", 2)
			res, err = conn.Session(context.TODO(), userpass[0], userpass[1])

		} else {
			res, err = conn.Up(context.TODO())
		}

		if err != nil {
			return conn, err
		}

		if res.Code >= 400 {
			return conn, errors.New(res.Status)
		}

	}

	return conn, err
}

//buildSecureTransport - Builds new Transport with provided certificates
func (b *builder) buildSecureTransport() (*http.Transport, error) {

	var certs []tls.Certificate = nil

	/*
		check if both client key and client certificate are provided if not,
		continue, however, a server may be configured to require a client's certificate and reject a connection without a client's certificate
	*/
	if b.certpath != "" && b.keypath != "" {

		cert, err := tls.LoadX509KeyPair(b.certpath, b.keypath)
		if err != nil {
			return nil, err
		}

		certs = []tls.Certificate{cert}
	}

	data, err := ioutil.ReadFile(b.caPath)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()

	pool.AppendCertsFromPEM(data)

	tlscfg := &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: b.skipVerify,
	}
	if certs != nil {
		tlscfg.Certificates = certs
	}

	tr := &http.Transport{
		TLSClientConfig:       tlscfg,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return tr, nil

}
