package todo

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/middlewares"
	"github.com/horlerdipo/todo-golang/utils"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	TodoService *Service
}

func NewHandler(todoService *Service) *Handler {
	return &Handler{
		TodoService: todoService,
	}
}

func (handler *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := utils.JsonValidate[dtos.CreateTodoDTO](w, r)
	if err != nil {
		return
	}

	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)
	jsonResponse.UserID = authDetails.UserId

	_, err = handler.TodoService.CreateTodo(r.Context(), &jsonResponse)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func (handler *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "id")
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)
	todoIdInt, err := strconv.ParseUint(todoId, 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "todo not found", nil)
		return
	}

	err = handler.TodoService.DeleteTodo(r.Context(), uint(todoIdInt), authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Todo not found", nil)
		return
	}

	jsonRequest, err := utils.JsonValidate[dtos.UpdateTodoDTO](w, r)
	if err != nil {
		return
	}

	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)

	err = handler.TodoService.UpdateTodo(r.Context(), uint(todoId), &jsonRequest, authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (handler *Handler) PinTodo(w http.ResponseWriter, r *http.Request) {
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)
	todoId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "todo not found", nil)
		return
	}

	err = handler.TodoService.PinTodo(r.Context(), uint(todoId), authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) UnPinTodo(w http.ResponseWriter, r *http.Request) {
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)
	todoId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "todo not found", nil)
		return
	}

	err = handler.TodoService.UnPinTodo(r.Context(), uint(todoId), authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) AddChecklistItem(w http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "id")
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)
	todoIdInt, err := strconv.ParseUint(todoId, 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "todo not found", nil)
		return
	}

	var checklistItem dtos.ChecklistItem
	checklistItem, err = utils.JsonValidate[dtos.ChecklistItem](w, r)
	if err != nil {
		return
	}

	err = handler.TodoService.AddChecklistItem(r.Context(), uint(todoIdInt), checklistItem.Item, authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (handler *Handler) DeleteChecklistItem(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "todo not found", nil)
		return
	}

	checklistItem, err := strconv.ParseUint(chi.URLParam(r, "itemId"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "checklist not found", nil)
		return
	}

	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)

	err = handler.TodoService.DeleteChecklistItem(r.Context(), uint(checklistItem), uint(todoId), authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) UpdateChecklistItem(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "todo not found", nil)
		return
	}

	checklistItem, err := strconv.ParseUint(chi.URLParam(r, "itemId"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "checklist not found", nil)
		return
	}

	var jsonRequest dtos.ChecklistItem
	jsonRequest, err = utils.JsonValidate[dtos.ChecklistItem](w, r)
	if err != nil {
		return
	}
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)

	_, err = handler.TodoService.UpdateChecklistItem(r.Context(), uint(checklistItem), jsonRequest.Item, uint(todoId), authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) UpdateChecklistItemStatus(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "todo not found", nil)
		return
	}

	checklistItem, err := strconv.ParseUint(chi.URLParam(r, "itemId"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, 400, "checklist not found", nil)
		return
	}

	var jsonRequest dtos.ChecklistStatus
	jsonRequest, err = utils.JsonValidate[dtos.ChecklistStatus](w, r)
	if err != nil {
		return
	}
	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)

	_, err = handler.TodoService.UpdateChecklistItemStatus(r.Context(), uint(checklistItem), jsonRequest.Done, uint(todoId), authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) FetchTodos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	perPage, _ := strconv.Atoi(query.Get("per_page"))
	sortBy := query.Get("sort_by")
	order := query.Get("order")

	filters := make(map[string]string)
	for key, values := range r.URL.Query() {
		if strings.HasPrefix(key, "filters[") {
			field := strings.TrimSuffix(strings.TrimPrefix(key, "filters["), "]")
			filters[field] = values[0]
		}
	}

	authDetails := r.Context().Value(middlewares.UserKey).(middlewares.AuthDetails)

	paginationOptions := dtos.PaginationOptions{
		Page:    page,
		PerPage: perPage,
		SortBy:  sortBy,
		Order:   dtos.Order(order),
		Filters: filters,
	}

	todos, err := handler.TodoService.FetchTodos(r.Context(), paginationOptions, authDetails.UserId)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error(), nil)
		return
	}

	utils.RespondWithPaginatedData(w, http.StatusOK, "Todos fetched successfully", todos.Data, todos.Meta)
	return
}

func (handler *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/todos", func(r chi.Router) {
		r.Use(middlewares.JwtAuthMiddleware(handler.TodoService.TokenBlacklistRepository))
		r.Post("/", handler.CreateTodo)
		r.Delete("/{id}", handler.DeleteTodo)
		r.Patch("/{id}", handler.UpdateTodo)
		r.Patch("/{id}/pin", handler.PinTodo)
		r.Patch("/{id}/unpin", handler.UnPinTodo)
		r.Get("/", handler.FetchTodos)

		//Checklist
		r.Group(func(r chi.Router) {
			r.Post("/{id}/checklist", handler.AddChecklistItem)
			r.Delete("/{id}/checklist/{itemId}", handler.DeleteChecklistItem)
			r.Put("/{id}/checklist/{itemId}", handler.UpdateChecklistItem)
			r.Patch("/{id}/checklist/{itemId}", handler.UpdateChecklistItemStatus)
		})
	})
}
