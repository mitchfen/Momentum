package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed templates static
var embedFS embed.FS

var (
	cfg  *AppConfig
	svc  *TodoService
	tmpl *template.Template
)

func main() {
	var err error
	cfg, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	svc = NewTodoService(cfg)

	tmpl = template.Must(template.ParseFS(embedFS, "templates/*.html"))

	staticFS, err := fs.Sub(embedFS, "static")
	if err != nil {
		log.Fatalf("Failed to create static sub-filesystem: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleIndex)
	mux.HandleFunc("POST /toggle", handleToggle)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Printf("Starting Momentum on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	today := cfg.GetToday()
	items := make([]TodoItem, 0, len(cfg.DailyTasks))
	for _, taskText := range cfg.DailyTasks {
		item := ParseTodoItem(taskText, func(subTaskText string) bool {
			return svc.GetSubTaskState(today, taskText, subTaskText)
		})
		items = append(items, item)
	}
	if err := tmpl.ExecuteTemplate(w, "index", items); err != nil {
		log.Printf("Error rendering index: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func handleToggle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	taskText := r.FormValue("taskText")
	subTaskText := r.FormValue("subTaskText")
	completed := r.FormValue("completed") == "true"

	today := cfg.GetToday()
	svc.SetSubTaskState(today, taskText, subTaskText, completed)

	item := ParseTodoItem(taskText, func(st string) bool {
		return svc.GetSubTaskState(today, taskText, st)
	})
	if err := tmpl.ExecuteTemplate(w, "todo-item", item); err != nil {
		log.Printf("Error rendering todo-item: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
