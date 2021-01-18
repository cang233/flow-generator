package log

import "log"

var logger log.Logger

func init() {
	logger.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	logger.SetPrefix("flowgen")
}

func Print(v ...interface{}) {
	logger.Print(v)
}
func Println(v ...interface{}) {
	logger.Println(v)
}

func Panicf(formatString string, v ...interface{}) {
	logger.Panicf(formatString, v)
}

func Panic(v ...interface{}) {
	logger.Panic(v)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v)
}

func Fatalf(formatString string, v ...interface{}) {
	logger.Fatalf(formatString, v)
}
