package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func render(w http.ResponseWriter, name string, status int, data PageData) {
	var body bytes.Buffer
	if err := pages.ExecuteTemplate(&body, name, data); err != nil {
		log.Println("Ошибка шаблона:", err)
		http.Error(w, "Не удалось отобразить страницу.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(body.Bytes())
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		http.Error(w, "Метод не поддерживается. Используйте "+method+".", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	render(w, "index.html", http.StatusOK, taskPageData())
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	render(w, "about.html", http.StatusOK, PageData{})
}

func showError(w http.ResponseWriter, status int, message, title string) {
	data := taskPageData()
	data.Error, data.Title = message, title
	render(w, "index.html", status, data)
}

func parseForm(w http.ResponseWriter, r *http.Request) bool {
	if !requireMethod(w, r, http.MethodPost) {
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		showError(w, http.StatusBadRequest, "Не удалось прочитать форму. Проверьте отправленные данные.", "")
		return false
	}
	return true
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !parseForm(w, r) {
		return
	}
	title := r.PostForm.Get("title")
	if err := addTask(title); err != nil {
		showError(w, http.StatusBadRequest, err.Error(), title)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func taskID(w http.ResponseWriter, r *http.Request) (int, bool) {
	if !parseForm(w, r) {
		return 0, false
	}
	id, err := strconv.Atoi(r.PostForm.Get("id"))
	if err != nil || id <= 0 {
		showError(w, http.StatusBadRequest, "ID задачи должен быть положительным целым числом.", "")
		return 0, false
	}
	return id, true
}

func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := taskID(w, r)
	if !ok {
		return
	}
	if !toggleTask(id) {
		showError(w, http.StatusNotFound, fmt.Sprintf("Задача с ID %d не найдена.", id), "")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := taskID(w, r)
	if !ok {
		return
	}
	if !deleteTask(id) {
		showError(w, http.StatusNotFound, fmt.Sprintf("Задача с ID %d не найдена.", id), "")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
