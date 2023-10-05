package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mattermost-pp-migration",
	Short: "Transfer profile pictures from one Mattermost server to another, matching by username",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if the variables are set
		srcServerURL := viper.GetString("src-server-url")
		if srcServerURL == "" {
			return errors.New("src-server-url is not set")
		}
		srcServerURL = strings.TrimSuffix(srcServerURL, "/")

		dstServerURL := viper.GetString("dst-server-url")
		if dstServerURL == "" {
			return errors.New("dst-server-url is not set")
		}
		dstServerURL = strings.TrimSuffix(dstServerURL, "/")

		if viper.GetString("src-access-token") == "" {
			return errors.New("src-access-token is not set")
		}

		if viper.GetString("dst-access-token") == "" {
			return errors.New("dst-access-token is not set")
		}

		srcMMClient := model.NewAPIv4Client(srcServerURL)
		srcMMClient.SetToken(viper.GetString("src-access-token"))
		mustBeOk, _, err := srcMMClient.GetPing(context.TODO())
		if err != nil {
			return fmt.Errorf("unable to ping src server: %w", err)
		}
		if mustBeOk != "OK" {
			return errors.New("unable to ping src server")
		}

		dstMMClient := model.NewAPIv4Client(dstServerURL)
		dstMMClient.SetToken(viper.GetString("dst-access-token"))
		mustBeOk, _, err = dstMMClient.GetPing(context.TODO())
		if err != nil {
			return fmt.Errorf("unable to ping dst server: %w", err)
		}
		if mustBeOk != "OK" {
			return errors.New("unable to ping dst server")
		}

		cfg, _, err := srcMMClient.GetConfig(context.TODO())
		if err != nil {
			return fmt.Errorf("unable to get src config: %w", err)
		}
		if cfg.RateLimitSettings.Enable != nil && *cfg.RateLimitSettings.Enable {
			return errors.New("rate limiting is enabled on the source server. Please turn off while using this tool.")
		}

		cfg, _, err = dstMMClient.GetConfig(context.TODO())
		if err != nil {
			return fmt.Errorf("unable to get dst config: %w", err)
		}
		if cfg.RateLimitSettings.Enable != nil && *cfg.RateLimitSettings.Enable {
			return errors.New("rate limiting is enabled on the destination server. Please turn off while using this tool.")
		}

		page := 0
		perPage := 1000
		errs := map[string]error{}
		for {
			users, _, err := srcMMClient.GetUsers(context.TODO(), page, perPage, "")
			if err != nil {
				return fmt.Errorf("unable to get users: %w", err)
			}

			for _, srcUser := range users {
				dstUser, _, err := dstMMClient.GetUserByUsername(context.TODO(), srcUser.Username, "")
				if err != nil {
					errs[srcUser.Username] = fmt.Errorf("unable to get user %s on dst server: %w", srcUser.Username, err)
					continue
				}

				pic, _, err := srcMMClient.GetProfileImage(context.TODO(), srcUser.Id, "")
				if err != nil {
					errs[srcUser.Username] = fmt.Errorf("unable to get profile image for %s on src server: %w", srcUser.Username, err)
					continue
				}

				_, err = dstMMClient.SetProfileImage(context.TODO(), dstUser.Id, pic)
				if err != nil {
					errs[srcUser.Username] = fmt.Errorf("unable to set profile image for %s on dst server: %w", srcUser.Username, err)
					continue
				}

				log.Println("Updated profile image for", srcUser.Username)
			}

			if len(users) < perPage {
				break
			}
			page++
		}

		if len(errs) > 0 {
			log.Println("Errors:")
			for username, err := range errs {
				log.Println(username, ":", err)
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// source server URL
	rootCmd.Flags().String("src-server-url", "", "Source server URL")
	viper.BindPFlag("src-server-url", rootCmd.Flags().Lookup("src-server-url"))
	// source server token secret
	rootCmd.Flags().String("src-access-token", "", "Source server access token")
	viper.BindPFlag("src-access-token", rootCmd.Flags().Lookup("src-access-token"))
	// source server URL
	rootCmd.Flags().String("dst-server-url", "", "Destination server URL")
	viper.BindPFlag("dst-server-url", rootCmd.Flags().Lookup("dst-server-url"))
	// source server token secret
	rootCmd.Flags().String("dst-access-token", "", "Destination server access token")
	viper.BindPFlag("dst-access-token", rootCmd.Flags().Lookup("dst-access-token"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match
}
