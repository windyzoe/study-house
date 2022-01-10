package util

import "flag"

var FLAG_ENV = flag.String("env", "local", "enviroment")

func Init() {
	flag.Parse()
}
