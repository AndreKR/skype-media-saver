package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/theckman/go-flock"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var exeDir = getExeDir()

var log = makeLogger()

func main() {

	log.Infoln("Start")

	fl := flock.NewFlock(filepath.Join(exeDir, "lock"))
	locked, _ := fl.TryLock() // ignore errors because a failed lock is also an error
	if !locked {
		// Another instance is already running

		pl := findMyOtherProcesses()
		if len(pl) != 1 {
			log.Fatalln("You are running two instances of Skype Media Saver from two different folders. In this case I can't help you to kill those instances, use the Task Manager.")
			os.Exit(1)
		}

		yes, _ := askMessage("Skype Media Saver", "Another instance is already running (PID: "+strconv.Itoa(pl[0])+"). Do you want to kill it?")

		if yes {
			p, err := os.FindProcess(pl[0])
			if err != nil {
				log.Fatalln("Error getting the other process: " + err.Error())
			}
			err = p.Kill()
			if err != nil {
				log.Fatalln("Error killing the other process: " + err.Error())
			}
		}

		os.Exit(0)
	}

	log.Infoln("Finding Skype media cache directories")

	mcs := findSkypeMediaCaches()

	if len(mcs) == 0 {
		log.Warnln("No media cache directory found")
		os.Exit(1)
	}

	for _, d := range mcs {
		go watch(d)
	}

	select {} // block main goroutine

}

func watch(dir SkypeProfileDirectory) {

	var err error

	log := log.WithField("profile", dir.Name)
	log.Infoln("Watching: " + dir.MediaCache)

	output := filepath.Join(exeDir, "output", dir.Name)
	err = os.MkdirAll(output, 0777)
	if err != nil {
		log.Errorln("Can't create output directory: " + err.Error())
	}
	log.Infoln("Output directory: " + output)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorln("Watching error: " + err.Error())
	}

	err = watcher.Add(dir.MediaCache)
	if err != nil {
		log.Errorln("Watching error: " + err.Error())
	}

	sync(dir.MediaCache, output)

	go func() {
		for err := range watcher.Errors {
			log.Warnln("Error from fsnotify: " + err.Error())
		}
	}()

	for e := range watcher.Events {

		// When Skype downloads an image, the following events happen:
		// Create: ^foo..._fullsize_distr
		// Write: ^foo..._fullsize_distr
		// Write: ^foo..._fullsize_distr
		// ...
		// Create: ^foo..._fullsize_distr.jpg
		// Write: ^foo..._fullsize_distr.jpg

		// So when we get a write event with a file extension, we assume the file to be complete

		if e.Op&fsnotify.Write == fsnotify.Write {
			if strings.HasSuffix(e.Name, "_fullsize_distr.jpg") ||
				strings.HasSuffix(e.Name, "_fullsize_distr.png") ||
				strings.HasSuffix(e.Name, "_video_distr.mp4") {

				name := filepath.Base(e.Name)
				err = copyFile(filepath.Join(output, name), e.Name)
				if err != nil {
					log.Errorln("Error copying file: " + err.Error())
				}
			}
		}

	}
}

func sync(from string, to string) {

	files, err := ioutil.ReadDir(from)
	if err != nil {
		log.Fatalln("Error doing initial sync: " + err.Error())
	}
	for _, f := range files {

		if f.IsDir() {
			continue
		}
		if !(strings.HasSuffix(f.Name(), "_fullsize_distr.jpg") ||
			strings.HasSuffix(f.Name(), "_fullsize_distr.png") ||
			strings.HasSuffix(f.Name(), "_video_distr.mp4")) {
			continue
		}

		dst, err := os.Stat(filepath.Join(to, f.Name()))
		if err != nil || dst.ModTime().Before(f.ModTime()) {
			err = copyFile(filepath.Join(to, f.Name()), filepath.Join(from, f.Name()))
			if err != nil {
				log.Errorln("Error copying file: " + err.Error())
			}
		}
	}

}

func getExeDir() string {
	if os.Getenv("EXEDIR") != "" {
		return os.Getenv("EXEDIR")
	} else {
		return filepath.Dir(os.Args[0])
	}
}

func makeLogger() *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.DebugLevel)

	log.Hooks.Add(DialogHook{})

	f, err := os.OpenFile(exeDir+"/skype-media-saver.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalln("Can't open log file: " + err.Error())
	}

	log.Out = f

	return log
}

type DialogHook struct{}

func (DialogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (DialogHook) Fire(e *logrus.Entry) error {
	showMessage("Skype Media Saver", e.Message, true)
	return nil
}
