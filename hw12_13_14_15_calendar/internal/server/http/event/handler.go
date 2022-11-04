package event

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/query"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/pkg/api/openapi"
)

type eventHandler struct {
	app    app.Application
	logger app.Logger
}

func NewEventHandler(app app.Application, logger app.Logger) openapi.ServerInterface {
	return &eventHandler{app: app, logger: logger}
}

func (h *eventHandler) GetEvents(w http.ResponseWriter, r *http.Request, params openapi.GetEventsParams) {
	var events []event.Event
	switch params.Period {
	case openapi.Day:
		response, err := h.app.GetDayEvents(r.Context(), query.GetDayEventsRequest{Date: time.Unix(params.Date, 0)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		events = response.Events
	case openapi.Week:
		response, err := h.app.GetWeekEvents(r.Context(), query.GetWeekEventsRequest{Date: time.Unix(params.Date, 0)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		events = response.Events
	case openapi.Month:
		response, err := h.app.GetMonthEvents(r.Context(), query.GetMonthEventsRequest{Date: time.Unix(params.Date, 0)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		events = response.Events
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.writeEvents(w, convertEventsDomainToAPI(events))
}

func (h *eventHandler) PostEvents(w http.ResponseWriter, r *http.Request) {
	var apiNewEvent openapi.NewEvent
	if err := json.NewDecoder(r.Body).Decode(&apiNewEvent); err != nil {
		h.logger.Error("openapi.NewEvent json decode error: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := command.CreateEventRequest{
		Title:    apiNewEvent.Title,
		DateTime: time.Unix(apiNewEvent.Date, 0),
		Duration: time.Duration(apiNewEvent.Duration) * time.Minute,
		UserID:   "10",
	}
	if apiNewEvent.Description != nil {
		request.Description = *apiNewEvent.Description
	}
	if apiNewEvent.NotifyBefore != nil {
		request.NotifyBefore = time.Duration(*apiNewEvent.NotifyBefore) * time.Minute
	}
	response, err := h.app.CreateEvent(r.Context(), request)
	if err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.writeEvent(w, response.Event)
}

func (h *eventHandler) PutEvents(w http.ResponseWriter, r *http.Request) {
	var apiEvent openapi.Event
	if err := json.NewDecoder(r.Body).Decode(&apiEvent); err != nil {
		h.logger.Error("openapi.Event json decode error: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := command.UpdateEventRequest{
		ID:       apiEvent.Id,
		Title:    apiEvent.Title,
		DateTime: time.Unix(apiEvent.Date, 0),
		Duration: time.Duration(apiEvent.Duration) * time.Minute,
		UserID:   "10",
	}
	if apiEvent.Description != nil {
		request.Description = *apiEvent.Description
	}
	if apiEvent.NotifyBefore != nil {
		request.NotifyBefore = time.Duration(*apiEvent.NotifyBefore) * time.Minute
	}
	if err := h.app.UpdateEvent(r.Context(), request); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *eventHandler) DeleteEventsID(w http.ResponseWriter, r *http.Request, id string) {
	err := h.app.DeleteEvent(r.Context(), command.DeleteEventRequest{ID: id})
	if err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func convertEventsDomainToAPI(events []event.Event) []openapi.Event {
	result := make([]openapi.Event, len(events))
	for i, e := range events {
		result[i] = openapi.Event{
			Date:         e.DateTime.Unix(),
			Description:  nil,
			Duration:     int(e.Duration.Minutes()),
			Id:           e.ID,
			NotifyBefore: nil,
			Title:        e.Title,
		}
		if e.Description != "" {
			d := e.Description
			result[i].Description = &d
		}
		if e.NotifyBefore.Minutes() > 0 {
			n := int(e.NotifyBefore.Minutes())
			result[i].NotifyBefore = &n
		}
	}
	return result
}

func (h *eventHandler) writeEvents(w http.ResponseWriter, events []openapi.Event) {
	resBuf, err := json.Marshal(events)
	if err != nil {
		h.logger.Error("Events json marshal error: " + err.Error())
	}
	_, err = w.Write(resBuf)
	if err != nil {
		h.logger.Error("Response write error: " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (h *eventHandler) writeEvent(w http.ResponseWriter, event event.Event) {
	responseEvent := openapi.Event{
		Date:     event.DateTime.Unix(),
		Duration: int(event.Duration / time.Minute),
		Id:       event.ID,
		Title:    event.Title,
	}
	if event.Description != "" {
		responseEvent.Description = &event.Description
	}
	if event.NotifyBefore > 0 {
		nb := int(event.NotifyBefore / time.Minute)
		responseEvent.NotifyBefore = &nb
	}
	resBuf, err := json.Marshal(responseEvent)
	if err != nil {
		h.logger.Error("Event json marshal error: " + err.Error())
	}
	_, err = w.Write(resBuf)
	if err != nil {
		h.logger.Error("Response write error: " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (h *eventHandler) writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, event.ErrEventNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, event.ErrDateBusy) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	var errValidate command.ErrValidate
	if errors.As(err, &errValidate) {
		w.WriteHeader(http.StatusBadRequest)
		errResponse := openapi.Error{Message: err.Error()}
		resBuf, errMarshal := json.Marshal(errResponse)
		if errMarshal != nil {
			h.logger.Error("Event json marshal error: " + errMarshal.Error())
		}
		_, errWrite := w.Write(resBuf)
		if errWrite != nil {
			h.logger.Error("Response write error: " + errWrite.Error())
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
