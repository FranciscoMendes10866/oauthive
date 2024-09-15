package repository

import (
	"context"
	"oauthive/db/entities"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
)

type SessionRepository struct {
	db rel.Repository
}

func NewSessionRepository(db rel.Repository) *SessionRepository {
	return &SessionRepository{db}
}

func (repo *SessionRepository) CreateSession(ctx context.Context, session *entities.Session) (entities.Session, error) {
	if err := repo.db.Insert(ctx, session); err != nil {
		if err != rel.ErrUniqueConstraint {
			return entities.Session{}, err
		}
		var existingSession entities.Session
		err := repo.db.Find(
			ctx,
			&existingSession,
			where.Eq("user_id", session.UserID).AndEq("expires_at", session.ExpiresAt),
		)
		if err != nil {
			return entities.Session{}, err
		}
		return existingSession, nil
	}
	return *session, nil
}

func (repo *SessionRepository) DeleteSessionByID(ctx context.Context, sessionID int) error {
	return repo.db.Delete(ctx, &entities.Session{ID: sessionID})
}

func (repo *SessionRepository) FindSessionByID(ctx context.Context, sessionID int) (entities.Session, error) {
	session := entities.Session{}
	if err := repo.db.Find(ctx, &session, where.Eq("id", sessionID)); err != nil {
		return session, err
	}
	return session, nil
}
