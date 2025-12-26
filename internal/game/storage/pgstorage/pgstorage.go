package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGStorage struct {
	db *pgxpool.Pool
}

func NewPGStorage(connString string) (*PGStorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}

	storage := &PGStorage{
		db: db,
	}

	err = storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGStorage) initTables() error {
	roomsSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s INTEGER PRIMARY KEY,
		%s VARCHAR(50) NOT NULL,
		%s TEXT NOT NULL
	)`, roomsTable, roomNumberColumn, nameColumn, descriptionColumn)

	_, err := s.db.Exec(context.Background(), roomsSQL)
	if err != nil {
		return errors.Wrap(err, "создание таблицы rooms")
	}

	gameSessionsSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s UUID PRIMARY KEY,
		%s UUID NOT NULL,
		%s UUID NOT NULL,
		%s INTEGER DEFAULT 1 NOT NULL
	)`, gameSessionsTable, partnershipIDColumn,
		user1IDColumn, user2IDColumn, currentRoomColumn)

	_, err = s.db.Exec(context.Background(), gameSessionsSQL)
	if err != nil {
		return errors.Wrap(err, "создание таблицы game_sessions")
	}

	indexesSQL := fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS idx_game_sessions_user1 ON %s(%s);
	CREATE INDEX IF NOT EXISTS idx_game_sessions_user2 ON %s(%s);
	`, gameSessionsTable, user1IDColumn,
		gameSessionsTable, user2IDColumn)

	_, err = s.db.Exec(context.Background(), indexesSQL)
	if err != nil {
		return errors.Wrap(err, "создание индексов")
	}

	return nil
}

func (s *PGStorage) InitializeRooms(ctx context.Context) error {
	var count int
	err := s.db.QueryRow(ctx, 
		fmt.Sprintf("SELECT COUNT(*) FROM %s", roomsTable)).Scan(&count)
	if err != nil {
		return errors.Wrap(err, "ошибка проверки количества комнат")
	}

	if count > 0 {
		return nil
	}

	rooms := []struct {
		number      int32
		name        string
		description string
	}{
		{1, "Прихожая", "Вы в тёмной прихожей старого дома."},
		{2, "Кабинет", "Старый кабинет с пыльными книгами."},
		{3, "Библиотека", "Пыльная библиотека с древними фолиантами."},
		{4, "Гостиная", "Простораная гостиная с камином."},
		{5, "Спальня", "Старинная спальня с кованой кроватью."},
	}

	for _, room := range rooms {
		query := fmt.Sprintf(`
			INSERT INTO %s (%s, %s, %s) 
			VALUES ($1, $2, $3)
			ON CONFLICT (%s) DO NOTHING`,
			roomsTable, roomNumberColumn, nameColumn, descriptionColumn, roomNumberColumn)

		_, err := s.db.Exec(ctx, query, room.number, room.name, room.description)
		if err != nil {
			return errors.Wrapf(err, "ошибка добавления комнаты %d", room.number)
		}
	}

	return nil
}

func (s *PGStorage) Close() {
	s.db.Close()
}