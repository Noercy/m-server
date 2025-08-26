package anilist

import (
	"fmt"
	"bytes"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
)

/* 
	Clean description function to clean up the HTML content
	It removes unnecessary HTML tags, replaces line breaks with newlines,
	and unescapes HTML entities.

	find a better way to clean up the description later.
*/
func CleanDescription(description string) string {

	description = strings.ReplaceAll(description, "<br>", "\n")
	description = strings.ReplaceAll(description, "<br/>", "\n")
	description = strings.ReplaceAll(description, "<br />", "\n")
	description = strings.ReplaceAll(description, "<p>", "")
	description = strings.ReplaceAll(description, "</p>", "\n")
	description = strings.ReplaceAll(description, "<i>", "")
	description = strings.ReplaceAll(description, "</i>", "")
	description = strings.ReplaceAll(description, "<b>", "")
	description = strings.ReplaceAll(description, "</b>", "")
	
	tagRemove := regexp.MustCompile(`(?s)<.*?>`)
	description = tagRemove.ReplaceAllString(description, "")

	// I dont know what this does 
	description = html.UnescapeString(description)
	newlinesRemove := regexp.MustCompile(`\n{3,}`)
	description = newlinesRemove.ReplaceAllString(description, "\n\n")
	description = strings.ReplaceAll(description, "\u00A0", " ") // non-breaking space
	return strings.TrimSpace(description)
}

func SearchMangaByTitle(title string) (int, error) {
	query := `
	query ($search: String) {
	  Page(page: 1, perPage: 1) {
	    media(search: $search, type: MANGA) {
	      id
	      title {
	        romaji
	        english
	        native
	      }
	    }
	  }
	}`

	variables := map[string]interface{}{
		"search": title,
	}

	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Struct to capture just the ID from the search
	var result struct {
		Data struct {
			Page struct {
				Media []struct {
					ID int `json:"id"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if len(result.Data.Page.Media) == 0 {
		return 0, fmt.Errorf("no results for %q", title)
	}

	return result.Data.Page.Media[0].ID, nil

}

func GetMangaByID(id int) (Media, error) {
	query := `
	query ($id: Int) {
	  Media (id: $id, type: MANGA) {
	    id
	    title {
	      romaji
	      english
	      native
	    }
		coverImage {
      		large
    	}
		description(asHtml: false)
		genres
		startDate {
			day
			month
			year
		}
    	status	
		volumes
    	chapters   
	  }
	}`

	variables := map[string]interface{}{
		"id": id,
	}

	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Media{}, err
	}

	url := "https://graphql.anilist.co"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Media{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Media{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Media{}, err
	}

	var result AniListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return Media{}, err
	}

	m := result.Data.Media
	m.Description = CleanDescription(m.Description)

	return m, nil
}