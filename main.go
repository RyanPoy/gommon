package main

import (
	"bufio"
	"gommon/ip"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// 你的应用逻辑
	runAppV4()
	//runAppV6()
}

func runAppV4() {
	fpath := "./ip/test_data/mgiplib-std.txt.latest"
	v4s, _ := ip.NewV4Table(fpath)
	ips, _ := LoadIP(fpath)
	for {
		for _, v := range ips {
			v4s.Search(v)
		}
	}
}

func runAppV6() {
	fpath := "./ip/test_data/mgiplib-v6-std.txt.latest"
	v6s, _ := ip.NewV6Table(fpath)
	ips, _ := LoadIP(fpath)
	for {
		for _, v := range ips {
			v6s.Search(v)
		}
	}
}

func LoadIP(fpath string) ([]string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]string, 0)
	for scanner := bufio.NewScanner(f); scanner.Scan(); {
		line := scanner.Text()
		if line[0] == '#' || len(line) == 0 || !strings.Contains(line, "|") {
			continue
		}
		lines = append(lines, strings.Split(line, "|")[0])
	}
	return lines, nil
}
