package scanner

import (
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/chai2010/webp"
	"github.com/maruel/natural"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nfnt/resize"
)

func initDB(db *sql.DB) {
    schemas := []string{ 
        `CREATE TABLE IF NOT EXISTS series (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            path TEXT NOT NULL UNIQUE,
            cover_image TEXT,
            num_vol INTEGER NOT NULL,
            num_images INTEGER NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,

       `CREATE TABLE IF NOT EXISTS volumes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            series_id INTEGER NOT NULL,
            number INTEGER NOT NULL,
            num_images INTEGER NOT NULL,
            title TEXT,
            path TEXT NOT NULL,
            cover TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE
        );`,

        `CREATE TABLE IF NOT EXISTS chapters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            volume_id INTEGER NOT NULL,
            number INTEGER NOT NULL,
            title TEXT,
            path TEXT NOT NULL,
            cover TEXT,
            image_count INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(volume_id) REFERENCES volumes(id) ON DELETE CASCADE
        );`,

        `CREATE TABLE IF NOT EXISTS genres (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE
        );`,

        `CREATE TABLE IF NOT EXISTS series_genres (
            series_id INTEGER NOT NULL,
            genre_id INTEGER NOT NULL,
            PRIMARY KEY (series_id, genre_id),
            FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE,
            FOREIGN KEY(genre_id) REFERENCES genres(id) ON DELETE CASCADE
        );`,

        `CREATE TABLE IF NOT EXISTS tags (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE
        );`,

        `CREATE TABLE IF NOT EXISTS series_tags (
            series_id INTEGER NOT NULL,
            tag_id INTEGER NOT NULL,
            PRIMARY KEY (series_id, tag_id),
            FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE,
            FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
        );`,

        `CREATE TABLE IF NOT EXISTS series_metadata (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            series_id INTEGER NOT NULL,
            title_romaji TEXT,
            title_english TEXT,
            title_native TEXT,
            description TEXT,
            release_date DATE,
            publisher TEXT,
            publication TEXT,
            total_vol INTEGER,
            total_ch INTEGER,
            release_status TEXT CHECK(release_status IN ('Releasing', 'Finished', 'Hiatus', 'Cancelled')),
            FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE
        );`,
    }

    for _, schema := range schemas {
        _, err := db.Exec(schema)
        if err != nil {
            log.Fatal(err)
        }
    }
}

type volumeInfo struct {
    number      int
    title       string
    path        string
    numImages   int
    coverImage  string
}

func generateWebPThumb(srcPath, cacheBase, seriesName, volumeName string) (string, error) {
    cacheDir := filepath.Join(cacheBase, seriesName)
    os.MkdirAll(cacheDir, 0755)

    outPath := filepath.Join(cacheDir, volumeName+"_cover.webp") 

    if _, err := os.Stat(outPath); err == nil {
        relPath, _ := filepath.Rel(cacheBase, outPath)
        return filepath.ToSlash(relPath), nil
    }

    f, err := os.Open(srcPath)
    if err != nil {
        return "", err
    }
    defer f.Close()

    var img image.Image

    switch ext := filepath.Ext(srcPath); ext {
    case ".jpg", ".jpeg", ".JPG", ".JPEG":
        img, err = jpeg.Decode(f)
    case ".png", ".PNG":
        img, err = png.Decode(f)
    case ".webp", ".WEBP":
        img, err = webp.Decode(f)
    default:
        return "", fmt.Errorf("unsupported image type: %s ", ext)
    }
    if err != nil {
        return "", err
    }

    thumb := resize.Resize(320, 0, img, resize.Lanczos2)

    outFile, err := os.Create(outPath)
    if err != nil {
        return "", err
    }
    defer outFile.Close()

    err = webp.Encode(outFile, thumb, &webp.Options{Quality: 100})
    if err != nil {
        return "", err
    }    

    relPath, _ := filepath.Rel(cacheBase, outPath)
    return filepath.ToSlash(relPath), nil
}

func scanSubfolders(path string, entry os.DirEntry, db *sql.DB, insertSerie string) {
    seriesPath := filepath.Join(path, entry.Name())
    cacheFolder := "E:/$otaku/mangaserver_cache/thumbnails"
   
    subfolders, err := os.ReadDir(seriesPath)
    if err != nil {
        fmt.Println("Error reading subfolders:", err)
        return
    }

    folderCount := 0
    totalImages := 0
    mainCoverSet := false
    var seriesMainCover string
    var volumes []volumeInfo

    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    sort.Slice(subfolders, func(i, j int) bool {
        return natural.Less(subfolders[i].Name(), subfolders[j].Name())
    })

    for _, sub := range subfolders {
        if sub.IsDir() {
            folderCount++
            volumePath := filepath.Join(seriesPath, sub.Name())

            files, err := os.ReadDir(volumePath)
            if err != nil {
                fmt.Println("Error reading files:", err)
                continue
            }

            imageCount := 0
            var volumeCoverPath string

            for _, f := range files {
                switch filepath.Ext(f.Name()) {
                case ".jpg", ".JPG", ".jpeg", ".png", ".PNG", ".gif", ".webp", ".bmp":
                    imageCount++
                    if volumeCoverPath == "" {
                        imgPath := filepath.Join(volumePath, f.Name())
                        relPath, err := generateWebPThumb(imgPath, cacheFolder, entry.Name(), sub.Name())
                        if err != nil {
                            fmt.Print("Thumbnail error: ", err)
                        }
                        volumeCoverPath = filepath.ToSlash(relPath)
                    }
                    if !mainCoverSet {
                        imgPath := filepath.Join(volumePath, f.Name())
                        relPath, err := generateWebPThumb(imgPath, cacheFolder, entry.Name(), sub.Name())
                        if err != nil {
                            fmt.Print("Thumbnail error: ", err)
                        }
                        seriesMainCover = filepath.ToSlash(relPath)
                        mainCoverSet = true
                    }
                }
            }
            totalImages += imageCount

            volumes = append(volumes, volumeInfo{
                number:     folderCount,
                title:      sub.Name(),
                path:       volumePath,
                numImages:  imageCount,
                coverImage: volumeCoverPath,
            })
        }
    }

    res, err := tx.Exec(insertSerie, entry.Name(), seriesPath, seriesMainCover, folderCount, totalImages)
    if err != nil {
        tx.Rollback()
        log.Fatal("Insert series error:", err)
    }

    seriesID, err := res.LastInsertId()
    if err != nil {
        tx.Rollback()
        log.Fatal("Fetch seriesID error:", err)
    }

    insertVolume := `INSERT INTO volumes (series_id, number, num_images, title, path, cover)
                     VALUES (?, ?, ?, ?, ?, ?)`
    
    for _, v := range volumes {
        _, err = tx.Exec(insertVolume, seriesID, v.number, v.numImages, v.title, v.path, v.coverImage)
        if err != nil {
            tx.Rollback()
            log.Fatal("Insert volume error:", err)
        }
    }

    if err := tx.Commit(); err != nil {
        log.Fatal("Commit error on transaction: ", err)
    }


    //fmt.Printf("  Total subfolders: %d, Total images: %d\n\n", folderCount, totalImages)
    //fmt.Println("Scanning series:", entry.Name())
}

func FullFolderScan() {
    // Define the root path where your manga series are stored, 
    // TODO let users pick this 
	rootPath := "E:/$otaku/Mangaserver"
    fmt.Println("Starting scan")

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
    
    // Open a database connection
    db, err := sql.Open("sqlite3", "./manga.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    initDB(db)
    
    insertSerie := `INSERT OR IGNORE INTO series (title, path, cover_image, num_vol, num_images) VALUES (?, ?, ?, ?, ?);`

    // Main loop to scan the root directory
    fmt.Println("Scanning for series...")
    const maxGoroutines = 10
    const maxjobs = 10

    jobs := make(chan os.DirEntry, maxjobs)
    var wg sync.WaitGroup

    for range maxGoroutines {
        go func() {
            for entry := range jobs {
                scanSubfolders(rootPath, entry, db, insertSerie)
                wg.Done()
            }
        }()
    }   
    
    for _, entry := range entries {
        if entry.IsDir() {
            wg.Add(1)
            jobs <- entry
        }
    }

    close(jobs)
    wg.Wait()
    fmt.Println("Scan complete.")
}

