package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type snippetItem struct {
	Scope  string   `json:"scope"`
	Body   []string `json:"body"`
	Prefix string   `json:"prefix"`
}

func parseFile() []byte {

	snippet := make(map[string]snippetItem)

	// walk dir, if found .txt then add to snippet
	cwd, _ := os.Getwd()
	err := filepath.WalkDir(cwd, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(filePath, ".txt") {

			// body
			txt, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			body := strings.Split(string(txt), "\n")

			// scope
			stripCwd := filePath[len(cwd)+1:]
			scope := ""
			idx := strings.Index(stripCwd, string(os.PathSeparator))
			if idx != -1 {
				scope = stripCwd[:idx]
			}
			if scope == "_global" { //file in _global folder is use for all language
				scope = ""
			}

			// prefix
			prefix := stripCwd[:len(stripCwd)-4]

			snippet[prefix] = snippetItem{scope, body, prefix}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// json encode
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(snippet)
	if err != nil {
		log.Fatal(err)
	}
	return buffer.Bytes()
}

func main() {
	json := parseFile()

	//write file
	localConfig := runtime.GOOS
	switch localConfig {
	case "windows":
		localConfig = os.Getenv("APPDATA")
	case "linux":
		localConfig = filepath.Join(os.Getenv("HOME"), ".config")
	case "darwin":
		localConfig = os.Getenv("HOME") + "Library/Application Support"
	}
	target := filepath.Join(localConfig, "Code", "User", "snippets", "snp.code-snippets")
	fmt.Println(string(json))
	fmt.Println("Updated to", target)
	ioutil.WriteFile(target, json, 0644)

}
