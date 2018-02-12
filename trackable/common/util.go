package common

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Convert reader to bytes
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func ParseString(str string) string {
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	str = strings.Trim(str, "\t")
	str = strings.Trim(str, " ")
	str = strings.Replace(str, "  ", " ", 100)
	str = strings.Replace(str, "\t\t", "\t", 100)
	str = strings.Replace(str, "\n", "", 100)
	return str
}

// Creates a map[string]string from a file with lines such as: "port: 80"
func LoadSettings(filename string) map[string]string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	settings := map[string]string{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		s := strings.SplitN(line, ":", 2)
		if len(s) == 2 {
			settings[strings.TrimSpace(s[0])] = strings.TrimSpace(s[1])
		}
	}
	return settings
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
