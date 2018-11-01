package utils

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/md4"
	"golang.org/x/text/encoding/unicode"
	ldap "gopkg.in/ldap.v2"
)

var (
	sambadomain = "SAMBA"
)

// LdapResult ... type
type LdapResult struct {
	DN         string              `json:"dn"`
	Attributes map[string][]string `json:"attributes"`
}

// LDAPClient ... type
type LDAPClient struct {
	Addr     string
	BaseDn   string
	BindDn   string
	BindPass string
	TLS      bool
	StartTLS bool
	Conn     *ldap.Conn
}

func createSambaNtpPwd(password string) (encpwd string, err error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	encoded, err := utf16.NewEncoder().String(password)
	if err != nil {
		return
	}
	h := md4.New()
	io.WriteString(h, encoded)
	encpwd = hex.EncodeToString(h.Sum(nil))
	return
}

// Close ... ldap close
func (lc *LDAPClient) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}

// Connect ... ldap connect
func (lc *LDAPClient) Connect() (err error) {
	if lc.TLS {
		lc.Conn, err = ldap.DialTLS("tcp", lc.Addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		lc.Conn, err = ldap.Dial("tcp", lc.Addr)
	}
	if err != nil {
		return err
	}
	if !lc.TLS && lc.StartTLS {
		err = lc.Conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			lc.Conn.Close()
			return err
		}
	}

	err = lc.Conn.Bind(lc.BindDn, lc.BindPass)
	if err != nil {
		lc.Conn.Close()
		return err
	}
	return err
}

// Search ... get groups or users
func (lc *LDAPClient) Search(filter string, attr []string, basedn string) (data []LdapResult, err error) {
	searchRequest := ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attr,
		nil,
	)
	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return
	}
	if len(sr.Entries) == 0 {
		err = errors.New("Cannot find such group")
		return
	}
	results := []LdapResult{}
	for _, entry := range sr.Entries {
		attributes := make(map[string][]string)
		for _, attr := range entry.Attributes {
			attributes[attr.Name] = attr.Values
		}

		var result LdapResult
		result.DN = entry.DN
		result.Attributes = attributes

		results = append(results, result)
	}
	data = results
	return
}

// Mod ... mod attr
func (lc *LDAPClient) Mod(basedn string, opt string, attrKey string, attrValue []string) (err error) {
	modify := ldap.NewModifyRequest(basedn)
	if opt == "add" {
		modify.Add(attrKey, attrValue)
	} else if opt == "del" {
		modify.Delete(attrKey, attrValue)
	} else if opt == "Replace" {
		modify.Replace(attrKey, attrValue)
	} else {
		err = errors.New("opt error")
		return
	}

	err = lc.Conn.Modify(modify)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// Exist ... check user or group exist
func (lc *LDAPClient) Exist(filter string) bool {
	searchRequest := ldap.NewSearchRequest(
		lc.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{},
		nil,
	)
	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return false
	}
	if len(sr.Entries) != 0 {
		return true
	}
	return false
}

// SambadomainSid ... get domain sid
func (lc *LDAPClient) SambadomainSid() (sid string, err error) {
	filter := fmt.Sprintf("(sambaDomainName=%s)", sambadomain)
	attrs := []string{"sambaSID"}
	data, err := lc.Search(filter, attrs, lc.BaseDn)
	sid = data[0].Attributes["sambaSID"][0]
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}

	return
}

