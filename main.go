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
	fpaths := []string{"./ip/test_data/mgiplib-std.txt.latest", "./ip/test_data/mgiplib-v6-std.txt.latest"}
	table, _ := ip.NewIPTable(fpaths...)
	ips, _ := loadIP(fpaths...)

	for {
		for _, ipstr := range ips {
			table.Search(ipstr)
		}
	}
}

func loadIP(fpaths ...string) ([]string, error) {
	lines := make([]string, 0)

	for _, fpath := range fpaths {
		f, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		for scanner := bufio.NewScanner(f); scanner.Scan(); {
			line := scanner.Text()
			if line[0] == '#' || len(line) == 0 || !strings.Contains(line, "|") {
				continue
			}
			lines = append(lines, strings.Split(line, "|")[0])
		}
	}

	return lines, nil
}
