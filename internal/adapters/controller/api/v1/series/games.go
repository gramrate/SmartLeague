package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) CreateGame(c echo.Context) error {
	// Create game in series
	//
	// @Summary Create game in series
	// @Tags game
	// @Accept json
	// @Produce json
	// @Param id path string true "Series ID"
	// @Param request body dto.CreateGameRequest true "Game data"
	// @Success 201 {object} dto.CreateGameResponse
	// @Failure 400 {object} dto.HTTPStatus
	// @Failure 401 {object} dto.HTTPStatus
	// @Failure 403 {object} dto.HTTPStatus
	// @Failure 500 {object} dto.HTTPStatus
	// @Router /api/v1/series/{id}/games [post]
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.CreateGameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.SeriesID = seriesID
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.gameService.Create(c.Request().Context(), requesterID, &req)
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *handler) GetSeriesGames(c echo.Context) error {
	// Get games in series (paginated)
	//
	// @Summary Get games in series
	// @Tags game
	// @Produce json
	// @Param id path string true "Series ID"
	// @Param limit query int false "limit"
	// @Param offset query int false "offset"
	// @Success 200 {object} dto.GetSeriesGamesResponse
	// @Failure 400 {object} dto.HTTPStatus
	// @Failure 403 {object} dto.HTTPStatus
	// @Failure 500 {object} dto.HTTPStatus
	// @Router /api/v1/series/{id}/games [get]
	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.GetSeriesGamesRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.SeriesID = seriesID
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.gameService.ListBySeries(c.Request().Context(), maybeRequesterID(c.Get("user_id")), &req)
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *handler) UpdateGame(c echo.Context) error {
	// Update game by id
	//
	// @Summary Update game by id
	// @Tags game
	// @Accept json
	// @Produce json
	// @Param id path string true "Game ID"
	// @Param request body dto.UpdateGameRequest true "Update data"
	// @Success 200 {object} dto.UpdateGameResponse
	// @Failure 400 {object} dto.HTTPStatus
	// @Failure 401 {object} dto.HTTPStatus
	// @Failure 403 {object} dto.HTTPStatus
	// @Failure 500 {object} dto.HTTPStatus
	// @Router /api/v1/game/{id} [patch]
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.UpdateGameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.ID = gameID
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.gameService.Update(c.Request().Context(), requesterID, &req)
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *handler) SetParticipants(c echo.Context) error {
	// Set game participants
	//
	// @Summary Set game participants
	// @Tags game
	// @Accept json
	// @Produce json
	// @Param id path string true "Game ID"
	// @Param request body dto.SetGameParticipantsRequest true "Participant IDs"
	// @Success 204
	// @Failure 400 {object} dto.HTTPStatus
	// @Failure 401 {object} dto.HTTPStatus
	// @Failure 403 {object} dto.HTTPStatus
	// @Failure 500 {object} dto.HTTPStatus
	// @Router /api/v1/game/{id}/participants [post]
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.SetGameParticipantsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.GameID = gameID
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := h.gameService.SetParticipants(c.Request().Context(), requesterID, &req); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *handler) UpsertResults(c echo.Context) error {
	// Upsert game results
	//
	// @Summary Upsert game results
	// @Tags game
	// @Accept json
	// @Produce json
	// @Param id path string true "Game ID"
	// @Param request body dto.UpsertGameResultsRequest true "Results"
	// @Success 204
	// @Failure 400 {object} dto.HTTPStatus
	// @Failure 401 {object} dto.HTTPStatus
	// @Failure 403 {object} dto.HTTPStatus
	// @Failure 500 {object} dto.HTTPStatus
	// @Router /api/v1/game/{id}/results [post]
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.UpsertGameResultsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.GameID = gameID
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := h.gameService.UpsertResults(c.Request().Context(), requesterID, &req); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *handler) GetLeaderboard(c echo.Context) error {
	// Get series leaderboard (paginated)
	//
	// @Summary Get series leaderboard
	// @Tags series
	// @Produce json
	// @Param id path string true "Series ID"
	// @Param limit query int false "limit"
	// @Param offset query int false "offset"
	// @Success 200 {object} dto.GetSeriesLeaderboardResponse
	// @Failure 400 {object} dto.HTTPStatus
	// @Failure 403 {object} dto.HTTPStatus
	// @Failure 500 {object} dto.HTTPStatus
	// @Router /api/v1/series/{id}/leaderboard [get]
	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.GetSeriesLeaderboardRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.SeriesID = seriesID
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.seriesService.GetLeaderboard(c.Request().Context(), maybeRequesterID(c.Get("user_id")), &req)
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