// AddUser ... add user
func (lc *LDAPClient) AddUser(username string, uidStr string, passwd string) (err error) {
	if lc.Exist("(&(uid=" + username + "))") {
		return errors.New("record has existed in ldap")
	}

	domainID, err := lc.SambadomainSid()
	if err != nil {
		return
	}
	curtime := fmt.Sprintf("%d", time.Now().Unix())
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return
	}
	sambaSid := fmt.Sprintf("%s-%d", domainID, uid*2+1000)
	userDn := fmt.Sprintf("uid=%s,ou=People,%s", username, lc.BaseDn)
	userAttr := make(map[string][]string)
	userAttr["uid"] = []string{username}
	userAttr["shadowMin"] = []string{"0"}
	userAttr["homeDirectory"] = []string{"/home/" + username}

	userAttr["objectClass"] = []string{"top", "person", "organizationalPerson",
		"inetOrgPerson", "sambaSamAccount", "posixAccount", "shadowAccount"}
	userAttr["cn"] = []string{username}
	userAttr["givenName"] = []string{username}
	userAttr["sn"] = []string{username}
	userAttr["uid"] = []string{username}
	userAttr["uidNumber"] = []string{uidStr}
	userAttr["gidNumber"] = []string{"0"}
	userAttr["displayName"] = []string{username}
	userAttr["homeDirectory"] = []string{"/home/" + username}
	userAttr["loginShell"] = []string{"/bin/bash"}
	userAttr["sambaSID"] = []string{sambaSid}
	userAttr["sambaAcctFlags"] = []string{"[U ]"}
	userAttr["userPassword"] = []string{passwd}
	// gen ntp pwd
	ntppwd, err := createSambaNtpPwd(passwd)
	if err != nil {
		return
	}
	userAttr["sambaNTPassword"] = []string{ntppwd}
	userAttr["sambaPwdLastSet"] = []string{curtime}

	addrequest := ldap.NewAddRequest(userDn)
	for k, v := range userAttr {
		addrequest.Attribute(k, v)
	}
	if err = lc.Conn.Add(addrequest); err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	passwordModifyRequest := ldap.NewPasswordModifyRequest(userDn, "", passwd)
	_, err = lc.Conn.PasswordModify(passwordModifyRequest)

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// ModifyPwd ... change pwd of user
func (lc *LDAPClient) ModifyPwd(username, password string) (err error) {
	user := fmt.Sprintf("uid=%s,ou=People,%s", username, lc.BaseDn)
	passwordModifyRequest := ldap.NewPasswordModifyRequest(user, "", password)
	_, err = lc.Conn.PasswordModify(passwordModifyRequest)

	if err != nil {
		log.Fatalf("Password could not be changed: %s", err.Error())
		return
	}

	ntppwd, err := createSambaNtpPwd(password)
	if err != nil {
		return
	}

	modify := ldap.NewModifyRequest(user)
	modify.Replace("sambaNTPassword", []string{ntppwd})
	err = lc.Conn.Modify(modify)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// DelUser ... del user
func (lc *LDAPClient) DelUser(username string) (err error) {
	userDn := fmt.Sprintf("uid=%s,ou=People,%s", username, lc.BaseDn)
	delrequest := ldap.NewDelRequest(userDn, nil)
	err = lc.Conn.Del(delrequest)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// DelGroup ... del group
func (lc *LDAPClient) DelGroup(groupname string) (err error) {
	userDn := fmt.Sprintf("cn=%s,ou=Group,%s", groupname, lc.BaseDn)
	delrequest := ldap.NewDelRequest(userDn, nil)
	err = lc.Conn.Del(delrequest)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// AddGroup ... add group
func (lc *LDAPClient) AddGroup(groupname string, gidStr string) (err error) {
	domainID, err := lc.SambadomainSid()
	if err != nil {
		return
	}
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		return
	}
	sambaSid := fmt.Sprintf("%s-%d", domainID, gid*2+1000)
	groupDn := fmt.Sprintf("cn=%s,ou=Group,%s", groupname, lc.BaseDn)
	groupAttr := make(map[string][]string)

	groupAttr["objectClass"] = []string{"top", "posixGroup", "sambaGroupMapping"}
	groupAttr["cn"] = []string{groupname}
	groupAttr["gidNumber"] = []string{gidStr}
	groupAttr["sambaSID"] = []string{sambaSid}
	groupAttr["sambaGroupType"] = []string{"2"}

	addrequest := ldap.NewAddRequest(groupDn)
	for k, v := range groupAttr {
		addrequest.Attribute(k, v)
	}
	if err = lc.Conn.Add(addrequest); err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// GetUsers ... get users
func GetUsers(lc *LDAPClient) (data string, err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	filter := "(objectClass=sambaSamAccount)"
	attrs := []string{"uid", "uidNumber"}

	basedn := fmt.Sprintf("ou=People,%s", lc.BaseDn)
	users, err := lc.Search(filter, attrs, basedn)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	usersbytes, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	data = string(usersbytes)
	return
}

// GetUserByName ... get user throuth name
func GetUserByName(lc *LDAPClient, name string) (data string, err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	filter := fmt.Sprintf("(&(uid=%s))", name)
	attrs := []string{}
	basedn := fmt.Sprintf("ou=People,%s", lc.BaseDn)
	user, err := lc.Search(filter, attrs, basedn)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	userbytes, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	data = string(userbytes)
	return
}

// GetUserByID ... get user through id
func GetUserByID(lc *LDAPClient, uidNumber int) (data string, err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	filter := fmt.Sprintf("(&(uidNumber=%d))", uidNumber)
	attrs := []string{}
	basedn := fmt.Sprintf("ou=People,%s", lc.BaseDn)
	user, err := lc.Search(filter, attrs, basedn)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	userbytes, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	data = string(userbytes)
	return
}

// AddUser ... add user
func AddUser(lc *LDAPClient, username string, uid string, pwd string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	err = lc.AddUser(username, uid, pwd)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// ModUserPwd ... mod pwd of user
func ModUserPwd(lc *LDAPClient, username string, pwd string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	err = lc.ModifyPwd(username, pwd)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// DelUser ... del user
func DelUser(lc *LDAPClient, username string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	err = lc.DelUser(username)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// GetGroups ... get groups
func GetGroups(lc *LDAPClient) (data string, err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	filter := "(objectClass=sambaGroupMapping)"
	attrs := []string{"cn", "gidNumber", "memberUid"}
	basedn := fmt.Sprintf("ou=Group,%s", lc.BaseDn)
	groups, err := lc.Search(filter, attrs, basedn)
	if err != nil {
		return
	}
	groupsbytes, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	data = string(groupsbytes)
	return
}

// GetGroupByName ... get group
func GetGroupByName(lc *LDAPClient, name string) (data string, err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	filter := fmt.Sprintf("(&(cn=%s))", name)
	attrs := []string{}
	basedn := fmt.Sprintf("ou=Group,%s", lc.BaseDn)
	group, err := lc.Search(filter, attrs, basedn)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	groupbytes, err := json.MarshalIndent(group, "", "  ")
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	data = string(groupbytes)
	return
}

// AddGroup ... add group
func AddGroup(lc *LDAPClient, groupname string, gid string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	err = lc.AddGroup(groupname, gid)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// DelGroup ... del user
func DelGroup(lc *LDAPClient, groupname string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	err = lc.DelGroup(groupname)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}

// GroupAddMember ... add user to group
func GroupAddMember(lc *LDAPClient, groupname string, username string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	// check if username existed
	if !lc.Exist("(&(uid=" + username + "))") {
		err = errors.New("username has not existed in ldap")
		fmt.Println("ERROR: ", err.Error())
		return
	}

	basedn := fmt.Sprintf("cn=%s,ou=Group,%s", groupname, lc.BaseDn)
	err = lc.Mod(basedn, "add", "memberUid", []string{username})
	return
}

// GroupAddMember ... del user from group
func GroupDelMember(lc *LDAPClient, groupname string, username string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	basedn := fmt.Sprintf("cn=%s,ou=Group,%s", groupname, lc.BaseDn)
	err = lc.Mod(basedn, "del", "memberUid", []string{username})
	return
}
