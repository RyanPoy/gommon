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
	runApp()
}

func runApp() {
	fpath4 := "./ip/test_data/mgiplib-std.txt.latest"
	v4s, _ := ip.NewTable(fpath4)
	ip4s, _ := LoadIP(fpath4)

	fpath6 := "./ip/test_data/mgiplib-v6-std.txt.latest"
	v6s, _ := ip.NewTable(fpath6)
	ip6s, _ := LoadIP(fpath6)

	go func() {
		for {
			for _, v := range ip4s {
				v4s.Search(v)
			}
		}
	}()

	go func() {
		for {
			for _, v := range ip6s {
				v6s.Search(v)
			}
		}
	}()
	select {}
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
