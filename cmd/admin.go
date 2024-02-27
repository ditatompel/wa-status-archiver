package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"wabot/internal/database"
	"wabot/internal/repo/admin"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Create Admin",
	Long:  `Create an admin account for WebUI access.`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: wabot admin create")
			os.Exit(1)
		}

		if args[0] == "create" {
			repo := admin.NewAdminRepo(database.GetDB())
			a := admin.Admin{
				Username: stringPrompt("Username:"),
				Password: passPrompt("Password:"),
			}

			_, err := repo.CreateAdmin(&a)
			if err != nil {
				fmt.Println("Error creating admin:", err)
				os.Exit(1)
			}

			fmt.Println("Admin created successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
}

func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func passPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}
