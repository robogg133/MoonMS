package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const VERSION_MANIFEST = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

type versionManifest struct {
	Latest struct {
		Release string `json:"release"`
	} `json:"latest"`

	Versions []struct {
		ID  string `json:"id"`
		Url string `json:"url"`
	} `json:"versions"`
}

type versionJson struct {
	Downloads struct {
		Server struct {
			Sha1 string `json:"sha1"`
			Url  string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}

func downloadServerJar(targetRelease *string) []byte {

	resp, err := http.Get(VERSION_MANIFEST)
	if err != nil {
		panic(err)
	}
	fmt.Println("=> (1/4) Download all versions json manifest from:", VERSION_MANIFEST)

	decoder := json.NewDecoder(resp.Body)

	var manifest versionManifest

	if err := decoder.Decode(&manifest); err != nil {
		resp.Body.Close()
		panic(err)
	}
	resp.Body.Close()
	resp = nil

	var dlUrl string

	if *targetRelease == "latest" {
		*targetRelease = manifest.Latest.Release
	}

	fmt.Println("=> (2/4) Searching for the latest release ")
	for _, v := range manifest.Versions {
		if v.ID == *targetRelease {
			dlUrl = v.Url
		}
	}

	resp, err = http.Get(dlUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("=> (3/4) Download latest release manifest from:", dlUrl)

	decoder = json.NewDecoder(resp.Body)

	var release versionJson

	if err := decoder.Decode(&release); err != nil {
		panic(err)
	}
	fmt.Println("=> (4/4) Starting server.jar download from:", release.Downloads.Server.Url)
	resp, err = http.Get(release.Downloads.Server.Url)
	if err != nil {
		panic(err)
	}

	serverJar, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	sum := sha1.Sum(serverJar)
	if hex.EncodeToString(sum[:]) != release.Downloads.Server.Sha1 {
		panic("Mismatched sha1 checksum")
	}

	return serverJar
}
