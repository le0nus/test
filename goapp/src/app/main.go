package main

// Импортируем основные модули - fmt, log, net/http, time, math/rand
import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Функция, которая будет обрабатывать все запросы на наш веб сервер
func handler(w http.ResponseWriter, r *http.Request) {
	// Засыпаем случайное количество секунд - от 0 до 2х
	time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	// В ответ отправляем путь, к которому обратился пользователь
	fmt.Fprintf(w, "%s\n", r.URL.Path)
}

// Главная функция - точка старта нашей программы
func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())
	// Определяем, что при запросах на "/" - то есть по сути на любой http путь, необходимо вызывать функцию handler
	http.HandleFunc("/", handler)
	// Запускаем наш сервер на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
