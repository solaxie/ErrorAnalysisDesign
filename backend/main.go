package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

var session *gocql.Session

type Session struct {
	ObjectType string              `json:"object_type"`
	UserID     string              `json:"user_id"`
	Progress   map[string]Progress `json:"progress"`
}

type Progress struct {
	ImageName string `json:"image_name"`
}

type Feedback struct {
	ObjectType string      `json:"object_type"`
	UserID     string      `json:"user_id"`
	Images     []ImageInfo `json:"images"`
}

type ImageInfo struct {
	Information map[string][]interface{} `json:"information"`
}

type Attitude struct {
	Name     string `json:"name"`
	Progress int    `json:"progress"`
	Total    int    `json:"total"`
}

func main() {
	// Initialize ScyllaDB connection
	cluster := gocql.NewCluster(os.Getenv("SCYLLA_HOST"))
	cluster.Keyspace = "swipe_mission"
	cluster.Consistency = gocql.Quorum

	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Set up router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/attitude/{attitude}", attitudeHandler).Methods("GET")
	r.HandleFunc("/api/attitudes", getAttitudesHandler).Methods("GET")
	r.HandleFunc("/api/next-image/{attitude}", getNextImageHandler).Methods("GET")
	r.HandleFunc("/api/feedback", feedbackHandler).Methods("POST")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend")))

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/index.html")
}

func attitudeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	attitude := vars["attitude"]

	// TODO: Implement session handling and image retrieval
	imagePath, err := getNextImage(attitude, "dummy-uuid")
	if err != nil {
		http.Error(w, "Error getting next image", http.StatusInternalServerError)
		return
	}

	// TODO: Implement proper template rendering
	fmt.Fprintf(w, `
		<div class="swipe-instructions">swipe: right(correct), left(wrong), down(undo), up(save & exit)</div>
		<div class="attitude-value">%s</div>
		<div class="image-container">
			<img src="%s" data-image-name="%s" data-attitude="%s">
		</div>
	`, attitude, imagePath, filepath.Base(imagePath), attitude)
}

func getAttitudesHandler(w http.ResponseWriter, r *http.Request) {
	attitudes, err := getAttitudes()
	if err != nil {
		http.Error(w, "Error getting attitudes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attitudes)
}

func getNextImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	attitude := vars["attitude"]

	// TODO: Implement proper user ID handling
	imagePath, err := getNextImage(attitude, "dummy-uuid")
	if err != nil {
		http.Error(w, "Error getting next image", http.StatusInternalServerError)
		return
	}

	attitudeValue, err := getAttitudeValue(imagePath, attitude)
	if err != nil {
		http.Error(w, "Error getting attitude value", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"imagePath":     imagePath,
		"imageName":     filepath.Base(imagePath),
		"attitudeValue": attitudeValue,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	var feedback struct {
		ImageName string `json:"imageName"`
		Attitude  string `json:"attitude"`
		Action    string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement feedback handling logic with ScyllaDB

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func getAttitudes() ([]Attitude, error) {
	files, err := ioutil.ReadDir("../test_result")
	if err != nil {
		return nil, err
	}

	attitudesMap := make(map[string]int)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".txt" {
			content, err := ioutil.ReadFile(filepath.Join("../test_result", file.Name()))
			if err != nil {
				return nil, err
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					attitudesMap[strings.TrimSpace(parts[0])]++
				}
			}
		}
	}

	attitudes := make([]Attitude, 0, len(attitudesMap))
	for name, count := range attitudesMap {
		attitudes = append(attitudes, Attitude{
			Name:     name,
			Progress: 0, // TODO: Implement progress tracking
			Total:    count,
		})
	}
	return attitudes, nil
}

func getNextImage(attitude, userID string) (string, error) {
	// TODO: Implement proper image selection based on session progress
	files, err := ioutil.ReadDir("../test_image")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".jpg" {
			return filepath.Join("../test_image", file.Name()), nil
		}
	}

	return "", fmt.Errorf("no images found")
}

func getAttitudeValue(imagePath, attitude string) (string, error) {
	txtFile := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + ".txt"
	txtFile = filepath.Join("../test_result", filepath.Base(txtFile))

	content, err := ioutil.ReadFile(txtFile)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == attitude {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("attitude not found in file")
}
