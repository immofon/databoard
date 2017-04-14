package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "databoard",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	initFromEnv()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var (
	passphareString string // default $DATABOARD_PASSPHARE
	owner           string // default $GITHUB_OWNER
	repo            string // default $GITHUB_REPO
	token           string // default $GITHUB_TOKEN
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&passphareString, "passphare", "", "", "AES256 passphare. $DATABOARD_PASSPHARE")
	RootCmd.PersistentFlags().StringVarP(&owner, "owner", "", "", "github user or organization. $GITHUB_OWNER")
	RootCmd.PersistentFlags().StringVarP(&repo, "repo", "", "", "github repository. $GITHUB_REPO")
	RootCmd.PersistentFlags().StringVarP(&token, "token", "", "", "github access token. $GITHUB_TOKEN")
}

func defaultValue(v interface{}, needDefault func(v interface{}) bool, defaultV interface{}) interface{} {
	if needDefault(v) {
		return defaultV
	}
	return v
}

func IsEmptyString(v interface{}) bool {
	s, _ := v.(string)
	return s == ""
}

func initFromEnv() {
	owner = defaultValue(owner, IsEmptyString, os.Getenv("GITHUB_OWNER")).(string)
	repo = defaultValue(repo, IsEmptyString, os.Getenv("GITHUB_REPO")).(string)
	token = defaultValue(token, IsEmptyString, os.Getenv("GITHUB_TOKEN")).(string)
}

func getPassphare() (pass []byte, err error) {
	switch {
	case passphareString != "":
		pass = []byte(passphareString)
	default:
		pass = []byte(os.Getenv("DATABOARD_PASSPHARE"))
		if len(pass) == 0 {
			err = errors.New("not set passphare")
		}
	}

	return
}

var (
	httpClient *http.Client
)

func getHttpClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{}

		if allProxy := os.Getenv("all_proxy"); allProxy != "" {
			httpClient.Transport = func() *http.Transport {
				u, err := url.Parse(allProxy)
				if err != nil {
					return nil
				}

				dialer, err := proxy.FromURL(u, proxy.Direct)
				if err != nil {
					return nil
				}

				return &http.Transport{
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return dialer.Dial(network, addr)
					},
				}
			}()
		}
	}
	return httpClient
}

func exit(code int, v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(code)
}

func exitf(code int, format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	if !strings.HasSuffix(format, "\n") {
		os.Stderr.Write([]byte{'\n'})
	}
	os.Exit(code)
}
