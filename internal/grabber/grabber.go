package grabber

import (
	"background-grabber/internal/background"
	"background-grabber/util"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"
)

type Grabber struct {
	config *Config
}

func New(config *Config) *Grabber {
	return &Grabber{
		config: config,
	}
}

func (g *Grabber) Run() error {

	if err := g.updateBackground(); err != nil {
		return fmt.Errorf("couldn't update: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(time.Duration(g.config.RefreshMinutes) * time.Minute)
	for {
		select {

		case <-interrupt:
			os.Exit(-1)

		case <-ticker.C:

			if g.config.DeleteOldPictures {
				if err := g.deleteOldPictures(); err != nil {
					return fmt.Errorf("couldn't delete old pictures: %v", err)
				}
			}

			if err := g.updateBackground(); err != nil {
				return fmt.Errorf("couldn't update: %v", err)
			}

			log.Info("Waiting until next tick to get new background")
		}
	}

	return nil
}


func deleteFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// is this a dir, is it a plain old file, and does the user have permission to delete this?
	if !info.IsDir() && info.Mode().IsRegular() && info.Mode().Perm() & (1 << (uint(7))) == 0 {
		if rmErr := os.Remove(path); rmErr != nil {
			return rmErr
		}
	}
	return nil
}

func (g *Grabber) deleteOldPictures() error {

	err := filepath.Walk(g.config.BackgroundsDirPath, deleteFile)
	if err != nil {
		return fmt.Errorf("failed to delete files: %v", err)
	}
	return nil
}

func (g *Grabber) updateBackground() error {
	log.Infof("updating background @ %s", time.Now().String())

	if util.Connected() {
		backgrounds, err := g.getNewBackgrounds()
		if err != nil {
			return err
		}

		if err := g.setBackgrounds(backgrounds); err != nil {
			return err
		}
	}

	return nil
}

func (g *Grabber) getNewBackgrounds() (*background.Set, error) {

	log.Debug("Getting new background")

	u := &url.URL{
		Path:   "/photos/random",
		Host:   "api.unsplash.com",
		Scheme: "https",
	}

	u = g.config.addQueryParams(u)
	URL := fmt.Sprintf("%s%s", u.String(), u.Query().Encode())
	log.Debugf("Sending request to get new backgrounds: %s", URL)
	resp, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("could not get images: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	log.Debug("Parsing JSON of backgrounds request")
	backgrounds, err := background.ParseFromJSON(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to get backgrounds from response body: %v", err)
	}

	log.Infof("got %d backgrounds to download", len(*backgrounds))

	return backgrounds, nil

}

func (g *Grabber) setBackgrounds(backgroundSet *background.Set) error {

	log.Debug("Setting backgrounds")
	backgrounds, err := g.downloadBackgrounds(backgroundSet)
	if err != nil {
		return fmt.Errorf("couldn't download backgrounds: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		return g.setWindowsBackground(backgrounds)
	case "linux":
		return g.setUbuntuBackground(backgrounds)
	case "darwin":
		return g.setMacOSBackground(backgrounds)
	default:
		panic(OSNotSupported{})
	}
}

func (g *Grabber) downloadBackgrounds(set *background.Set) (files []string, err error) {

	log.Debug("Downloading backgrounds")

	for _, v := range *set {
		log.Infof("downloading background from %v", v.Urls.Raw)
		resp, err := http.Get(v.Urls.Raw)
		if err != nil {
			return nil, fmt.Errorf("couldn't get the image from its url - %s | err: %v", v.Urls.Raw, err)
		}

		path := g.getNewFileName(v.ID)
		log.Debugf("Creating file at %s", path)
		newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return nil, fmt.Errorf("couldn't create file for new background image to live in: %v", err)
		}

		if _, err := io.Copy(newFile, resp.Body); err != nil {
			return nil, fmt.Errorf("couldn't copy into new file: %v", err)
		}
		log.Infof("copied image into %s", path)
		if err := resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("couldn't close response body: %v", err)
		}

		log.Infof("closing file %s", path)

		if err := newFile.Close(); err != nil {
			return nil, fmt.Errorf("couldn't close file: %v", err)
		}

		files = append(files, path)
	}

	return
}

func (g *Grabber) getNewFileName(id string) string {
	return g.config.BackgroundsDirPath + id + ".png"
}

func (g *Grabber) setUbuntuBackground(files []string) error {

	log.Debug("Using ubuntu method for setting background")

	path := fmt.Sprintf(`file://%s`, files[0])

	var b bytes.Buffer
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", path)
	cmd.Stderr = &b

	log.Infof("running '%s' to set image as background", cmd.Args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("encountered error while changing background: %v", err)
	}

	stderr, err := ioutil.ReadAll(&b)
	if err != nil {
		return fmt.Errorf("could not read from stderr: %v", err)
	}

	if len(stderr) > 0 {
		return fmt.Errorf("stderr from dconf: %s", string(stderr))
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("failed to run successfully")
	} else {
		log.Infof("Successfully set background to: %s", path)
	}

	return nil
}

func (g *Grabber) setMacOSBackground(_ []string) error {

	// set your backgrounds setting to point at your backgrounds folder!

	return nil
}

func (g *Grabber) setWindowsBackground(_ []string) error {

	// set your backgrounds setting to point at your backgrounds folder!

	return nil
}


type OSNotSupported struct{}

func (o OSNotSupported) Error() string {
	return "operating system not supported"
}
