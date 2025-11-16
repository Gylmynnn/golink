package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/Gylmynnn/golink/models"
	"github.com/Gylmynnn/golink/storage"
	"github.com/go-chi/chi/v5"
)

var store = storage.NewRedisStore()

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := store.SaveURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	res := models.ShortenResponse{
		ShortURL: shortURL,
	}
	json.NewEncoder(w).Encode(res)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")

	originalURL, err := store.GetOriginalURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func GetTopDomains(w http.ResponseWriter, r *http.Request) {
	domainCounts, err := store.GetDomainCounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type kv struct {
		Key   string
		Value int
	}

	var shortedDomains []kv
	for k, v := range domainCounts {
		shortedDomains = append(shortedDomains, kv{k, v})
	}
	sort.Slice(shortedDomains, func(i, j int) bool { return shortedDomains[i].Value > shortedDomains[j].Value })

	topDomains := make(map[string]int)
	for i, domain := range shortedDomains {
		if i >= 3 {
			break
		}
		topDomains[domain.Key] = domain.Value
	}
	json.NewEncoder(w).Encode(topDomains)
}

func ShortenURLHTML(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='error'>URL is required</div>"))
		return
	}
	shortURL, err := store.SaveURL(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='error'>" + err.Error() + "</div>"))
		return
	}
	w.Write([]byte("<div class='result'>Short URL: <a href='" + shortURL + "' target='_blank'>" + shortURL + "</a></div>"))
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func GetTopDomainsHTML(w http.ResponseWriter, r *http.Request) {
	domainCounts, err := store.GetDomainCounts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div class='error'>" + err.Error() + "</div>"))
		return
	}
	type kv struct {
		Key   string
		Value int
	}
	var shortedDomains []kv
	for k, v := range domainCounts {
		shortedDomains = append(shortedDomains, kv{k, v})
	}
	sort.Slice(shortedDomains, func(i, j int) bool { return shortedDomains[i].Value > shortedDomains[j].Value })
	w.Write([]byte("<ul>"))
	for i, domain := range shortedDomains {
		if i >= 3 {
			break
		}
		w.Write([]byte("<li>" + domain.Key + " (" + fmt.Sprintf("%d", domain.Value) + ")</li>"))
	}
	w.Write([]byte("</ul>"))
}
