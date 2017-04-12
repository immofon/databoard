package databoard

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
)

func TestNew(t *testing.T) {
	d := New(nil, "")
	d.c = (*github.Client)(nil)
	d.Owner = ""
	d.Repo = ""
}

func initFromEnv(t *testing.T, d *Databoard) (*Databoard, bool) {
	d.Owner = os.Getenv("GITHUB_OWNER")
	d.Repo = os.Getenv("GITHUB_REPO")

	if os.Getenv("GITHUB_TOKEN") == "" || d.Owner == "" || d.Repo == "" {
		t.Error("expect $GITHUB_TOKEN,$GITHUB_OWNER,$GITHUB_REPO")
		return nil, false
	}
	return d, true
}

func getDataboardFromEnv(t *testing.T) (d *Databoard, ok bool) {
	d, ok = initFromEnv(t, New(nil, os.Getenv("GITHUB_TOKEN")))
	return
}

func TestDataboard_GetLatestRelease(t *testing.T) {
	d, ok := getDataboardFromEnv(t)
	if !ok {
		return
	}
	_, err := d.GetLatestRelease(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestDataboard_GetReleases(t *testing.T) {
	d, ok := getDataboardFromEnv(t)
	if !ok {
		return
	}
	releases, _, err := d.GetReleases(context.TODO(), 10)
	if err != nil {
		t.Fatal(err)
	}

	for r := range releases {
		fmt.Println(*r.TagName, r.PublishedAt)
	}
}

func TestDataboard_GetReleasesByTag(t *testing.T) {
	d, ok := getDataboardFromEnv(t)
	if !ok {
		return
	}

	_, err := d.GetReleaseByTag(context.TODO(), os.Getenv("GITHUB_REPO_TAG"))
	err = errors.Annotate(err, "expect $GITHUB_REPO_TAG")
	if err != nil {
		t.Fatal(errors.ErrorStack(err))
	}
}
