package data

import (
	"github.com/jinzhu/gorm"
)

const (
	UnknownSource int = iota
	LocalPassword
)

type CredentialsDAO struct {
	dao *DAO
	db  *gorm.DB
}

func NewCredentialsDAO(dao *DAO, db *gorm.DB) CredentialsDAO {
	return CredentialsDAO{dao, db}
}

type Credential struct {
	ID     int    `gorm:"primary_key" json:"id"`
	UserID int    `json:"user_id"`
	Record string `json:"record,omitempty"`
	Source int    `json:"source"`
}

func (d *CredentialsDAO) GetByUser(userId int) ([]Credential, error) {
	t := make([]Credential, 0)
	err := d.db.Where(&Credential{UserID: userId}).Find(&t).Error

	for i := range t {
		if t[i].Source == LocalPassword {
			t[i].Record = ""
		}
	}
	return t, err
}

func (d *CredentialsDAO) Delete(id int) error {
	err := d.db.Delete(&Credential{}, id).Error
	logError(err)
	return err
}

func (d *CredentialsDAO) Save(u *Credential) error {
	var err error
	if u.ID != 0 {
		// update record field only
		err = d.db.Model(u).Update(Credential{Record: u.Record}).Error
	} else {
		err = d.db.Save(u).Error
	}

	logError(err)
	return err
}
