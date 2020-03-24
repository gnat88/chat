package broker

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var Worlds []string
func LoadWorlds(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	buff := bufio.NewReader(f)
	for {
		line, _, err := buff.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		Worlds = append(Worlds, string(line))
	}
}

func Replace(s string) string {
	for _, v := range Worlds {
		if strings.Contains(s, v) {
			s = strings.ReplaceAll(s, v, "***")
		}
	}
	return s
}