package databoard

import (
	"context"

	"github.com/google/go-github/github"
)

type Databoard struct {
	c     *github.Client
	Token string
	Owner string
	Repo  string
}

func New(c *github.Client) *Databoard {
	return &Databoard{
		c: c,
	}
}

func (d *Databoard) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	ret, _, err := d.c.Repositories.GetLatestRelease(ctx, d.Owner, d.Repo)
	return ret, err
}
