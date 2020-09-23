package user

import (
	"github.com/google/uuid"
	"time"
)

type user struct {
	id uuid.UUID
	stripeId string
	linkId string
	email string
	status string
	createdAt time.Time
}
