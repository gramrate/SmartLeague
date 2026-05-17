package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetSeriesFull Get series with participants, games and leaderboard
//
// @Summary Get series full data
// @Tags series
// @Produce json
// @Param id path string true "Series ID"
// @Param participants_limit query int false "participants limit"
// @Param participants_offset query int false "participants offset"
// @Param games_limit query int false "games limit"
// @Param games_offset query int false "games offset"
// @Param leaderboard_limit query int false "leaderboard limit"
// @Param leaderboard_offset query int false "leaderboard offset"
// @Success 200 {object} dto.GetSeriesFullResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Router /api/v1/series/{id}/full [get]
func (h *handler) GetSeriesFull(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.GetSeriesFullRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.ID = id
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	requesterID := maybeRequesterID(c.Get("user_id"))

	seriesResp, err := h.seriesService.GetSeries(c.Request().Context(), requesterID, &dto.GetSeriesRequest{ID: req.ID})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	participantsResp, err := h.seriesService.GetParticipants(c.Request().Context(), requesterID, &dto.GetSeriesParticipantsRequest{
		SeriesID: req.ID,
		Limit:    req.ParticipantsLimit,
		Offset:   req.ParticipantsOffset,
	})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	gamesResp, err := h.gameService.ListBySeries(c.Request().Context(), requesterID, &dto.GetSeriesGamesRequest{
		SeriesID: req.ID,
		Limit:    req.GamesLimit,
		Offset:   req.GamesOffset,
	})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	leaderboardResp, err := h.seriesService.GetLeaderboard(c.Request().Context(), requesterID, &dto.GetSeriesLeaderboardRequest{
		SeriesID: req.ID,
		Limit:    req.LeaderboardLimit,
		Offset:   req.LeaderboardOffset,
	})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &dto.GetSeriesFullResponse{
		Series:       (*dto.Series)(seriesResp),
		Participants: participantsResp,
		Games:        gamesResp,
		Leaderboard:  leaderboardResp,
	})
}
