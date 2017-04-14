package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/immofon/databoard"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use: "get",
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var (
		tag string
		dst string
	)

	getCmd.Flags().StringVarP(&tag, "tag", "", "", "tag name of release")
	getCmd.Flags().StringVarP(&dst, "dst", "d", "", "dst dir")

	getCmd.Run = func(cmd *cobra.Command, args []string) {
		dst = defaultValue(dst, IsEmptyString, ".").(string)
		passphare, err := getPassphare()
		if err != nil {
			packetCmd.Usage()
			exit(1, "require passphare\n", errors.ErrorStack(errors.Trace(err)))
		}

		d := databoard.New(nil, token)
		d.Owner = owner
		d.Repo = repo

		// get release assets by tag or latest release assets
		var (
			ctx     = context.Background()
			release *github.RepositoryRelease
		)
		if tag != "" {
			release, err = d.GetReleaseByTag(ctx, tag)
			if err != nil {
				exit(2, errors.ErrorStack(errors.Trace(err)))
			}
		} else { // latest tag
			release, err = d.GetLatestRelease(ctx)
			if err != nil {
				exit(3, errors.ErrorStack(errors.Trace(err)))
			}
			tag = release.GetTagName()
		}

		// download asset named "data.tar.gz.gpg"
		var asset *github.ReleaseAsset
		for _, a := range release.Assets {
			if a.GetName() == "data.tar.gz.gpg" {
				asset = &a
			}
		}

		if asset == nil {
			exitf(4, "no file named data.tar.gz.gpg in release tag %q", tag)
		}

		tmpdir := fmt.Sprintf("/tmp/databoard-%d", os.Getpid())
		defer os.RemoveAll(tmpdir)

		resp, err := http.Get(asset.GetBrowserDownloadURL())
		if err != nil {
			exit(5, errors.ErrorStack(errors.Trace(err)))
		}
		defer resp.Body.Close()

		srcfile := tmpdir + "/data.tar.gz.gpg"
		out, err := os.OpenFile(srcfile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			exit(6, errors.ErrorStack(errors.Trace(err)))
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			exit(7, errors.ErrorStack(errors.Trace(err)))
		}
		out.Close()

		// unpack "data.tar.gz.gpg" to dst dir
		err = Unpack(srcfile, passphare, dst)
		if err != nil {
			exit(8, errors.ErrorStack(errors.Trace(err)))
		}

		exit(0)
	}
}
