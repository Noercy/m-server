package models

type OldSeriesMetaData struct {
	ID           int      `json:"id"`
	TitleRomaji  string   `json:"title_romaji"`
	TitleEnglish string   `json:"title_english"`
	TitleNative  string   `json:"title_native"`
	CoverImage   string   `json:"cover_image"`
	Description  string   `json:"description"`
	Genres       []string `json:"genres"`
	StartDate    struct {
		Day   int `json:"day"`
		Month int `json:"month"`
		Year  int `json:"year"`
	} `json:"start_date"`
	Status       string `json:"status"`
	Volumes      int    `json:"volumes"`
	Chapters     int    `json:"chapters"`
	Publications string `json:"publications"`
	Publishers   string `json:"publishers"`
}


type Volume struct {
	ID        int    `json:"ID"`
	Number    int    `json:"Number"`
	NumImages int    `json:"Num_images"`
	Title     string `json:"Title"`
	Path      string `json:"Path"`
	Cover     string `json:"Cover"`
	CreatedAt string `json:"Created_at"`
}

type DbMetadata struct {
	TitleRomaji   string `json:"Title_romaji"`
	TitleEnglish  string `json:"Title_english"`
	TitleNative   string `json:"Title_native"`
	Description   string `json:"Description"`
	ReleaseDate   string `json:"Release_date"`
	Publisher     string `json:"Publisher"`
	Publication   string `json:"Publication"`
	TotalVol      int    `json:"Total_vol"`
	TotalCh       int    `json:"Total_ch"`
	ReleaseStatus string `json:"Release_status"`
}

type Series struct {
	ID        int    `json:"ID"`
	Title     string `json:"Title"`
	Path      string `json:"Path"`
	Cover     string `json:"Cover"`
	NumVol    int    `json:"Num_vol"`
	NumImages int    `json:"Num_images"`
	CreatedAt string `json:"Created_at"`

	Metadata DbMetadata `json:"Metadata"`
	Genres   []string   `json:"Genres"`
	Tags     []string   `json:"Tags"`
	Volumes  []Volume   `json:"Volumes"`
}
