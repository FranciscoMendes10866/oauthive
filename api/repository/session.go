package repository

import (
	"context"
	"oauthive/db/entities"

	"github.com/go-rel/rel"
)

type SessionRepository struct {
	db rel.Repository
}

func NewSessionRepository(db rel.Repository) *SessionRepository {
	return &SessionRepository{db}
}

func (self *SessionRepository) CreateSession(ctx context.Context, session *entities.Session) (*entities.Session, error) {
	err := self.db.Insert(ctx, session)
	if err != nil {
		existingSession := &entities.Session{}
		err = self.db.Find(
			ctx,
			existingSession,
			rel.Eq("user_id", session.UserID).
				AndEq("expires_at", session.ExpiresAt),
		)
		if err == rel.ErrNotFound {
			return nil, err
		}

		return existingSession, nil
	}

	return nil, err
}

func (self *SessionRepository) DeleteSessionByID(ctx context.Context, sessionID int) error {
	return self.db.Delete(ctx, &entities.Session{ID: sessionID})
}
