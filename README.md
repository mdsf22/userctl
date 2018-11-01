userctl
Simple ldap tool to retrieve basic information and groups for a user.

Usage:
  userctl [command]
  
Available Commands:
  group       group related commands
  
  help        Help about any command
  
  user        user related commands

Flags:
      --admin string     ldap admin (default "cn=manager,dc=test,dc=com")
      --adminPw string   ldap admin password (default "123456")
      --baseDn string    ldap basedn (default "dc=test,dc=com")
  -h, --help             help for userctl
      --url string       ldap address (default "127.0.0.1:389")

(1) user commands

Usage:
  userctl user [command]

Available Commands:
  add         add user
  del         del user
  id          get user through ID
  list        get all users
  name        get user through name
  putpwd      mod password of user

Flags:
  -h, --help   help for user

Global Flags:
      --admin string     ldap admin (default "cn=manager,dc=test,dc=com")
      --adminPw string   ldap admin password (default "scutech")
      --baseDn string    ldap basedn (default "dc=test,dc=com")
      --url string       ldap address (default "127.0.0.1:389")

(2) group commands
Usage:
  userctl group [command]

Available Commands:
  add         add group
  addMember   add user to group
  del         del group
  delMember   del user from group
  list        get all groups
  name        get group through name

Flags:
  -h, --help   help for group

Global Flags:
      --admin string     ldap admin (default "cn=manager,dc=test,dc=com")
      --adminPw string   ldap admin password (default "scutech")
      --baseDn string    ldap basedn (default "dc=test,dc=com")
      --url string       ldap address (default "127.0.0.1:389")
