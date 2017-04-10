package databoard

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestNew(t *testing.T) {
	d := New(nil)
	d.c = (*github.Client)(nil)
	d.Token = ""
	d.Owner = ""
	d.Repo = ""
}
