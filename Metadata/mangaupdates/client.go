package mangaupdates
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SearchResult struct {
	Results []struct {
		Record struct {
            SeriesID int `json:"series_id"`
		} `json:"record"`
	} `json:"results"`
}

type SeriesData struct {
	Publications []struct {
		PublicationName string `json:"publication_name"`
		PublisherName   string `json:"publisher_name"`
	} `json:"publications"`
}

func SearchSeriesID(title string) (int, error) {
	url := "https://api.mangaupdates.com/v1/series/search"

	payload := map[string]string{"search": title}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result SearchResult 
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if len(result.Results) == 0 {
		return 0, fmt.Errorf("no results for %q", title)
	}

	return result.Results[0].Record.SeriesID, nil
}


func GetSeriesPublications(seriesID int) (string, string, error) {
	url := fmt.Sprintf("https://api.mangaupdates.com/v1/series/%d", seriesID)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data SeriesData
	if err := json.Unmarshal(body, &data); err != nil {
		return "", "", err
	}

	var pubs []string
	var publishers []string

	for _, pub := range data.Publications {
		if pub.PublicationName != "" {
			pubs = append(pubs, pub.PublicationName)
		}
		if pub.PublisherName != "" {
			publishers = append(publishers, pub.PublisherName)
		}
	}

	return strings.Join(pubs, ", "), strings.Join(publishers, ", "), nil
}