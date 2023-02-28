package entities

import (
	"math/rand"
	"time"
)

type URL struct {
	UserID        string
	ShortURL      string
	FullURL       string
	CorrelationID string
}

func (u *URL) GenerateRandomString(n int) {
	var res string
	line := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomSrc := rand.NewSource(time.Now().UnixMicro())
	rnd := rand.New(randomSrc)
	for i := 0; i < n; i++ {
		res += string(line[rnd.Intn(len(line))])
	}
	u.ShortURL = res

}
