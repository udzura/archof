package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

const VERSION = "0.1.0"

func usage(err error) {
	fmt.Println(err.Error())

	format := `
usage:
    %s [container.reg.example/org/image:tag] [OPTIONS]

flags:
    --bearer <token> A http bearer auth token to pass to registory

version:
    v%s
`
	fmt.Printf(format, os.Args[0], VERSION)
	os.Exit(1)
}

func parseFlags(image *string, token *string) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Arguments too short")
	}
	var args []string
	if !strings.HasPrefix(os.Args[1], "-") {
		*image = os.Args[1]
		args = os.Args[2:]
	} else {
		args = os.Args[1:]
	}

	fs := flag.NewFlagSet("Flags", flag.ContinueOnError)
	fs.StringVar(token, "bearer", "", "")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		*image = fs.Arg(0)
	}

	if *image == "" {
		return fmt.Errorf("Invalid arguments")
	}
	return nil
}

func main() {
	target := ""
	token := ""
	if err := parseFlags(&target, &token); err != nil {
		usage(err)
	}

	ref, err := name.ParseReference(target)
	if err != nil {
		panic(err)
	}

	var img v1.Image
	if token == "" {
		img, err = remote.Image(ref)
	} else {
		img, err = remote.Image(ref, remote.WithAuth(&authn.Bearer{Token: token}))
	}
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
