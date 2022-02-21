package env

import (
	"bufio"
	"firebase-sso/helpers"
	"fmt"
	"log"
	"os"
	"strings"
)

var keys = make(map[string]string)

type Environment struct {
	EnvPath string
}

func (e *Environment) LoadEnv() {

	if e.EnvPath == "" {
		e.EnvPath = ".env"
	}
	_, err := os.Stat(e.EnvPath)
	if err != nil {
		_, err = os.Stat(helpers.BasePath() + "/" + e.EnvPath)
		if err == nil {
			e.EnvPath = helpers.BasePath() + "/" + e.EnvPath
		} else {
			log.Fatal(err)
		}
	}

	file, err := os.Open(e.EnvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line != "" && string(line[0]) != "#" {
			envVar := strings.Split(line, "=")
			if os.Getenv(strings.TrimSpace(envVar[0])) != "" {
				keys[strings.TrimSpace(envVar[0])] = os.Getenv(strings.TrimSpace(envVar[0]))
			} else {
				keys[strings.TrimSpace(envVar[0])] = strings.TrimSpace(envVar[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("ENV LOADED")
	fmt.Println("----------")
	for key, value := range keys {
		fmt.Println(key + "  :  " + value)
	}
	fmt.Println("----------")

}

func Get(key string) string {
	return keys[key]
}
