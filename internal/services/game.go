package services

import (
	"context"
	"mine-game/internal/model"
	repository "mine-game/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameService struct {
	gameRepository *repository.GameRepository
}

// NewGameService creates a new game service.
func NewGameService(gameRepository *repository.GameRepository) *GameService {
	return &GameService{
		gameRepository: gameRepository,
	}
}

// CreateNewGame creates a new game.
func (s *GameService) CreateNewGame(ctx context.Context, mines int) (*model.Game, error) {
	game, err := model.NewGame(mines)
	if err != nil {
		return nil, err
	}

	err = s.gameRepository.CreateNewGame(ctx, game)
	return game, err
}

// GetGame retrieves a game by ID.
func (s *GameService) GetGame(ctx context.Context, id primitive.ObjectID) (*model.Game, error) {
	game, err := s.gameRepository.GetGame(ctx, id)
	return game, err
}

// UpdateGame updates a game.
func (s *GameService) UpdateGame(ctx context.Context, game *model.Game) error {
	return s.gameRepository.UpdateGame(ctx, game)
}
