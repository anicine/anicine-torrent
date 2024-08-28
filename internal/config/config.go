package config

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	MongoDB      string
	Proxy        string
	CertFile     string
	KeyFile      string
	TrackersFile string
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		config  = new(Config)
		scanner = bufio.NewScanner(file)
		comment = regexp.MustCompile(`\s*#.*$|^\s+|\s+$`)
	)

	for scanner.Scan() {
		txt := scanner.Text()
		if !strings.Contains(txt, "=") {
			continue
		}

		line := strings.TrimSpace(comment.ReplaceAllString(txt, ""))
		kv := strings.SplitN(line, "=", 2)
		if len(kv) < 1 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.Trim(strings.Trim(strings.TrimSpace(kv[1]), `"`), `'`)

		switch key {
		case "MONGO_DB_URL":
			if value == "" {
				return nil, errors.New("you need to set the Mongo DB URL")
			}
			config.MongoDB = value
		case "TRACKERS_FILE":
			if value == "" {
				return nil, errors.New("you need to set the path to trackers URLs file")
			}
			config.TrackersFile = value
		case "PROXY_URL":
			config.Proxy = value
		case "CERT_FILE":
			config.CertFile = value
		case "KEY_FILE":
			config.KeyFile = value
		}
	}

	return config, nil
}
