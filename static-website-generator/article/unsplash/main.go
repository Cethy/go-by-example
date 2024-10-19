package unsplash

import (
	"context"
	"github.com/hbagdi/go-unsplash/unsplash"
	"golang.org/x/oauth2"
)

func GetRandomPhotoUrl(accessToken string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "Client-ID " + accessToken},
	)
	client := oauth2.NewClient(ctx, ts)
	u := unsplash.New(client)
	randomPhoto, _, err := u.Photos.Random(&unsplash.RandomPhotoOpt{
		Width:       1200,
		Height:      800,
		SearchQuery: "code",
		Orientation: "landscape",
	})
	if err != nil {
		return "", err
	}
	slice := *randomPhoto

	return slice[0].Urls.Raw.URL.String(), nil
}
