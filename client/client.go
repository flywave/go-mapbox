package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flywave/go-mapbox/style"
	"github.com/flywave/go-mapbox/tilejson"
	"github.com/pkg/errors"
)

type Client struct {
	username string
	token    string
	baseURL  url.URL

	HttpClient *http.Client
}

func NewClient(username, token string) *Client {
	baseURL, _ := url.Parse("https://api.mapbox.com")

	httpClient := http.Client{
		Timeout: 15 * time.Second,
	}

	return &Client{
		username:   username,
		token:      token,
		baseURL:    *baseURL,
		HttpClient: &httpClient,
	}
}

func (c *Client) do(method string, url url.URL, body io.Reader, value interface{}) (*http.Response, error) {
	url = c.addAuthentication(url)
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "requesting")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(value)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json")
	}

	return resp, nil
}

func (c *Client) addAuthentication(url url.URL) url.URL {
	q := url.Query()
	q.Add("access_token", c.token)
	url.RawQuery = q.Encode()
	return url
}

func (c *Client) nextPageURL(header http.Header) *url.URL {
	nextRegex := regexp.MustCompile("<(.*)>")

	link := header.Get("Link")
	if strings.Contains(link, "next") {
		requestURL, err := url.Parse(nextRegex.FindStringSubmatch(link)[1])
		if err != nil {
			return nil
		}
		return requestURL
	} else {
		return nil
	}
}

func (c *Client) ListStyles(draft bool) ([]style.ListStyle, error) {

	url := c.baseURL
	url.Path = path.Join(url.Path, "styles/v1/", c.username)
	q := url.Query()
	if draft {
		q.Add("draft", "true")
	} else {
		q.Add("draft", "false")
	}
	q.Add("fresh", "true")
	url.RawQuery = q.Encode()

	var allStyles []style.ListStyle

	requestURL := &url

	for requestURL != nil {

		var styles []style.ListStyle

		resp, err := c.do("GET", *requestURL, nil, &styles)
		if err != nil {
			return nil, errors.Wrap(err, "making request")
		}
		requestURL = c.nextPageURL(resp.Header)

		allStyles = append(allStyles, styles...)
	}

	return allStyles, nil
}

func (c *Client) GetStyle(id string, draft bool) (style.Style, error) {

	url := c.baseURL
	url.Path = path.Join(url.Path, "styles/v1/", c.username, id)
	if draft {
		url.Path = path.Join(url.Path, "draft")
	}

	var s style.Style

	_, err := c.do("GET", url, nil, &s)
	if err != nil {
		return style.Style{}, errors.Wrap(err, "making request")
	}

	return s, nil
}

func (c *Client) ListTilesets(params tilejson.ListTilesetsParams) ([]tilejson.Tileset, error) {
	url := c.baseURL
	url.Path = path.Join(url.Path, "tilesets/v1/", c.username)
	q := url.Query()
	if params.Type != nil {
		q.Add("type", string(*params.Type))
	}
	if params.Visibility != nil {
		q.Add("visibility", string(*params.Visibility))
	}
	if params.SortBy != nil {
		q.Add("sortby", string(*params.SortBy))
	}
	if params.Limit != nil {
		q.Add("limit", strconv.Itoa(*params.Limit))
	}
	url.RawQuery = q.Encode()

	var allTilesets []tilejson.Tileset

	requestURL := &url

	for requestURL != nil {

		var tilesets []tilejson.Tileset

		resp, err := c.do("GET", *requestURL, nil, &tilesets)
		if err != nil {
			return nil, errors.Wrap(err, "making request")
		}
		requestURL = c.nextPageURL(resp.Header)

		allTilesets = append(allTilesets, tilesets...)
	}

	return allTilesets, nil
}

func (c *Client) GetTileJSON(tilesetIDs ...string) (tilejson.TileJSON, error) {
	ids := strings.Join(tilesetIDs, ",")

	url := c.baseURL
	url.Path = path.Join(url.Path, "v4", ids) + ".json"

	var metadata tilejson.TileJSON

	_, err := c.do("GET", url, nil, &metadata)
	if err != nil {
		return tilejson.TileJSON{}, err
	}

	return metadata, nil
}
