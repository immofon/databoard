package databoard

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type Databoard struct {
	c     *github.Client
	Owner string
	Repo  string
}

func New(hc *http.Client, token string) *Databoard {
	if hc == nil {
		hc = &http.Client{}
	}
	hc.Transport = &oauth2.Transport{
		Base: hc.Transport,
		Source: oauth2.ReuseTokenSource(nil, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)),
	}

	client := github.NewClient(hc)

	return &Databoard{
		c: client,
	}
}

func (d *Databoard) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	ret, _, err := d.c.Repositories.GetLatestRelease(ctx, d.Owner, d.Repo)
	return ret, err
}

func (d *Databoard) GetReleases(ctx context.Context) (<-chan *github.RepositoryRelease, error) {
	ch := make(chan *github.RepositoryRelease)

	go func(ch chan<- *github.RepositoryRelease) {
		defer close(ch)
		listOptions := &github.ListOptions{
			Page:    1,
			PerPage: 1,
		}
		for {
			releases, resp, err := d.c.Repositories.ListReleases(ctx, d.Owner, d.Repo, listOptions)
			if err != nil {
				return
			}

			for _, rel := range releases {
				select {
				case <-ctx.Done():
					return
				case ch <- rel:
				}
			}

			if resp.NextPage == 0 {
				return
			}

			listOptions.Page = resp.NextPage
		}
	}(ch)
	return ch, nil
}
