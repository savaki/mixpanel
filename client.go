package mixpanel

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	get func(string, url.Values) (io.ReadCloser, error)
}

const ONE_HOUR = 3600

func getFn(client *http.Client, apiKey, apiSecret string) func(string, url.Values) (io.ReadCloser, error) {
	return func(path string, params url.Values) (io.ReadCloser, error) {
		params.Add("api_key", apiKey)
		params.Add("expire", strconv.Itoa(int(time.Now().Unix() + ONE_HOUR)))

		// sort our keys
		keys := []string{}
		for key, _ := range params {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// construct the data to sign
		argsConcat := []string{}
		for _, key := range keys {
			argsConcat = append(argsConcat, fmt.Sprintf("%s=%s", key, params.Get(key)))
		}

		// generate our signature
		hash := md5.New()
		hash.Write([]byte(strings.Join(argsConcat, "") + apiSecret))
		signature := hex.EncodeToString(hash.Sum(nil))

		// add the signature to the query params
		params.Add("sig", signature)

		uri := fmt.Sprintf("https://data.mixpanel.com%s?%s", path, params.Encode())
		request, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return nil, err
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}

		return response.Body, nil
	}
}

func New(apiKey, apiSecret string) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}

	return &Client{
		get: getFn(client, apiKey, apiSecret),
	}
}
