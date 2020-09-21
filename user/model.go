package user

import (
	"github.com/google/uuid"
	"time"
)

type user struct {
	id uuid.UUID
	stripeId string
	email string
	createdAt time.Time
}
