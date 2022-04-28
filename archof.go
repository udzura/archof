package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func main() {
	token := ""
	flag.StringVar(&token, "bearer", "", "")
	flag.Parse()
	target := flag.Arg(0)

	ref, err := name.ParseReference(target)
	if err != nil {
		panic(err)
	}

	img, err := remote.Image(ref, remote.WithAuth(&authn.Bearer{Token: token}))
	if err != nil {
		panic(err)
	}
	// fmt.Println("Image:", img)

	manifest, err := img.Manifest()
	if err != nil {
		panic(err)
	}
	sha := manifest.Config.Digest.String()

	desc, err := remote.Get(ref, remote.WithAuth(&authn.Bearer{Token: token}))
	if err != nil {
		panic(err)
	}

	reg := ref.Context()
	url := fmt.Sprintf("%s://%s/v2/%s/blobs/%s",
		reg.Scheme(),
		reg.RegistryStr(),
		reg.RepositoryStr(),
		sha,
	)
	res, err := desc.Client.Get(url)
	if err != nil {
		panic(err)
	}

	type BlobResponse struct {
		Architecture string `json:"architecture"`
	}
	b := BlobResponse{}

	if err := json.NewDecoder(res.Body).Decode(&b); err != nil {
		panic(err)
	}

	fmt.Println(b.Architecture)
}
