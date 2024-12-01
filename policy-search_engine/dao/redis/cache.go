package redis

import (
	"PolicySearchEngine/config"
	"github.com/gocolly/colly"
	"hash/fnv"
	"time"
)

func SetRedisStorage(c *colly.Collector, prefix string, urls []string) {

	storage := &Storage{
		Address:     config.V.GetString("redis.addr"),
		Password:    config.V.GetString("redis.password"),
		DB:          config.V.GetInt("redis.db"),
		Prefix:      prefix,
		ExceptionID: UrlToRequestID(urls),
		Expires:     time.Hour,
	}

	err := c.SetStorage(storage)
	if err != nil {
		panic(err)
	}

}

func UrlToRequestID(urls []string) (requestID []uint64) {
	for _, u := range urls {
		h := fnv.New64a()
		_, _ = h.Write([]byte(u))
		requestID = append(requestID, h.Sum64())
	}
	return
}
