package main

// Импортируем основные модули - fmt, log, net/http, time, math/rand
import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	// Добавляем импорт клиентской библиотеки Prometheus
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Метрика для подсчета времени обработки запроса
	apiDurations = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "app_api_durations_seconds",
			Help:       "API latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})

	// Метрика для подсчета количества входящих запросов
	apiProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "app_api_processed_ops_total",
		Help: "The total number of processed requests",
	})
)

// Функция, которая будет обрабатывать все запросы на наш веб-сервер
func handler(w http.ResponseWriter, r *http.Request) {
	// Увеличиваем счетчик количества входящих запросов
	apiProcessed.Inc()

	// Засекаем время начала обработки запроса
	start := time.Now()

	// Засыпаем на случайное количество секунд - от 0 до 2х
	time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	// В ответ отправляем путь, к которому обратился пользователь
	fmt.Fprintf(w, "%s\n", r.URL.Path)

	// После окончания обработки считаем сколько времени прошло с начала обработки
	duration := time.Since(start)
	// Сохраняем время обработки в метрику
	apiDurations.Observe(duration.Seconds())
}

func main() {
	// Регистрируем наши метрики
	prometheus.MustRegister(apiDurations)
	prometheus.MustRegister(apiProcessed)

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())
	// Определяем, что при запросах на / - то есть по сути на любой http путь, необходимо вызывать функцию handler
	http.HandleFunc("/", handler)

	// При запросе пути /metrics будем выдавать метрики в формате Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Запускаем наш сервер на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
