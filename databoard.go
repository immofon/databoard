package databoard

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
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

	hc4github := &http.Client{
		Transport: &oauth2.Transport{
			Base: hc.Transport,
			Source: oauth2.ReuseTokenSource(nil, oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)),
		},
		CheckRedirect: hc.CheckRedirect,
		Jar:           hc.Jar,
		Timeout:       hc.Timeout,
	}

	client := github.NewClient(hc4github)

	return &Databoard{
		c: client,
	}
}

func (d *Databoard) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	ret, _, err := d.c.Repositories.GetLatestRelease(ctx, d.Owner, d.Repo)
	return ret, err
}

func (d *Databoard) GetReleases(ctx context.Context, perpage int) (rch <-chan *github.RepositoryRelease, cancel func(), err error) {
	ctx, cancel = context.WithCancel(ctx)
	if perpage < 1 {
		perpage = 1
	}

	ch := make(chan *github.RepositoryRelease, perpage)

	go func(ch chan<- *github.RepositoryRelease) {
		defer close(ch)
		var (
			listOptions = &github.ListOptions{
				Page:    1,
				PerPage: perpage,
			}
		)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

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
	return ch, cancel, nil
}

func (d *Databoard) UploadReleaseAsset(ctx context.Context, release_id int, file *os.File, name string) (*github.ReleaseAsset, error) {
	if name == "" {
		name = file.Name()
	}
	asset, _, err := d.c.Repositories.UploadReleaseAsset(ctx, d.Owner, d.Repo, release_id, &github.UploadOptions{
		Name: name,
	}, file)

	return asset, err
}

func (d *Databoard) GetReleaseByTag(ctx context.Context, tag string) (*github.RepositoryRelease, error) {
	if tag == "" {
		return nil, errors.New("empty tag")
	}

	release, _, err := d.c.Repositories.GetReleaseByTag(ctx, d.Owner, d.Repo, tag)
	return release, errors.Annotatef(err, "github.RepositoriesService.GetReleaseByTag <tag: %q>", tag)
}

func (d *Databoard) CreateRelease(ctx context.Context, tag string) (*github.RepositoryRelease, error) {
	if tag == "" {
		return nil, errors.New("empty tag")
	}

	release, _, err := d.c.Repositories.CreateRelease(ctx, d.Owner, d.Repo, &github.RepositoryRelease{
		TagName: github.String(tag),
	})

	return release, errors.Annotatef(err, "github.RepositoriesService.CreateRelease <tag: %q>", tag)
}
