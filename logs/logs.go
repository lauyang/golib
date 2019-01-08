package logs

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	LOG_ERROR = iota
	LOG_WARING
	LOG_INFO
	LOG_DEBUG
)

var log *mylog

/*
 * 初始化
 */
func init() {
	log = newMylog()
}

func Init(dir string, file string, level int, savefile bool) {
	log.setDir(dir)
	log.setFile(file)
	log.setLevel(level)
	log.setSavefile(savefile)
}

func Error(err ...interface{}) {
	log.write(LOG_ERROR, fmt.Sprint(err...))
}

func Waring(war ...interface{}) {
	log.write(LOG_WARING, fmt.Sprint(war...))
}

func Info(info ...interface{}) {
	log.write(LOG_INFO, fmt.Sprint(info...))
}

func Debug(deb ...interface{}) {
	log.write(LOG_DEBUG, fmt.Sprint(deb...))
}

/*
 * 日志执行函数
 */
type mylog struct {
	log      chan string // 日志chan
	dir      string      // 日志存放目录
	file     string      // 日志文件名
	savefile bool        // 是否保存到文件
	level    int         // 日志级别
}

func newMylog() *mylog {
	log := &mylog{}

	log.log = make(chan string, 100)
	log.dir = "/opt/logs"
	log.file = "out"
	log.savefile = false

	go log.run()
	return log
}

func (l *mylog) setDir(dir string) {
	l.dir = dir
}

func (l *mylog) setFile(file string) {
	l.file = file
}

func (l *mylog) setSavefile(b bool) {
	l.savefile = b
}

func (l *mylog) setLevel(level int) {
	l.level = level
}

func (l *mylog) getLevelString(level int) string {
	switch level {
	case LOG_ERROR:
		return "ERROR"
	case LOG_WARING:
		return "WARING"
	case LOG_INFO:
		return "INFO"
	case LOG_DEBUG:
		return "DEBUG"
	}

	return "unknown"
}

func (l *mylog) write(level int, str string) {
	// 判断级别
	if level > l.level {
		return
	}

	// 输出日志
	pc, _, line, _ := runtime.Caller(2)
	p := runtime.FuncForPC(pc)
	t := time.Now()
	str = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d] [%s] %s(%d): %s\n",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(),
		l.getLevelString(level), p.Name(), line, str)
	// 输出到控制台
	if false == l.savefile {
		fmt.Print(str)
		return
	}

	// 输出到文件
	l.log <- str
}

func (l *mylog) run() {
	for {
		str := <-l.log

		// 判断文件夹是否存在
		_, err := os.Stat(l.dir)
		if nil != err {
			os.MkdirAll(l.dir, os.ModePerm)
		}

		// 获取时间
		t := time.Now()
		path := fmt.Sprintf("%s/%s-%04d-%02d-%02d.log", l.dir, l.file,
			t.Year(), t.Month(), t.Day())
		fp, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if nil == err {
			fp.WriteString(str)
			fp.Close()
		}
	}
}
