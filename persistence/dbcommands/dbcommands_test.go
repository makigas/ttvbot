package dbcommands

import (
	"testing"

	"github.com/go-redis/redismock/v9"
)

func TestDbCommandsGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("cmd:socials").RedisNil()
	mock.ExpectGet("cmd:youtube").SetVal("http://youtube.com")

	cmd := NewDbCommands(db)
	if _, err := cmd.GetCommand("socials"); err != ErrCommandNotFound {
		t.Errorf("Expected !socials to fail with 'CommandNotFound'")
	}

	res, err := cmd.GetCommand("youtube")
	if err != nil {
		t.Error(err)
	}
	if res != "http://youtube.com" {
		t.Errorf("!youtube yielded wrong value: %s", res)
	}
}

func TestDbCommandsSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectSet("cmd:social", "http://linktree.com", 0).SetVal("OK")

	cmd := NewDbCommands(db)
	if err := cmd.SetCommand("social", "http://linktree.com"); err != nil {
		t.Error(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDbCommandsDelete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectDel("cmd:social").SetVal(1)

	cmd := NewDbCommands(db)
	if err := cmd.DeleteCommand("social"); err != nil {
		t.Error(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
