package helper

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ThinkGoPath(args ...string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic(errors.New("Can not load ThinkGo path info"))
	}

	dir := filepath.Dir(filepath.Dir(file)) + "/"

	if 1 == len(args) {
		dir = dir + strings.TrimLeft(args[0], "/")
	}

	return dir
}

func WorkPath(args ...string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if 1 == len(args) {
		dir = filepath.Join(dir, strings.TrimLeft(args[0], "/"))
	}

	return dir
}

func AppPath(args ...string) string {

	dir := WorkPath("app")

	if 1 == len(args) {
		dir = filepath.Join(dir, strings.TrimLeft(args[0], "/"))
	}

	return dir
}

func ConfigPath(args ...string) string {
	dir := WorkPath("config")

	if 1 == len(args) {
		dir = filepath.Join(dir, strings.TrimLeft(args[0], "/"))
	}

	return dir
}
