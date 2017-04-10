package databorad

import (
	"context"

	"github.com/google/go-github/github"
)

type Databorad struct {
	c     *github.Client
	Token string
	Owner string
	Repo  string
}

func NewDataboard(c *github.Client) *Databorad {
	return &Databorad{
		c: c,
	}
}

func (d *Databorad) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	ret, _, err := d.c.Repositories.GetLatestRelease(ctx, d.Owner, d.Repo)
	return ret, err
}
