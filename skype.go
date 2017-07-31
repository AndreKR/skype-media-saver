package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type SkypeProfileDirectory struct {
	Name       string
	MediaCache string
}

// findSkypeMediaCaches returns one or more (if the user has multiple Skype profiles) paths that could be Skype media caches
func findSkypeMediaCaches() []SkypeProfileDirectory {

	var profileDirectories []SkypeProfileDirectory

	appdata := os.Getenv("AppData")
	if appdata == "" {
		log.Fatalln("%AppData% is not set, I have no way of finding the Skype media cache")
	}

	files, err := ioutil.ReadDir(filepath.Join(appdata, "Skype"))
	if err != nil {
		log.Fatalln("Can't read %AppData%\\Skype: " + err.Error())
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		// No one will have those as their Sykpe username
		if f.Name() == "Content" || f.Name() == "DataRv" || f.Name() == "logs" || f.Name() == "shared_dynco" || f.Name() == "shared_httpfe" || f.Name() == "SkypeRT" {
			continue
		}

		_, err := os.Stat(filepath.Join(appdata, "Skype", f.Name(), "media_messaging/media_cache_v3"))
		if err != nil {
			continue
		}

		profileDirectories = append(profileDirectories, SkypeProfileDirectory{f.Name(), filepath.Join(appdata, "Skype", f.Name(), "media_messaging/media_cache_v3")})
	}

	return profileDirectories
}
