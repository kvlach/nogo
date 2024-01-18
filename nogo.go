package nogo

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type set map[string]struct{}

type NoGo struct {
	fakenews bool
	gambling bool
	porn     bool
	social   bool

	hostnames set
}

func Init() *NoGo {
	return &NoGo{}
}

func (n *NoGo) Fakenews() *NoGo {
	n.fakenews = true
	return n
}

func (n *NoGo) Gambling() *NoGo {
	n.gambling = true
	return n
}

func (n *NoGo) Porn() *NoGo {
	n.porn = true
	return n
}

func (n *NoGo) Social() *NoGo {
	n.social = true
	return n
}

func (n *NoGo) generateDownloadURL() string {
	url := "https://raw.githubusercontent.com/StevenBlack/hosts/master"

	if !(n.fakenews || n.gambling || n.porn || n.social) {
		return url + "/hosts"
	}

	categories := []string{}

	if n.fakenews {
		categories = append(categories, "fakenews")
	}

	if n.gambling {
		categories = append(categories, "gambling")
	}

	if n.porn {
		categories = append(categories, "porn")
	}

	if n.social {
		categories = append(categories, "social")
	}

	return url + "/alternates/" + strings.Join(categories, "-") + "/hosts"
}

func (n *NoGo) Download() (*NoGo, error) {
	n.hostnames = set{}

	resp, err := http.Get(n.generateDownloadURL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	reLine := regexp.MustCompile(`^0\.0\.0\.0 [0-9a-z\.\-]+\.[a-z]+`)
	for scanner.Scan() {
		line := scanner.Text()

		if !reLine.MatchString(line) {
			continue
		}

		// TODO: probably change this to just index first space and slice string
		fields := strings.Split(line, " ")
		hostname := fields[1]
		n.hostnames[hostname] = struct{}{}
	}

	return n, scanner.Err()
}

// Check if a hostname is blacklisted
func (n *NoGo) Safe(hostname string) bool {
	// if it's in the set of urls then that means that it's not safe
	if _, ok := n.hostnames[hostname]; ok {
		return false
	}
	return true
}
