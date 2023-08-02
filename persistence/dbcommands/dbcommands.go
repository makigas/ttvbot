package dbcommands

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func NewDbCommands(rdb *redis.Client) *DbCommands {
	return &DbCommands{rdb: rdb}
}

type DbCommands struct {
	rdb *redis.Client
}

var ErrCommandNotFound = errors.New("command not found")

const commandPrefix = "cmd:"

func commandKey(name string) string {
	return commandPrefix + name
}

func (dbc *DbCommands) ListCommands() ([]string, error) {
	res, err := dbc.rdb.Keys(context.Background(), "cmd:*").Result()
	if err != nil {
		return nil, err
	}
	result := make([]string, len(res))
	for i, r := range res {
		result[i], _ = strings.CutPrefix(commandPrefix, r)
	}
	return result, nil
}

func (dbc *DbCommands) GetCommand(name string) (string, error) {
	res, err := dbc.rdb.Get(context.Background(), commandKey(name)).Result()
	if err == redis.Nil {
		return "", ErrCommandNotFound
	}
	if err != nil {
		return "", err
	}
	return res, nil
}

func (dbc *DbCommands) SetCommand(name, value string) error {
	return dbc.rdb.Set(context.Background(), commandKey(name), value, 0).Err()
}

func (dbc *DbCommands) DeleteCommand(name string) error {
	return dbc.rdb.Del(context.Background(), commandKey(name)).Err()
}

type DbCommandsResult struct {
	fx.Out
	Commands *DbCommands
}

type DbCommandsParams struct {
	fx.In
	Rdb *redis.Client
}

func dbCommand(par *DbCommandsParams) *DbCommandsResult {
	return &DbCommandsResult{Commands: NewDbCommands(par.Rdb)}
}

var Module = fx.Module("dbcommands", fx.Provide(dbCommand))
