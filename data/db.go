package data

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"log"
)

var Debug = true

func logSQL(s string) {
	if Debug {
		log.Printf("sql: %s", s)
	}
}

func logError(e error) {
	if e != nil {
		log.Printf("[ERROR]\n%s\n", e)
	}
}

type DAO struct {
	db *gorm.DB

	Credentials CredentialsDAO
	Users       UsersDAO
	Roles       RolesDAO
	Meta        MetaDAO
	Logs        LogsDAO
}

func (d *DAO) GetDB() *gorm.DB {
	return d.db
}

func NewDAO(db *gorm.DB) *DAO {
	d := DAO{}
	d.db = db
	d.Logs = NewLogsDAO(&d, db)
	d.Credentials = NewCredentialsDAO(&d, db)
	d.Users = NewUsersDAO(&d, db)
	d.Roles = NewRolesDAO(&d, db)
	d.Meta = NewMetaDAO(&d, db)


	d.db.AutoMigrate(&User{}, &UserRole{}, &UserRule{})
	d.db.AutoMigrate(&Role{}, &RoleRule{})
	d.db.AutoMigrate(&Credential{})
	d.db.AutoMigrate(&Log{})
	return &d
}

func rowsToIntSlice(rows *sql.Rows) ([]int, error){
	var temp int
	out := make([]int, 0)

	for rows.Next() {
		if err := rows.Scan(&temp); err != nil {
			return nil, err
		}
		out = append(out, temp)
	}

	return out, nil
}