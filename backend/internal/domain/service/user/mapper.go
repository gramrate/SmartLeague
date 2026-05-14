package user

import (
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
)

func toDTO(u *model.User) dto.User {
	return dto.User{
		ID:          u.ID,
		Nickname:    u.Nickname,
		Name:        u.Name,
		ShowName:    u.ShowName,
		Description: u.Description,
		Email:       u.Email,
		ClubID:      u.ClubID,
		ClubState:   u.ClubState,
		Role:        u.Role,
	}
}
