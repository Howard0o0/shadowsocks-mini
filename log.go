package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// dir should be an absolute path
func createDir(dir string) error {
	_, err := os.Stat(dir)

	if err == nil {
		//directory exists
		return nil
	}

	err2 := os.MkdirAll(dir, 0666)
	if err2 != nil {
		return err2
	}

	return nil
}

type logfileHook struct {
	dir string
}

func (hk *logfileHook) Fire(entry *log.Entry) error {
	date := time.Now().Format("2006-01-02")
	time := time.Now().Format("15:04:05")

	f, err := os.OpenFile(hk.dir+"/"+date+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		f.Write([]byte(entry.Level.String() + "\t[" + date + " " + time + "]\t" + entry.Message + "\n"))
	}
	defer f.Close()
	return err
}

func (df *logfileHook) Levels() []log.Level {
	return log.AllLevels[:5] // level >= InfoLevel
}

type warnHook struct {
	logfileHook
}

func (hk *warnHook) Fire(entry *log.Entry) error {
	date := time.Now().Format("2006-01-02")
	time := time.Now().Format("15:04:05")

	f, err := os.OpenFile(hk.dir+"/"+"WARN-"+date+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		f.Write([]byte(entry.Level.String() + "\t[" + date + " " + time + "]\t" + entry.Message + "\n"))
	}
	defer f.Close()
	return err
}

func (hk *warnHook) Levels() []log.Level {
	return log.AllLevels[3:4] // level >= InfoLevel
}

func init() {

	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true                    // 显示完整时间
	customFormatter.TimestampFormat = "2006-01-02 15:04:05" // 时间格式
	customFormatter.DisableTimestamp = false                // 禁止显示时间
	customFormatter.DisableColors = false                   // 禁止颜色显示
	logrus.SetFormatter(customFormatter)

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)

}

func setLogDir(dir string) error {

	hook := &logfileHook{dir: dir}
	logrus.AddHook(hook)
	logrus.AddHook(&warnHook{logfileHook: *hook})

	return createDir(dir)
}
