package main

import (
	"encoding/json"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "jetstyle/docs"
	"jetstyle/internal/taskstore"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) taskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		// Запрос направлен к "/task/", без идущего в конце ID.
		if req.Method == http.MethodPost {
			ts.createTaskHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ts.deleteTaskHandler(w, req, id)
		} else if req.Method == http.MethodGet {
			ts.getTaskHandler(w, req, id)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /task/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ts *taskServer) swaggerHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/swagger/index.html" {
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		)
	}
}

// Task godoc
//
//	@Summary		Task
//	@Description	Get one task
//	@Tags			Tasks
//	@Param			id	path		int		true	"Id"
//	@Success		200	{string}	string	"ok"
//	@Failure		400	{string}	string	"error"
//
//	@Router			/task/{id} [get]
func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

// Task godoc
//
//	@Summary		Task
//	@Description	delete one task
//	@Tags			Tasks
//	@Param			id	path		int		true	"Id"
//	@Success		200	{string}	string	"ok"
//	@Failure		400	{string}	string	"error"
//
//	@Router			/task/{id} [delete]
func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling delete task at %s\n", req.URL.Path)

	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

type RequestTask struct {
	Text string    `json:"text"`
	Due  time.Time `json:"due"`
}

// Task godoc
//
//	@Summary		Task
//	@Description	Create one task
//	@Tags			Tasks
//	@Param			task	body		RequestTask	true	"Add account"
//	@Success		200		{string}	string		"ok"
//	@Failure		400		{string}	string		"error"
//
//	@Router			/task/ [post]
func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.

	type ResponseId struct {
		Id int `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Due)
	renderJSON(w, ResponseId{Id: id})
}

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
func main() {
	log.Printf("Start docker")
	mux := http.NewServeMux()
	server := NewTaskServer()
	mux.HandleFunc("/task/", server.taskHandler)
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(":8080/swagger/doc.json"),
	))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
