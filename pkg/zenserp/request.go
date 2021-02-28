package zenserp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	zenserpBaseURL = "https://app.zenserp.com"
	searchPath     = "api/v2/search?q=%s"
)

func (c *client) do(ctx context.Context, method string, endpoint string, body []byte, contentType string) ([]byte, error) {
	p := strings.Split(endpoint, "?")

	rel := &url.URL{Path: p[0]}

	if len(p) > 1 {
		rel.RawQuery = p[1]
	}

	u := c.baseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("apikey", c.apiKey)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Content-Length", strconv.Itoa(len(string(body))))

	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error on Zenserp API %s method (%w)", method, err)
	}
	defer resp.Body.Close()

	rspBody, rspBodyErr := ioutil.ReadAll(resp.Body)
	if rspBodyErr != nil {
		rspBodyErr = fmt.Errorf("error while reading response body (%w)", err)
	}

	// Check whether response status is not 2xx
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if rspBodyErr == nil {
			log.
				WithField("response", string(rspBody)).
				WithField("method", method).
				WithField("endpoint", u.String()).
				WithField("status", resp.StatusCode).
				Warn("Failed Zenserp request")
		}

		err := fmt.Errorf("%d %s", resp.StatusCode, resp.Status)

		return []byte{}, fmt.Errorf("server returned non OK status(%d): %w", resp.StatusCode, err)
	}

	return rspBody, rspBodyErr
}

func (c *client) getJSON(ctx context.Context, endpoint string, result interface{}) error {
	bytes, err := c.do(ctx, "GET", endpoint, []byte{}, "application/json")
	if err != nil {
		return fmt.Errorf("error while calling Zenserp GET JSON API (%w)", err)
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Zenserp response (%w)", err)
	}

	log.WithField("endpoint", endpoint).WithField("response", result).Debug("api call succeeded")
	return nil
}
