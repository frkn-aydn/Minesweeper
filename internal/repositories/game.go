package repository

import (
	"context"
	"mine-game/internal/database"
	"mine-game/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameRepository struct {
	db *database.MongoDB
}

// NewGameRepository creates a new game repository.
func NewGameRepository(db *database.MongoDB) *GameRepository {
	return &GameRepository{
		db: db,
	}
}

// CreateNewGame creates a new game.
func (r *GameRepository) CreateNewGame(ctx context.Context, game *model.Game) error {
	_, err := r.db.GameCollection().InsertOne(ctx, game)
	return err
}

// GetGame retrieves a game by ID.
func (r *GameRepository) GetGame(ctx context.Context, id primitive.ObjectID) (*model.Game, error) {
	var game model.Game
	err := r.db.GameCollection().FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&game)

	return &game, err
}

// UpdateGame updates a game.
func (r *GameRepository) UpdateGame(ctx context.Context, game *model.Game) error {
	filter := bson.M{
		"_id": game.ID,
	}

	update := bson.M{
		"$set": game,
	}

	_, err := r.db.GameCollection().UpdateOne(ctx, filter, update)
	return err
}
