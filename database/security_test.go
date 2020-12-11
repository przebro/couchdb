package database

import (
	"context"
	"testing"
)

func TestGetSecurity(t *testing.T) {

	_, database, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}
	res, err := database.Security(context.Background())
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if res.Code != 200 {
		t.Error("unexpected result")
	}
}
func TestSetAdminSecurity(t *testing.T) {

	_, database, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	res, err := database.SetAdminSecurity(context.Background(), nil, []string{})
	if err != errSecurityDataEmpty {
		t.Error("unexpected result:", err)
	}

	res, err = database.SetAdminSecurity(context.Background(), []string{}, nil)
	if err != errSecurityDataEmpty {
		t.Error("unexpected result:", err)
	}

	res, err = database.SetAdminSecurity(context.Background(), []string{"user"}, []string{"role"})
	if err != nil {
		t.Error("unexpected result:")
	}

	if res.Code != 200 {
		t.Error("unexpected result")
	}

	expect := map[string]bool{}

	err = res.Decode(&expect)
	if err != nil {
		t.Error("unexpected result")
	}

}

func TestSetMemberSecurity(t *testing.T) {

	_, database, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	res, err := database.SetMemberSecurity(context.Background(), nil, []string{})
	if err != errSecurityDataEmpty {
		t.Error("unexpected result:", err)
	}

	res, err = database.SetMemberSecurity(context.Background(), []string{}, nil)
	if err != errSecurityDataEmpty {
		t.Error("unexpected result:", err)
	}

	res, err = database.SetMemberSecurity(context.Background(), []string{"user"}, []string{"role"})
	if err != nil {
		t.Error("unexpected result:")
	}

	if res.Code != 200 {
		t.Error("unexpected result")
	}

	expect := map[string]bool{}

	err = res.Decode(&expect)
	if err != nil {
		t.Error("unexpected result")
	}

}
