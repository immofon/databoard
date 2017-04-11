package databoard

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/github"
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
