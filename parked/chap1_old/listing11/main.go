package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var Version string = "0.1"

type versionInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Url     string `json:"url"`
}

func versionsAvailable() ([]versionInfo, error) {
	var versions []versionInfo
	versionMetadataUrl := "https://gist.githubusercontent.com/amitsaha/1d65876a0a2d9b5bae5a299a94dffae7/raw/7b9b6ab54c9434ac34731cc457a1bd10a6e2aab1/versions.json"
	resp, err := http.Get(versionMetadataUrl)
	if err != nil {
		return versions, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &versions)
	return versions, err
}

func downloadAndReplace(v versionInfo) error {
	resp, err := http.Get(v.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	f, err := os.Create(v.Name)
	if err != nil {
		return err
	}
	_, err = f.Write(body)
	if err != nil {
		return err
	}

	currentBinPath, err := os.Executable()
	if err != nil {
		return err
	}
	p, err := filepath.EvalSymlinks(currentBinPath)
	if err != nil {
		return err
	}
	fmt.Println(p)
	// rename current executable to a new executable name
	tmpExec := path.Dir(p) + ".tmp"
	err = os.Rename(p, tmpExec)
	if err != nil {
		return err
	}
	err = os.Rename(v.Name, p)
	if err != nil {
		return err
	}

	err = os.Remove(tmpExec)
	return err
}

func main() {
	fmt.Printf("Version: %s\n", Version)
	v, err := versionsAvailable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Latest version: %s\n", v[0].Version)
	err = downloadAndReplace(v[0])
	fmt.Println(err)
}
