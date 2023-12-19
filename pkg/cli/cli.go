package cli

import (
	"fmt"

	"github.com/nsmith5/super-insecure/pkg/client"
	"github.com/nsmith5/super-insecure/pkg/server"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(getItemCmd)
	rootCmd.AddCommand(setItemCmd)
	rootCmd.AddCommand(delItemCmd)

	return rootCmd
}

var rootCmd = &cobra.Command{
	Use: "super-insecure",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := server.New()
		return s.ListenAndServe()
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register new user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Register(cmd.Context(), args[0])
	},
}

var getItemCmd = &cobra.Command{
	Use:   "get",
	Short: "get item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		value, err := client.ItemGet(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}

		fmt.Println("Value:", value)
		return nil
	},
}

var setItemCmd = &cobra.Command{
	Use:   "set",
	Short: "set item",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.ItemSet(cmd.Context(), args[0], args[1], args[2])
	},
}

var delItemCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.ItemDelete(cmd.Context(), args[0], args[1])
	},
}
