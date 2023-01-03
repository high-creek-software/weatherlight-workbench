package storage

import "fyne.io/fyne/v2"

type SetIconRepo interface {
	Image(code string) fyne.Resource
}

type SetIconRepoImpl struct {
}