package main

import (
	"fmt"
	"os"
	"strconv"
	"userctl/utils"

	"github.com/spf13/cobra"
)

var (
	url     string
	basedn  string
	admin   string
	adminpw string
)

var (
	cliName        = "userctl"
	cliDescription = "A simple command line tool for user manage."
	rootCmd        = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"userctl"},
	}
)

func userCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user <subcommand>",
		Short: "user related commands",
	}
	cmd.AddCommand(getAllUsersCommand())
	cmd.AddCommand(getUserByIDCommand())
	cmd.AddCommand(getUserByNameCommand())
	cmd.AddCommand(addUserCommand())
	cmd.AddCommand(modUserPwdCommand())
	cmd.AddCommand(delUserCommand())
	return cmd
}

func getAllUsersCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "get all users",
		Run:   getAllUsers,
	}
	return &cmd
}

func getUserByIDCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "id <id>",
		Short: "get user through ID",
		Run:   getUserByID,
	}
	return &cmd
}

func getUserByNameCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "name <name>",
		Short: "get user through name",
		Run:   getUserByName,
	}
	return &cmd
}

func addUserCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add <name> <id> <password>",
		Short: "add user",
		Run:   addUser,
	}
	return &cmd
}

func delUserCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "del <name>",
		Short: "del user",
		Run:   delUser,
	}
	return &cmd
}

func modUserPwdCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "putpwd <name> <password>",
		Short: "mod password of user",
		Run:   modUserPwd,
	}
	return &cmd
}

func getAllUsers(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}
	data, err := utils.GetUsers(client)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(data)
}

func getUserByID(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		os.Exit(1)
	}
	data, err := utils.GetUserByID(client, id)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(data)
}

func getUserByName(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	data, err := utils.GetUserByName(client, args[0])
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(data)
}

func addUser(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	err := utils.AddUser(client, args[0], args[1], args[2])
	if err != nil {
		os.Exit(1)
	}
}

func delUser(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	err := utils.DelUser(client, args[0])
	if err != nil {
		os.Exit(1)
	}
}

func modUserPwd(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	err := utils.ModUserPwd(client, args[0], args[1])
	if err != nil {
		os.Exit(1)
	}
}

func groupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group <subcommand>",
		Short: "group related commands",
	}
	cmd.AddCommand(getAllGroupsCommand())
	cmd.AddCommand(getGroupByNameCommand())
	cmd.AddCommand(addGroupCommand())
	cmd.AddCommand(delGroupCommand())
	cmd.AddCommand(addGroupMemberCommand())
	cmd.AddCommand(delGroupMemberCommand())
	return cmd
}

func getAllGroupsCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "get all groups",
		Run:   getAllGroups,
	}
	return &cmd
}

func getGroupByNameCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "name <name>",
		Short: "get group through name",
		Run:   getGroupByName,
	}
	return &cmd
}

func addGroupCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add <name> <id>",
		Short: "add group",
		Run:   addGroup,
	}
	return &cmd
}

func delGroupCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "del <name>",
		Short: "del group",
		Run:   delGroup,
	}
	return &cmd
}

func addGroupMemberCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "addMember <groupname> <username>",
		Short: "add user to group",
		Run:   addGroupMember,
	}
	return &cmd
}

func delGroupMemberCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "delMember <groupname> <username>",
		Short: "del user from group",
		Run:   delGroupMember,
	}
	return &cmd
}

func getAllGroups(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}
	data, err := utils.GetGroups(client)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(data)
}

func getGroupByName(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}
	data, err := utils.GetGroupByName(client, args[0])
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(data)
}

func addGroup(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	err := utils.AddGroup(client, args[0], args[1])
	if err != nil {
		os.Exit(1)
	}
}

func delGroup(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	err := utils.DelGroup(client, args[0])
	if err != nil {
		os.Exit(1)
	}
}

func addGroupMember(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	err := utils.GroupAddMember(client, args[0], args[1])
	if err != nil {
		os.Exit(1)
	}
}

func delGroupMember(cmd *cobra.Command, args []string) {
	client := &utils.LDAPClient{
		Addr:     url,
		BaseDn:   basedn,
		BindDn:   admin,
		BindPass: adminpw,
		TLS:      false,
		StartTLS: false}

	err := utils.GroupDelMember(client, args[0], args[1])
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	rootCmd.AddCommand(userCommand())
	rootCmd.AddCommand(groupCommand())
	rootCmd.PersistentFlags().StringVar(&url, "url", "127.0.0.1:389", "ldap address")
	rootCmd.PersistentFlags().StringVar(&basedn, "baseDn", "dc=test,dc=com", "ldap basedn")
	rootCmd.PersistentFlags().StringVar(&admin, "admin", "cn=manager,dc=test,dc=com", "ldap admin")
	rootCmd.PersistentFlags().StringVar(&adminpw, "adminPw", "123456", "ldap admin password")
	rootCmd.Execute()
}
