package ports

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/builder"
	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type HttpCalendarHandler struct {
	app *builder.Application
}

func NewHttpCalendarHandler(app *builder.Application) HttpCalendarHandler {
	return HttpCalendarHandler{app: app}
}

func (h HttpCalendarHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// Проверка на соответствие метода запроса
	if r.Method != http.MethodPost {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Валиадции и парсинг параметров
	event, statusCode, errMessage := h.validationAndParse(r)
	if statusCode != 200 {
		// Исправить описание, потому что на текущий момент мы возвращаем при валидации не только 400, но и 500 и 503! -----------------
		// Если ошибка во входных данных, возвращаем HTTP 400
		h.mapToResponse(w, statusCode, nil, errMessage)
		return
	}

	// Проверка обязательных полей
	if event.UserID == 0 || event.Date == (time.Time{}) || event.Description == "" {
		// Если обязательные входные данные отсутсвуют, возвращаем HTTP 400
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	_, err := h.app.CreateEvent.Execute(r.Context(), event)
	if err != nil {
		// Если ошибка в бизнес-логике, возвращаем HTTP 503
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	h.mapToResponse(w, http.StatusOK, event, "")
}

func (h HttpCalendarHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// Проверка на соответствие метода запроса
	if r.Method != http.MethodPost {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Валиадции и парсинг параметров
	event, statusCode, errMessage := h.validationAndParse(r)
	if statusCode != 200 {
		// Исправить описание, потому что на текущий момент мы возвращаем при валидации не только 400, но и 500 и 503! -----------------
		// Если ошибка во входных данных, возвращаем HTTP 400
		h.mapToResponse(w, statusCode, nil, errMessage)
		return
	}

	// Проверка обязательных полей
	if event.ID == 0 || event.UserID == 0 || event.Date == (time.Time{}) || event.Description == "" {
		// Если обязательные входные данные отсутсвуют, возвращаем HTTP 400
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	err := h.app.UpdateEvent.Execute(r.Context(), event)
	if err != nil {
		// Если ошибка в бизнес-логике, возвращаем HTTP 503
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	h.mapToResponse(w, http.StatusOK, event, "")
}

func (h HttpCalendarHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	// Проверка на соответствие метода запроса
	if r.Method != http.MethodPost {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Валиадции и парсинг параметров
	event, statusCode, errMessage := h.validationAndParse(r)
	if statusCode != 200 {
		// Исправить описание, потому что на текущий момент мы возвращаем при валидации не только 400, но и 500 и 503! -----------------
		// Если ошибка во входных данных, возвращаем HTTP 400
		h.mapToResponse(w, statusCode, nil, errMessage)
		return
	}

	// Проверка обязательных полей
	if event.ID == 0 {
		// Если обязательные входные данные отсутсвуют, возвращаем HTTP 400
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	// Нужна ли эта промежуточная проверка по нахождению эвента в БД? ------------------------------------------
	// event, err := h.app.GetEventByID.Execute(r.Context(), event.ID)
	// if err != nil {
	// 	// Если ошибка в бизнес-логике, возвращаем HTTP 503
	// 	h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
	// 	return
	// }

	err := h.app.DeleteEvent.Execute(r.Context(), event.ID)
	if err != nil {
		// Если ошибка в бизнес-логике, возвращаем HTTP 503
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	h.mapToResponse(w, http.StatusOK, event, "")
}

func (h HttpCalendarHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	// Проверка на соответствие метода запроса
	if r.Method != http.MethodGet {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Валиадции и парсинг параметров
	event, statusCode, errMessage := h.validationAndParse(r)
	if statusCode != 200 {
		// Исправить описание, потому что на текущий момент мы возвращаем при валидации не только 400, но и 500 и 503! -----------------
		// Если ошибка во входных данных, возвращаем HTTP 400
		h.mapToResponse(w, statusCode, nil, errMessage)
		return
	}

	// Проверка обязательных полей
	if event.Date == (time.Time{}) {
		// Если обязательные входные данные отсутсвуют, возвращаем HTTP 400
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	events, err := h.app.GetEventsForDay.Execute(r.Context(), event.Date)
	if err != nil {
		// Если ошибка в бизнес-логике, возвращаем HTTP 503
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	// Возможно потребуется изменить формат выдачи events -------------------------------------------------------------------
	h.mapToResponse(w, http.StatusOK, events, "")
}

func (h HttpCalendarHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	// Проверка на соответствие метода запроса
	if r.Method != http.MethodGet {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Валиадции и парсинг параметров
	event, statusCode, errMessage := h.validationAndParse(r)
	if statusCode != 200 {
		// Исправить описание, потому что на текущий момент мы возвращаем при валидации не только 400, но и 500 и 503! -----------------
		// Если ошибка во входных данных, возвращаем HTTP 400
		h.mapToResponse(w, statusCode, nil, errMessage)
		return
	}

	// Проверка обязательных полей
	if event.Date == (time.Time{}) {
		// Если обязательные входные данные отсутсвуют, возвращаем HTTP 400
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	events, err := h.app.GetEventsForWeek.Execute(r.Context(), event.Date)
	if err != nil {
		// Если ошибка в бизнес-логике, возвращаем HTTP 503
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	// Возможно потребуется изменить формат выдачи events -------------------------------------------------------------------
	h.mapToResponse(w, http.StatusOK, events, "")
}

func (h HttpCalendarHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	// Проверка на соответствие метода запроса
	if r.Method != http.MethodGet {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Валиадции и парсинг параметров
	event, statusCode, errMessage := h.validationAndParse(r)
	if statusCode != 200 {
		// Исправить описание, потому что на текущий момент мы возвращаем при валидации не только 400, но и 500 и 503! -----------------
		// Если ошибка во входных данных, возвращаем HTTP 400
		h.mapToResponse(w, statusCode, nil, errMessage)
		return
	}

	// Проверка обязательных полей
	if event.Date == (time.Time{}) {
		// Если обязательные входные данные отсутсвуют, возвращаем HTTP 400
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	events, err := h.app.GetEventsForMonth.Execute(r.Context(), event.Date)
	if err != nil {
		// Если ошибка в бизнес-логике, возвращаем HTTP 503
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	// Возможно потребуется изменить формат выдачи events -------------------------------------------------------------------
	h.mapToResponse(w, http.StatusOK, events, "")
}

// Нужна ли эта функция? -----------------------------------------------------------
// func (h HttpCalendarHandler) GetEventByUID(w http.ResponseWriter, r *http.Request, eventID int) (domain.Event, error) {
// 	event, err := h.app.GetEventByUID.Execute(r.Context(), eventID)
// 	if err != nil {
// 		log.Println("здесь косяк!")
// 		return domain.Event{}, fmt.Errorf("%v, %v", http.StatusInternalServerError, err.Error())
// 	}

// 	return event, nil
// }

type jsonEvent struct {
	ID          int       `json:"event_id"`
	UserID      int       `json:"user_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

// Валидация и парсинг параметров
func (h HttpCalendarHandler) validationAndParse(r *http.Request) (domain.Event, int, string) {
	event := domain.Event{}

	if r.Method == http.MethodGet || (r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			// Или нужно возвращать HTTP 400 ??? Или HTTP 503 ??? ---------------------------------------------------------
			// Если ошибка при парсинге данных, возвращаем HTTP 500
			return domain.Event{}, http.StatusInternalServerError, err.Error()
		}

		eventID := 0
		if r.Form.Get("event_id") != "" {
			eventID, err = strconv.Atoi(r.Form.Get("event_id"))
			if err != nil {
				// Если ошибка валидации входных данных, возвращаем HTTP 400
				return domain.Event{}, http.StatusBadRequest, err.Error()
			}
		}

		userID, err := strconv.Atoi(r.Form.Get("user_id"))
		if err != nil {
			// Если ошибка валидации входных данных, возвращаем HTTP 400
			return domain.Event{}, http.StatusBadRequest, err.Error()
		}

		date, err := time.Parse("2006-01-02", r.Form.Get("date"))
		if err != nil {
			// Если ошибка валидации входных данных, возвращаем HTTP 400
			return domain.Event{}, http.StatusBadRequest, err.Error()
		}

		description := r.Form.Get("description")

		event.ID = eventID
		event.UserID = userID
		event.Date = date
		event.Description = description
	} else if r.Header.Get("Content-Type") == "application/json" {
		jEvent := jsonEvent{}

		err := json.NewDecoder(r.Body).Decode(&jEvent)
		if err != nil {
			// Или нужно возвращать HTTP 400 ??? Или HTTP 503 ??? ---------------------------------------------------------
			// Если ошибка при декодировании данных, возвращаем HTTP 500
			return domain.Event{}, http.StatusInternalServerError, err.Error()
		}

		event.ID = jEvent.ID
		event.UserID = jEvent.UserID
		event.Date = jEvent.Date
		event.Description = jEvent.Description
	} else {
		// Если необходимые входные данные отсутсвуют, возвращаем HTTP 400
		return domain.Event{}, http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
	}

	return event, http.StatusOK, ""
}

// Парсинг ответа
func (h HttpCalendarHandler) mapToResponse(w http.ResponseWriter, statusCode int, data interface{}, errMessage string) {
	// Задаём JSON формат в Content-Type заголовка ответа
	w.Header().Set("Content-Type", "application/json")

	// Записываем HTTP код статуса ответа в заголовок ответа
	w.WriteHeader(statusCode)

	// Создаём JSON объект для ответа
	response := make(map[string]interface{})

	// В зависимости от статуса запроса, формируем JSON
	if statusCode >= 200 && statusCode < 300 {
		response["result"] = data
	} else {
		response["error"] = errMessage
	}

	// Преобразуем данные в JSON и записываем в тело ответа
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Или нужно возвращать HTTP 400 ??? Или HTTP 503 ??? ---------------------------------------------------------
		// Если ошибка при кодировании данных, возвращаем HTTP 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h HttpCalendarHandler) MiddlewareLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем информацию о запросе
		log.Printf("Request: %s %s\nBody: %s", r.Method, r.URL.Path, r.Form)
		// Передаем запрос следующему обработчику
		next(w, r)
	}
}

func CustomRegisterHandlers(router *http.ServeMux, h HttpCalendarHandler) {
	router.HandleFunc("/create_event", h.MiddlewareLogger(h.CreateEvent))
	router.HandleFunc("/update_event", h.MiddlewareLogger(h.UpdateEvent))
	router.HandleFunc("/delete_event", h.MiddlewareLogger(h.DeleteEvent))
	router.HandleFunc("/events_for_day", h.MiddlewareLogger(h.GetEventsForDay))
	router.HandleFunc("/events_for_week", h.MiddlewareLogger(h.GetEventsForWeek))
	router.HandleFunc("/events_for_month", h.MiddlewareLogger(h.GetEventsForMonth))
}
