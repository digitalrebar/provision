package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	pbuf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	pat := string(pbuf)
	res := []*regexp.Regexp{}
	for _, line := range strings.Split(pat[4:], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		re, err := regexp.Compile(`^[\s]*(` + line + `)[\s]*$`)
		if err != nil {
			log.Fatalf("Regex compile error: %v", err)
		} else {
			res = append(res, re)
		}
	}
	so, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	soLines := strings.Split(string(so), "\n")
	soIdx := 0
	reMatches := []bool{}
	for _, re := range res {
		for {
			matched := false
			if re.MatchString(soLines[soIdx]) {
				reMatches = append(reMatches, true)
				log.Printf("Matched %s at %d", re.String(), soIdx)
				matched = true
			} else {
				log.Printf("Failed match %s at %d", re.String(), soIdx)
			}
			soIdx++
			if soIdx == len(soLines) || matched {
				break
			}
		}
		if soIdx == len(soLines) {
			break
		}
	}
	if len(reMatches) != len(res) {
		os.Exit(1)
	}
	log.Printf("Matched overall")
}
