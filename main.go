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
	runApp1()
}

func runApp1() {
	fpath := "./ip/test_data/mgiplib-v6-std.txt.latest"
	v6s, _ := ip.NewV6s(fpath)
	ips, _ := LoadV6(fpath)
	for {
		for _, v := range ips {
			v6s.Search(v)
		}
	}
}

func runApp2() {
	fpath := "./ip/test_data/mgiplib-v6-std.txt.latest"
	v6s, _ := ip.NewV6Table(fpath)
	ips, _ := LoadV6(fpath)
	for {
		for _, v := range ips {
			v6s.Search(v)
		}
	}
}

func LoadV6(fpath string) ([]string, error) {
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
