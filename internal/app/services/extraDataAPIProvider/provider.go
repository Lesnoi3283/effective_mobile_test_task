package extraDataAPIProvider

import (
	"encoding/json"
	"fmt"
	"io"
	"musiclib/internal/app/entities"
	"net/http"
	"net/url"
)

const songDataPath = "/info"

// ExtraDataAPIProvider uses API to get extra data.
type ExtraDataAPIProvider struct {
	address string
}

func NewExtraDataAPIProvider(address string) *ExtraDataAPIProvider {
	return &ExtraDataAPIProvider{
		address: address,
	}
}

type extraSongData struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// GetExtraSongData sends a GET request to and ExtraDataAPIProvider.address and returns extra data.
func (p *ExtraDataAPIProvider) GetExtraSongData(song entities.Song) (releaseDate, text, link string, err error) {

	//ask api
	query := url.Values{}
	query.Add("group", song.Group)
	query.Add("song", song.Song)
	res, err := http.Get(p.address + songDataPath + "?" + query.Encode())
	if err != nil {
		return "", "", "", fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("error sending request: status code %d", res.StatusCode)
	}

	//parse data
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("error reading response body: %w", err)
	}
	data := extraSongData{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return "", "", "", fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return data.ReleaseDate, data.Text, data.Link, nil
}
