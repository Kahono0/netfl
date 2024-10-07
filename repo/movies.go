package repo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"encoding/json"

	"github.com/gabriel-vasile/mimetype"
	"github.com/kahono0/netfl/utils"
)

var acceptedMimeTypes = []string{
	"video/mp4",           // MP4
	"video/x-msvideo",     // AVI
	"video/x-matroska",    // MKV
	"video/quicktime",     // MOV
	"video/x-ms-wmv",      // WMV
	"video/x-flv",         // FLV
	"video/webm",          // WEBM
	"video/mpeg",          // MPEG
	"video/3gpp",          // 3GP
	"video/ogg",           // OGV
	"video/x-matroska-3d", // MK3D
	"video/x-m4v",         // M4V
	"video/x-ms-asf",      // ASF
	"video/mp2t",          // TS
	"video/x-ms-dvr",      // DVR-MS
	"application/mxf",     // MXF
	"video/avchd",         // AVCHD
}

var playableByWebBrowsers = []string{
	"video/mp4",       // MP4
	"video/webm",      // WebM
	"video/ogg",       // OGG
	"video/x-flv",     // FLV (limited support)
	"video/quicktime", // MOV (limited support)
	"video/mpeg",      // MPEG
	"video/3gpp",      // 3GP
}

var thumbNailTime = "00:00:01"

// returns mime type string if it is an accepted mime type
// else returns an empty string
func getMimeType(mime *mimetype.MIME) string {
	for _, t := range acceptedMimeTypes {
		if mime.Is(t) {
			return t
		}
	}

	return ""
}

func isPlayableByWebBrowser(mime *mimetype.MIME) bool {
	for _, t := range playableByWebBrowsers {
		if mime.Is(t) {
			return true
		}
	}

	return false
}

type Movie struct {
	Name                   string
	Path                   string
	MovieUrl               string
	MimeType               string
	ThumbNailUrl           string
	IsAcceptedMimeType     bool
	IsPlayableByWebBrowser bool
}

func (m *Movie) CreateThumbnail(hostAdrr string) error {
	outputFile := fmt.Sprintf("%s.thumb.jpg", m.Path)
	err := utils.ExtractThumbnail(m.Path, outputFile, thumbNailTime)
	if err != nil {
		return err
	}

	m.ThumbNailUrl = fmt.Sprintf("%s/%s", hostAdrr, filepath.Base(outputFile))
	return nil
}

type MovieRepo struct {
	Dir      string
	HostAddr string
	Log      bool
	Movies   []Movie
	Loaded   bool
}

func NewMovieRepo(dir, hostAddr string, log bool) *MovieRepo {
	return &MovieRepo{
		Dir:      dir,
		HostAddr: hostAddr,
		Log:      log,
		Loaded:   false,
	}
}

func (r *MovieRepo) Load() error {
	var movies []Movie

	err := filepath.Walk(r.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		timeToCreate := time.Now()
		mime, err := mimetype.DetectFile(path)
		if err != nil {
			return err
		}

		mimeTypeStr := getMimeType(mime)
		if mimeTypeStr == "" {
			if r.Log {
				log.Println("Skipping file", path, "because it has an unsupported mime type:", mime.String())
			}
			return nil
		}

		movie := Movie{
			Name:                   filepath.Base(path),
			Path:                   path,
			MovieUrl:               fmt.Sprintf("%s/%s", r.HostAddr, filepath.Base(path)),
			MimeType:               mimeTypeStr,
			IsAcceptedMimeType:     true,
			IsPlayableByWebBrowser: isPlayableByWebBrowser(mime),
		}

		err = movie.CreateThumbnail(r.HostAddr)
		if err != nil {
			return err
		}

		movies = append(movies, movie)

		fmt.Println("Time to create:", time.Since(timeToCreate))

		return nil
	})

	if err != nil {
		return err
	}

	r.Movies = movies
	r.Loaded = true

	return nil
}

func (r *MovieRepo) String() string {
	if !r.Loaded {
		r.Load()
	}
	var str string

	for _, m := range r.Movies {
		str += fmt.Sprintf("Name: %s\n", m.Name)
	}

	return str
}

func (r *MovieRepo) GetMovies(w http.ResponseWriter, rq *http.Request) {
	if !r.Loaded {
		r.Load()
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(r.Movies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *MovieRepo) DetectMimeType(file string) {
	fmt.Println(http.DetectContentType([]byte(file)))
}
