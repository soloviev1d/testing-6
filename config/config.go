package config

import (
	"flag"
)

var (
	TestSuiteDir string
)
func FlagsInit() {
	flag.StringVar(&TestSuiteDir, "d", "/home/user/test-suites",
		"Укажите директорию в которой содержатся тест-сьюты")
	flag.StringVar(&TestSuiteDir, "directory", "/home/user/test-suites",
		"Укажите директорию в которой содержатся тест-сьюты")

	flag.Parse()
}
