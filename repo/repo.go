package repo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	ID                     int
	Name                   string
	MovieUrl               string
	MimeType               string
	ThumbNailUrl           string
	IsAcceptedMimeType     bool
	IsPlayableByWebBrowser bool
	Hash                   string
}

func (m *Movie) CreateThumbnail(hostAdrr string, path string, dir string) error {
	outputFile := fmt.Sprintf("%s.thumb.jpg", path)
	err := utils.ExtractThumbnail(path, outputFile, thumbNailTime)
	if err != nil {
		return err
	}

	m.ThumbNailUrl = fmt.Sprintf("%s/thumb%s", hostAdrr, outputFile[len(dir):])
	return nil
}

type MovieRepo struct {
	Dir      string
	HostAddr string
	Log      bool
	Movies   []Movie
	Loaded   bool
	NextID   int
}

func NewMovieRepo(dir, hostAddr string, log bool) *MovieRepo {
	repo := &MovieRepo{
		Dir:      dir,
		HostAddr: hostAddr,
		Log:      log,
		Loaded:   false,
		NextID:   0,
	}

	go repo.Load()

	return repo
}

func (r *MovieRepo) Load() error {
	if r.Dir == "" {
		return nil
	}

	worker := NewThumbNailGenWorker(r)
	worker.Start()

	defer worker.Close()

	err := filepath.Walk(r.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

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
			ID:                     r.NextID,
			Name:                   filepath.Base(path),
			MovieUrl:               fmt.Sprintf("%s/movies%s", r.HostAddr, path[len(r.Dir):]),
			MimeType:               mimeTypeStr,
			IsAcceptedMimeType:     true,
			IsPlayableByWebBrowser: isPlayableByWebBrowser(mime),
		}

		hash := utils.GetFileHash(path)
		movie.Hash = hash

		worker.AddJob(Job{
			Movie: movie,
			Path:  path,
		})

		r.AddMovie(movie)

		return nil
	})

	if err != nil {
		return err
	}

	r.Loaded = true

	return nil
}

func (r *MovieRepo) AddMovie(m Movie) {
	m.ID = r.NextID
	r.Movies = append(r.Movies, m)
	r.NextID++
}

func (r *MovieRepo) UpdateMovie(m Movie) {
	for i, movie := range r.Movies {
		if movie.Hash == m.Hash {
			r.Movies[i] = m
			break
		}
	}
}

func (r *MovieRepo) GetMovies() []Movie {
	if !r.Loaded {
		r.Load()
	}

	return r.Movies
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

func (r *MovieRepo) ToJSON() string {
	jsonB, _ := json.Marshal(r.Movies)
	return string(jsonB)
}

func (r *MovieRepo) ContainsFile(hash string) bool {
	for _, m := range r.Movies {
		if m.Hash == hash {
			return true
		}

	}

	return false
}

func (r *MovieRepo) AddMovies(movies []Movie) {
	for _, m := range movies {
		if m.Hash == "" || r.ContainsFile(m.Hash) {
			continue
		}
		r.AddMovie(m)
	}
}

func (r *MovieRepo) AddFromJSON(jsonStr string) error {
	var movies []Movie
	err := json.Unmarshal([]byte(jsonStr), &movies)
	if err != nil {
		return err
	}

	r.AddMovies(movies)

	return nil
}

func (r *MovieRepo) DetectMimeType(file string) {
	fmt.Println(http.DetectContentType([]byte(file)))
}
