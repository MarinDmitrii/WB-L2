package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	calendarBuilder "github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/builder"
	calendarPorts "github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/ports"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503.
		В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400.
		В случае остальных ошибок сервер должен возвращать HTTP 500.
		Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

type HttpServer struct {
	calendarPorts.HttpCalendarHandler
}

type Application struct {
	httpServer *http.Server
}

func (a *Application) Run(addr string, debug bool) error {
	router := http.NewServeMux()

	ctx := context.Background()

	calendarApp := calendarBuilder.NewApplication(ctx)
	calendarHttpHandler := calendarPorts.NewHttpCalendarHandler(calendarApp)
	calendarPorts.CustomRegisterHandlers(router, calendarHttpHandler)

	a.httpServer = &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    180 * time.Second,
		WriteTimeout:   180 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Server is running...")

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	doneCh := make(chan struct{})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	close(doneCh)

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func main() {
	app := &Application{}
	err := app.Run(":8080", false)
	if err != nil {
		panic(err)
	}
}
