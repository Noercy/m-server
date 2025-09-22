package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	//"html/template"
	"log"
	"net/http"

	"project/metadata/anilist"
	"project/metadata/mangaupdates"
	"project/metadata/models"

	//"project/scanner"

	_ "github.com/mattn/go-sqlite3"
	//"github.com/NYTimes/gziphandler"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
)

type smallSeries struct {
    ID        int    
    Title     string 
    Cover     string 
}

type config struct {
	RootPath		string	`json:"rootPath"`
	ThumbnailPath	string	`json:"thumbnailPath"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetFullMangaInfo(title string) (models.OldSeriesMetaData, error) {
	var fullInfo models.OldSeriesMetaData

	aniId, err := anilist.SearchMangaByTitle(title)
	if err != nil {	
		return fullInfo, err
	}

	aniData, err := anilist.GetMangaByID(aniId)
	if err != nil {
		return fullInfo, err
	}
	
	fullInfo.ID = aniData.ID
	fullInfo.TitleRomaji = aniData.Title.Romaji
	fullInfo.TitleEnglish = aniData.Title.English
	fullInfo.TitleNative = aniData.Title.Native
	fullInfo.CoverImage = aniData.CoverImage.Large
	fullInfo.Description = aniData.Description
	fullInfo.Genres = aniData.Genres
	fullInfo.StartDate.Year = aniData.StartDate.Year
	fullInfo.StartDate.Month = aniData.StartDate.Month
	fullInfo.StartDate.Day = aniData.StartDate.Day
	fullInfo.Status = aniData.Status
	fullInfo.Volumes = aniData.Volumes
	fullInfo.Chapters = aniData.Chapters

	muID, err := mangaupdates.SearchSeriesID(title)
	if err == nil {
		fmt.Printf("Found MangaUpdates ID: %d\n", muID)
		publications, publishers, err := mangaupdates.GetSeriesPublications(muID)
		if err == nil {
			fullInfo.Publications = publications
			fullInfo.Publishers = publishers
		}
	}

	return fullInfo, nil
}

func GetSerieDataByID(id int, db *sql.DB) (*models.Series, error) {
	rows, err := db.Query(`
    SELECT 
		s.id, s.title, s.path, s.cover_image, s.num_vol, s.num_images, s.created_at,

		COALESCE(m.title_romaji, '') AS title_romaji,
		COALESCE(m.title_english, '') AS title_english,
		COALESCE(m.title_native, '') AS title_native,
		COALESCE(m.description, '') AS description,
		COALESCE(m.release_date, '') AS release_date,
		COALESCE(m.publisher, '') AS publisher,
		COALESCE(m.publication, '') AS publication,
		COALESCE(m.total_vol, 0) AS total_vol,
		COALESCE(m.total_ch, 0) AS total_ch,
		COALESCE(m.release_status, '') AS release_status,
		
		g.name AS genre,
		t.name AS tag,

		v.id AS volume_id, v.number, v.num_images, v.title, v.path, v.cover, v.created_at
        FROM series s
        LEFT JOIN series_metadata m ON m.series_id = s.id
        LEFT JOIN series_genres sg ON sg.series_id = s.id
        LEFT JOIN genres g ON g.id = sg.genre_id
        LEFT JOIN series_tags st ON st.series_id = s.id
        LEFT JOIN tags t ON t.id = st.tag_id
        LEFT JOIN volumes v ON v.series_id = s.id
        WHERE s.id = ?;
    `, id)
	if err != nil {
		log.Println("Db query failed")
		return nil, err;
	}
	defer rows.Close()

	var series models.Series
	genres := make(map[string]bool)
	tags := make(map[string]bool)
	var volumeList []models.Volume

	rowCount := 0
	for rows.Next() {
		rowCount++
        var (
            genre, tag sql.NullString
            vol models.Volume
        )

        err := rows.Scan(
            &series.ID, &series.Title, &series.Path, &series.Cover,
            &series.NumVol, &series.NumImages, &series.CreatedAt,
            &series.Metadata.TitleRomaji, &series.Metadata.TitleEnglish, &series.Metadata.TitleNative,
            &series.Metadata.Description, &series.Metadata.ReleaseDate, &series.Metadata.Publisher,
            &series.Metadata.Publication, &series.Metadata.TotalVol, &series.Metadata.TotalCh,
            &series.Metadata.ReleaseStatus,
            &genre, &tag,
            &vol.ID, &vol.Number, &vol.NumImages, &vol.Title, &vol.Path, &vol.Cover, &vol.CreatedAt,
        )
        if err != nil {
			log.Printf("❌ Row scan failed: %v", err)
            return nil, err
        }
		log.Printf("✅ Row %d scanned: genre=%v, tag=%v, volID=%d", rowCount, genre, tag, vol.ID)


        if genre.Valid {
            genres[genre.String] = true
        }
        if tag.Valid {
            tags[tag.String] = true
        }
        if vol.ID != 0 { // only add if exists
        	volumeList = append(volumeList, vol)
        }
    }

	sort.Slice(volumeList, func(i, j int) bool {
    	return volumeList[i].Number < volumeList[j].Number
	})
	series.Volumes = volumeList

	if err := rows.Err(); err != nil {
        log.Printf("❌ Row iteration error: %v", err)
        return nil, err
    }

    if rowCount == 0 {
        log.Printf("⚠️ No series found with id=%d", id)
        return nil, sql.ErrNoRows
    }
	
	for g := range genres {
        series.Genres = append(series.Genres, g)
    }
    for t := range tags {
        series.Tags = append(series.Tags, t)
    }

	log.Println("Invidual DB query done")
    return &series, nil
}

func printMetaData() {
	metaData, err := GetFullMangaInfo("My Hero Academia")
	if err != nil {
		panic(err)
	} 
	
	fmt.Println("Title:", metaData.TitleEnglish)
	fmt.Println("Description:", metaData.Description)
	fmt.Println("Genres:", metaData.Genres)
	fmt.Println("Status:", metaData.Status)
	fmt.Println("Published:", metaData.StartDate.Day, metaData.StartDate.Month, metaData.StartDate.Year)
	fmt.Println("Volumes:", metaData.Volumes)
	fmt.Println("Chapters:", metaData.Chapters)
	fmt.Println("Cover Image:", metaData.CoverImage)
	fmt.Println("Publishers:", metaData.Publishers)
	fmt.Println("Publications:", metaData.Publications)
}

func GenerateUserDBSchema(db *sql.DB) {
	schema := 
	`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username STRING NOT NULL UNIQUE,
		password STRING NOT NULL
	);`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO want to create a folder called thumbnails whereever the cache thumbnail folder is choosen
func serverSetup() {
	var serverConfig config 
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter desired Rootpath: ")
	rootPath, _ := reader.ReadString('\n')
	serverConfig.RootPath = strings.TrimSpace(rootPath)

	fmt.Println("Enter desired thumbnail location: ")
	thumbnailPath, _ := reader.ReadString('\n')
	serverConfig.ThumbnailPath = strings.TrimSpace(thumbnailPath)	

	data, err := json.MarshalIndent(serverConfig, "", " ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("serverConfig.json", data, 0644);
	if err != nil {
		panic(err)
	}
}

func configExists() bool {
	_, err := os.Stat("serverConfig.json")
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func getConfigPaths() (config, error) {
	var serverConfig config

	data, err := os.ReadFile("serverConfig.json")
	if err != nil {
		return serverConfig, err
	}

	err = json.Unmarshal(data, &serverConfig)
	if err != nil {
		return serverConfig, err
	}

	return serverConfig, nil
}

func getAllSeries(db *sql.DB) []smallSeries {
	rows, err := db.Query("SELECT id, title, cover_image FROM series ORDER BY created_at DESC")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var series []smallSeries
    for rows.Next() {
        var s smallSeries
        if err := rows.Scan(&s.ID, &s.Title, &s.Cover); err != nil {
            log.Println("Error scanning series:", err)
            continue
        }
        series = append(series, s)
    }
    return series
}

// test
func testHandler() gin.HandlerFunc {
	 return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hey",
		})
	}
}

// register user
func registerUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		GenerateUserDBSchema(db)
		var newUser User
		if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		insertNewUser := `INSERT INTO users (username, password) VALUES (?, ?);`
		db.Exec(insertNewUser, newUser.Username, newUser.Password)

		c.JSON(http.StatusOK, gin.H{
			"message": "User Regged successful",
			"username": newUser.Username,
			"password": newUser.Password,
		})
	}
}

// login user
func loginUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginReq User
		if err := c.ShouldBindBodyWithJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		var storedPass string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", loginReq.Username).Scan(&storedPass)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		if storedPass != loginReq.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}

		session := sessions.Default(c)
		session.Set("user", loginReq.Username)
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "login successful"})
	}
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(401, gin.H{"error:": "unauthorized"})
			c.Abort()
			return 
		}
		c.Set("user", user)
		c.Next()
	}
}

func currentUserHandler(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, gin.H{"user": user})
}

func getAllSeriesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		series := getAllSeries(db)
		c.JSON(http.StatusOK, series)
	}
}

// get a single series and the volumes associated with it
func getSerieHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		log.Printf("id: %d", id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		s, err := GetSerieDataByID(id, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if s.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "series not found"})
			return
		}
		
    	c.JSON(http.StatusOK, s)
	}
}

// get the images of a single volume from a series for the reader
func getReaderVolPages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		seriesId := c.Param("id")
		volumeId := c.Param("vId")

		var volumePath, seriesTitle string 
		 err := db.QueryRow(`
			SELECT s.title, v.path 
			FROM volumes v
			JOIN series s ON s.id = v.series_id
			WHERE v.id = ? AND v.series_id = ?`,
			volumeId, seriesId).Scan(&seriesTitle, &volumePath)
		if err != nil {
			log.Printf("❌ Volume lookup failed: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Volume not found for this series"})
			return
		}

		files, err := os.ReadDir(volumePath)
		if err != nil {
			log.Printf("❌ Failed to read volume folder: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read volume directory"})
			return
		}

		var images []string
		for _, f := range files {
			if !f.IsDir() {
				ext := strings.ToLower(filepath.Ext(f.Name()))
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
	            	relPath := filepath.Join(seriesTitle, filepath.Base(volumePath), f.Name())
					relPath = strings.ReplaceAll(relPath, "\\", "/") 
					images = append(images, "/pages/"+relPath)
				}
			}
		}

		// sort
		sort.Slice(images, func(i, j int) bool {
			return images[i] < images[j]
		})
	
		c.JSON(http.StatusOK, gin.H{
			"series_id": seriesId,
			"volume_id": volumeId,
			"images":    images,
		})
	}
}

func main() {
	//scanner.FullFolderScan()
	if !configExists() {
		serverSetup()
	}
	
	configData, err := getConfigPaths()
	if err != nil {log.Fatal(err)}

	db, err := sql.Open("sqlite3", "./manga.db")
    if err != nil {log.Fatal(err)}
    defer db.Close()

	router := gin.Default();

	store := cookie.NewStore([]byte("sessionKeys"))
	router.Use(sessions.Sessions("TheSession", store))

	router.Static("/static", "./static")
	router.Static("/thumbnails", configData.ThumbnailPath)
	pages := router.Group("/pages")
	pages.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=2000")
		c.Next()
	})
	pages.Static("/", configData.RootPath)
	router.GET("/me", authRequired(), currentUserHandler)

	router.POST("/api/register", registerUserHandler(db))

	router.POST("/api/login", loginUserHandler(db))

	router.GET("/api/allseries", getAllSeriesHandler(db))
	
	router.GET("/api/series/:id", getSerieHandler(db))

	router.GET("/api/series/:id/reader/:vId", getReaderVolPages(db))

	router.GET("/api/test", testHandler())
	router.Run("localhost:8080")
/*
	tmpl := template.Must(template.ParseFiles("Templates/index.html", "Templates/series_page.html"))
	
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("E:/$otaku/Mangaserver"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", http.FileServer(http.Dir("E:/$otaku/mangaserver_cache/thumbnails"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, series)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Scan started"}`))
	})

	http.HandleFunc("/series", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
        	http.Error(w, "Missing series ID", http.StatusBadRequest)
        	return
    	}

		fmt.Println(id)
		row := db.QueryRow("SELECT id, title, cover_image FROM series WHERE id = ?", id)
    	var s Series
		if err := row.Scan(&s.ID, &s.Title, &s.Cover); err != nil {
			http.Error(w, "Not Found", 404)
			return
		}

		fmt.Println(s)
		tmpl.ExecuteTemplate(w, "series_page.html", s)
	})

	compressedHandler := gziphandler.GzipHandler(http.DefaultServeMux)
	log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe("localhost:8080", compressedHandler))
	*/
}