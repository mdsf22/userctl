package utils

import (
	"fmt"
	"testing"
)

func Test_getUsers(t *testing.T) {
	client := &LDAPClient{
		Addr:     "127.0.0.1:3899",
		BaseDn:   "dc=test,dc=com",
		BindDn:   "cn=manager,dc=test,dc=com",
		BindPass: "123456",
		TLS:      false,
		StartTLS: false}
	data, err := GetGroups(client)
	if err != nil {
		t.Fatalf("error sending message: %v", err)
	}
	fmt.Println(data)

}

func Test_addUser(t *testing.T) {
	lc := &LDAPClient{
		Addr:     "127.0.0.1:3899",
		BaseDn:   "dc=test,dc=com",
		BindDn:   "cn=manager,dc=test,dc=com",
		BindPass: "123456",
		TLS:      false,
		StartTLS: false}
	err := lc.Connect()
	defer lc.Close()
	if err != nil {
		return
	}
	err = lc.AddUser("test1", "50000", "123456")
	if err != nil {
		t.Fatalf("error sending message: %v", err)
	}
}
