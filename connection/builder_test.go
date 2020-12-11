package connection

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/przebro/couchdb/client"
)

func TestBuildNoUser(t *testing.T) {
	builder := NewBuilder()
	_, err := builder.WithAddress(host, port).WithAuthentication(client.Cookie, "", "").Build(false)
	if err == nil {
		t.Error("Unexpected result")
	}

}

func TestRootDir(t *testing.T) {
	dir, _ := os.Getwd()
	pwd := os.Getenv("PWD")
	fmt.Println("getpwd", dir, ":", pwd)
}

func TestSecureConConn(t *testing.T) {

	dir, _ := os.Getwd()
	caPath := filepath.Join(filepath.Dir(dir), "docker", "etc", "root_ca.crt")
	ckeyPath := filepath.Join(filepath.Dir(dir), "docker", "etc", "client.key")
	ccertPath := filepath.Join(filepath.Dir(dir), "docker", "etc", "client.crt")

	builder := NewBuilder()
	conn, err := builder.WithAuthentication(client.Cookie, username, password).
		WithCertificate(caPath, ckeyPath, ccertPath, true).
		WithAddress(host, sport).Build(true)
	if err != nil {
		t.Error(err)
	}

	result, _ := conn.Up(context.TODO())
	if result.Code != 200 {
		t.Error("unexpected result")
	}

	conn, err = builder.WithAuthentication(client.Cookie, username, password).
		WithCertificate(caPath, ckeyPath, "", true).
		WithAddress(host, sport).Build(false)

	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	_, err = conn.Session(ctx, username, password)

	if err == nil {
		t.Error(err)
	}
	cancel()

	conn, err = builder.WithAuthentication(client.Cookie, username, password).
		WithCertificate(caPath, ckeyPath, "not_valid_path", true).
		WithAddress(host, sport).Build(false)

	if err == nil {
		t.Error(err)
	}

	conn, err = builder.WithAuthentication(client.Cookie, username, password).
		WithCertificate("", ckeyPath, ccertPath, true).
		WithAddress(host, sport).Build(false)

	if err == nil {
		t.Error(err)
	}

}
