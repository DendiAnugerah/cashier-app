package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SessionsRepository struct {
	db *gorm.DB
}

func NewSessionsRepository(db *gorm.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	result := u.db.Create(&session)
	return result.Error
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	result := u.db.Where("token = ?", tokenTarget).Delete(&model.Session{})
	return result.Error
}

func (u *SessionsRepository) UpdateSessions(session model.Session) error {
	result := u.db.Where("username = ?", session.Username).Updates(&session)
	return result.Error
}

func (u *SessionsRepository) TokenValidity(token string) (model.Session, error) {
	var session model.Session
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}
	if u.TokenExpired(session) {
		return model.Session{}, errors.New("token expired")
	}
	return session, nil
}

func (u *SessionsRepository) SessionAvailName(name string) (model.Session, error) {
	var session model.Session
	availName := u.db.Where("username = ?", name).First(&session)
	if availName.Error != nil {
		return model.Session{}, errors.New("username not found")
	}

	return session, nil
}

func (u *SessionsRepository) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session
	result := u.db.Where("token = ?", token).First(&session)
	if result.Error != nil {
		return model.Session{}, errors.New("token not found")
	}

	return session, nil
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
