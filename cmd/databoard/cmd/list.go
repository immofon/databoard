package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var (
		all   = false
		first = 20
	)

	listCmd.Flags().BoolVarP(&all, "all", "", all, "show all of releases")
	listCmd.Flags().IntVarP(&first, "first", "", first, "show first size of releases")

	listCmd.Run = func(cmd *cobra.Command, args []string) {
		switch {
		case owner == "":
			exit(1, "require owner")
		case repo == "":
			exit(2, "require repo")
		case token == "":
			exit(3, "require token")
		}

		var (
			d       = newDataboard()
			ctx     = context.Background()
			rch     <-chan *github.RepositoryRelease
			err     error
			perpage = 20
		)
		if all {
			rch, _, err = d.GetReleases(ctx, perpage)

		} else {
			if first < 1 {
				first = 1
			}
			if perpage > first {
				perpage = first
			}
			rch, _, err = d.GetReleases(ctx, perpage)
		}

		if err != nil {
			exit(4, errors.ErrorStack(errors.Trace(err)))
		}

		var size = 0
		for release := range rch {
			if size >= first {
				break
			}

			var asset *github.ReleaseAsset
			for _, a := range release.Assets {
				if a.GetName() == "data.tar.gz.gpg" {
					asset = &a
					size++
					break
				}
			}

			if asset == nil {
				continue
			}

			n, unit := Bytesize(float64(asset.GetSize()))
			assetSize := fmt.Sprintf("%6.2f %s", n, unit)
			fmt.Printf("%12s    %s\n", assetSize, release.GetTagName())
		}
	}
}

func Bytesize(size float64) (n float64, unit string) {
	if size/1024 < 1 { // B
		n = size
		unit = "B"
		return
	}
	if size/1024/1024 < 1 { // KB
		n = size / 1024
		unit = "KB"
		return
	}
	if size/1024/1024/1024 < 1 { // MB
		n = size / 1024 / 1024
		unit = "MB"
		return
	}

	n = size / 1024 / 1024 / 1024
	unit = "GB"

	return
}
