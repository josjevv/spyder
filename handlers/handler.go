package handler

import (
	db "github.com/changer/spyder/db"
)

type Handler func(fly *db.Fly)
