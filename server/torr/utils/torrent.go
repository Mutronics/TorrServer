package utils

import (
	"encoding/base32"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"server/settings"
	"golang.org/x/time/rate"
)

func LoadFromUrl(url string) []string {
    var ret []string
    resp, err := http.Get(url)
    if err == nil {
    buf, err := ioutil.ReadAll(resp.Body)
    if err == nil {
        arr := strings.Split(string(buf), "\n")
        for _, s := range arr {
	s = strings.TrimSpace(s)
	if len(s) > 0 {
	    ret = append(ret, s)
	}
        }
    }
    }
    return ret
}

func LoadFromFile() []string {
    var ret []string
    ref := filepath.Join(settings.Path, "retrackers.txt")
    content, err := ioutil.ReadFile(ref)
    if err == nil {
    ret = strings.Split(string(content), "\n")
    }
    return ret
}

func PeerIDRandom(peer string) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return peer + base32.StdEncoding.EncodeToString(randomBytes)[:20-len(peer)]
}

func Limit(i int) *rate.Limiter {
	l := rate.NewLimiter(rate.Inf, 0)
	if i > 0 {
		b := i
		if b < 16*1024 {
			b = 16 * 1024
		}
		l = rate.NewLimiter(rate.Limit(i), b)
	}
	return l
}
