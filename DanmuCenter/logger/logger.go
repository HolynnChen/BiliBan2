package logger

import "log"

var silentLog = false

func SetSilent(silent bool) {
	silentLog = true
}

func Println(v ...interface{}) {
	if silentLog {
		return
	}
	log.Println(v...)
}

func Printf(format string, v ...interface{}) {
	if silentLog {
		return
	}
	log.Printf(format, v...)
}
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

func Panic(v ...interface{}) {
	log.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	log.Panicln(v...)
}
