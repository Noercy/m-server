package anilist 

type MediaTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type CoverImage struct {
	Large string `json:"large"`
}

type StartDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type Media struct {
	ID    int        `json:"id"`
	Title MediaTitle `json:"title"`
	CoverImage  CoverImage `json:"coverImage"`
	Description string     `json:"description"`
	Genres      []string   `json:"genres"`
	StartDate   StartDate  `json:"startDate"`
	Status      string     `json:"status"`
	Volumes     int        `json:"volumes"`
	Chapters    int        `json:"chapters"`
	Publications string `json:"publications"`
	Publishers   string `json:"publishers"`
}

type Data struct {
	Media Media `json:"Media"`
}

type AniListResponse struct {
	Data Data `json:"data"`
}