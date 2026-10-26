package dto

import "library-api/internal/model"

type UserDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func UserModelToDTO(u *model.User) UserDTO {
	return UserDTO{
		ID:   u.ID,
		Name: u.Name,
	}
}
