package service

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/fayupable/chessfut-be/domain"
)

type mockChessComClient struct {
	mock.Mock
}

func (m *mockChessComClient) GetProfile(ctx context.Context, username string) (domain.Player, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.Player), args.Error(1)
}

func (m *mockChessComClient) GetStats(ctx context.Context, username string) (domain.PlayerStats, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.PlayerStats), args.Error(1)
}

func (m *mockChessComClient) GetGames(ctx context.Context, username string, from, to time.Time) ([]domain.Game, error) {
	args := m.Called(ctx, username, from, to)
	return args.Get(0).([]domain.Game), args.Error(1)
}

func (m *mockChessComClient) GetTitledUsernames(ctx context.Context, title domain.Title) ([]string, error) {
	args := m.Called(ctx, title)
	return args.Get(0).([]string), args.Error(1)
}

type mockCardRepository struct {
	mock.Mock
}

func (m *mockCardRepository) Save(ctx context.Context, card domain.Card) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}

func (m *mockCardRepository) FindByUsername(ctx context.Context, username string) (domain.Card, bool, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.Card), args.Bool(1), args.Error(2)
}

func (m *mockCardRepository) FindStale(ctx context.Context, before time.Time, limit int) ([]domain.Card, error) {
	args := m.Called(ctx, before, limit)
	return args.Get(0).([]domain.Card), args.Error(1)
}

func (m *mockCardRepository) FindTopByOVR(ctx context.Context, limit, offset int) ([]domain.Card, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Card), args.Error(1)
}

type mockCache struct {
	mock.Mock
}

func (m *mockCache) GetCard(ctx context.Context, username string) (domain.Card, bool, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.Card), args.Bool(1), args.Error(2)
}

func (m *mockCache) SetCard(ctx context.Context, card domain.Card) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}

func (m *mockCache) IncrementViewCount(ctx context.Context, username string) (int, error) {
	args := m.Called(ctx, username)
	return args.Int(0), args.Error(1)
}

func (m *mockCache) IsPromoted(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *mockCardRepository) CountAll(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}
func (m *mockCache) GetTotalCardsCount(ctx context.Context) (int, bool, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Bool(1), args.Error(2)
}

func (m *mockCache) SetTotalCardsCount(ctx context.Context, count int) error {
	args := m.Called(ctx, count)
	return args.Error(0)
}

func (m *mockCardRepository) SearchByUsername(ctx context.Context, prefix string, limit, offset int) ([]domain.Card, error) {
	args := m.Called(ctx, prefix, limit, offset)
	return args.Get(0).([]domain.Card), args.Error(1)
}
