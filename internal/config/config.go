package config

import (
	"os"
	"path/filepath"
)

var DataDir = "./data"

func SubDir(sub string) string {
	dir := filepath.Join(DataDir, sub)
	os.MkdirAll(dir, 0755)
	return dir
}

func KVDir() string     { return SubDir("kv") }
func SQLDir() string    { return SubDir("sql") }
func GraphDir() string  { return SubDir("graph") }
func SearchDir() string { return SubDir("search") }
func TSDir() string     { return SubDir("ts") }
