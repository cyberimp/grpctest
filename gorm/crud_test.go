package gorm

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strconv"
	"testing"
	"time"
)

var (
	testDb *gorm.DB
	mockDb sqlmock.Sqlmock
)

//AnyTime for matching any timestamps
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

//init creates mock db
func init() {
	var (
		db  *sql.DB
		err error
	)
	db, mockDb, err = sqlmock.New()
	if err != nil {
		log.Fatalf("error creating db: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:        "sqlmock_db_0",
		DriverName: "postgres",
		Conn:       db,

		PreferSimpleProtocol: true,
	})
	testDb, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("error creating db: %v", err)
	}
}

func TestConn_Create(t *testing.T) {
	tests := []struct {
		name    string
		newPost *Post
	}{
		{name: "should add post", newPost: &Post{Title: "test", Text: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Conn{db: testDb}
			mockDb.NewRows([]string{"id", "deleted_at"})
			mockDb.ExpectBegin()
			mockDb.ExpectQuery("^INSERT INTO \"posts\" (.+)$").WillReturnRows(mockDb.NewRows([]string{"", ""}))
			mockDb.ExpectCommit()
			_, err := s.Create(tt.newPost)
			if err != nil {
				t.Errorf("Create() got error %q", err)
			}
			got := mockDb.ExpectationsWereMet()
			if got != nil {
				t.Errorf("Create() got error %q", got)
			}
		})
	}
}

func TestConn_Delete(t *testing.T) {
	tests := []struct {
		name string
		id   uint
	}{
		{name: "should delete post", id: 1},
		{name: "should another post", id: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Conn{db: testDb}
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(
				`UPDATE "posts"  SET "deleted_at"=$1 WHERE "posts"."id" = $2 AND "posts"."deleted_at" IS NULL`)).
				WithArgs(AnyTime{}, tt.id).WillReturnResult(sqlmock.NewResult(1, 1))
			mockDb.ExpectCommit()
			err := s.Delete(tt.id)
			if err != nil {
				t.Errorf("Delete() got error %q", err)
			}
			got := mockDb.ExpectationsWereMet()
			if got != nil {
				t.Errorf("Delete() got error %q", got)
			}
		})
	}
}

func TestConn_List(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		size      int
		wantError bool
	}{
		{"Should list with page 1, size 1", 1, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Conn{db: testDb}
			mockDb.ExpectQuery(
				regexp.QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."deleted_at" IS NULL LIMIT ` +
					strconv.Itoa(tt.size) + ` OFFSET ` + strconv.Itoa(tt.page*tt.size))).
				WillReturnRows(sqlmock.NewRows([]string{"", ""}))
			_, err := s.List(tt.page, tt.size)
			if (err != nil) != tt.wantError {
				t.Errorf("List got unwanted error:%q", err)
			}
			got := mockDb.ExpectationsWereMet()
			if got != nil {
				t.Errorf("List() expectations where not met %q", got)
			}
		})
	}
}

func TestConn_Read(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{"should read row with id 0", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Conn{db: testDb}
			mockRow := mockDb.NewRows([]string{"id", "created_at", "updated_at", "name", "text"}).
				AddRow(tt.id, time.Now(), time.Now(), "1", "2")
			mockDb.ExpectQuery(regexp.QuoteMeta(
				`SELECT * FROM "posts" WHERE "posts"."id" = $1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT 1`)).
				WithArgs(tt.id).WillReturnRows(mockRow)
			_, err := s.Read(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = mockDb.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Read() error %q", err)
			}
		})
	}
}

func TestConn_Update(t *testing.T) {
	tests := []struct {
		name    string
		post    *Post
		newPost *Post
		wantErr bool
	}{
		{"Should update post",
			&Post{gorm.Model{ID: 1}, "name_old", "value_old"},
			&Post{gorm.Model{ID: 1}, "name", "value"},
			false},
		{"Should not update post",
			&Post{gorm.Model{ID: 1}, "name_old", "value_old"},
			&Post{gorm.Model{ID: 2}, "name", "value"},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Conn{db: testDb}

			mockRow := mockDb.NewRows([]string{"id", "created_at", "updated_at", "title", "text"}).
				AddRow(tt.post.ID, time.Now(), time.Now(), tt.post.Title, tt.post.Text)

			mockDb.ExpectQuery(regexp.
				QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."id" = $1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT 1`)).
				WithArgs(tt.newPost.ID).WillReturnRows(mockRow)

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.
				QuoteMeta(`UPDATE "posts" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"title"=$4,"text"=$5 WHERE "id" = $6`)).
				WithArgs(AnyTime{}, AnyTime{}, nil, tt.newPost.Title, tt.newPost.Text, tt.newPost.ID).
				WillReturnResult(sqlmock.NewResult(int64(tt.newPost.ID), 1))
			mockDb.ExpectCommit()
			_, err := s.Update(tt.newPost)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error %q", err)
				return
			}

			err = mockDb.ExpectationsWereMet()
			if (err != nil) != tt.wantErr {
				t.Errorf("expectations were unmet, error: %q", err)
			}
		})
	}
}
