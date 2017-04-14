package cmd

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/immofon/databoard"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use: "publish",
}

func init() {
	RootCmd.AddCommand(publishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	publishCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			publishCmd.Usage()
			exit(1)
		}

		switch {
		case owner == "":
			exit(2, "require owner")
		case repo == "":
			exit(3, "require repo")
		case token == "":
			exit(4, "require token")
		}

		passphare, err := getPassphare()
		if err != nil {
			exit(5, "require passphare\n", errors.ErrorStack(errors.Trace(err)))
		}

		// packet data.tar.gz.gpg
		basedir := fmt.Sprintf("/tmp/databoard-%d", os.Getpid())
		defer os.RemoveAll(basedir)
		asset, err := Packet(basedir, "", passphare, args)
		if err != nil {
			exit(6, errors.ErrorStack(errors.Trace(err)))
		}

		// generate tag name via sha1 sum
		tagName, err := sha1file(asset)
		if err != nil {
			exit(7, errors.ErrorStack(errors.Trace(err)))
		}

		tagName = "v." + tagName

		// create new release by tagName
		ctx := context.Background()

		d := databoard.New(nil, token)
		d.Owner = owner
		d.Repo = repo
		release, err := d.CreateRelease(ctx, tagName)
		if err != nil {
			exit(7, errors.ErrorStack(errors.Annotatef(err, "Databoard.CreateRelease <tag_name:%q>", tagName)))
		}
		release_id := *(release.ID)

		// upload data.tar.gz.gpg as release asset
		assetFile, err := os.Open(asset)
		if err != nil {
			exit(8, errors.ErrorStack(errors.Trace(err)))
		}

		_, err = d.UploadReleaseAsset(ctx, release_id, assetFile, "data.tar.gz.gpg")
		if err != nil {
			exit(9, errors.ErrorStack(errors.Annotatef(err, "Databoard.UploadReleaseAsset <release_id:%d>", release_id)))
		}
	}
}

func sha1file(file string) (sum string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Trace(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
