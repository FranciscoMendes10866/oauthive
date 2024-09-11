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

func (repo *SessionRepository) CreateSession(ctx context.Context, session *entities.Session) (*entities.Session, error) {
	err := repo.db.Insert(ctx, session)
	if err == nil {
		return nil, nil
	}
	existingSession := &entities.Session{}
	err = repo.db.Find(
		ctx,
		existingSession,
		rel.Eq("user_id", session.UserID).AndEq("expires_at", session.ExpiresAt),
	)
	if err == rel.ErrNotFound {
		return nil, err
	}
	return existingSession, nil
}

func (repo *SessionRepository) DeleteSessionByID(ctx context.Context, sessionID int) error {
	return repo.db.Delete(ctx, &entities.Session{ID: sessionID})
}

func (repo *SessionRepository) FindSessionByID(ctx context.Context, sessionID int) (*entities.Session, error) {
	session := &entities.Session{}
	err := repo.db.Find(ctx, session, where.Eq("id", sessionID))
	if err != nil {
		return nil, err
	}
	return session, nil
}
