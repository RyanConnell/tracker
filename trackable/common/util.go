package common

import (
	"io/ioutil"
	"net/http"
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
