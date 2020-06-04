package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UsersDAO struct {
	dao *DAO
	db  *gorm.DB
}

func NewUsersDAO(dao *DAO, db *gorm.DB) UsersDAO {
	return UsersDAO{dao, db}
}

type User struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Details    string    `json:"details"`
	Visited    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"visited"`
	Registered time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"registered"`
	Avatar     string    `json:"avatar"`
	Status     int       `json:"status"`
}

type UserRole struct {
	ID     int `gorm:"primary_key"`
	UserID int
	RoleID int
}

type UserRule struct {
	ID     int `gorm:"primary_key"`
	UserID int
	RuleID int
}

func (d *UsersDAO) GetAll() ([]User, error) {
	t := make([]User, 0)
	err := d.db.Find(&t).Error
	logError(err)
	return t, err
}

func (d *UsersDAO) GetOne(id int) (*User, error) {
	t := User{}
	err := d.db.Find(&t, id).Error
	logError(err)
	return &t, err
}

func (d *UsersDAO) Delete(id int) error {
	d.LogRemoveUser(id)
	err := d.db.Delete(&User{}, id).Error
	d.RemoveAllRoleRules(id)
	logError(err)
	return err
}

func (d *UsersDAO) RemoveAllRoleRules(userId int) error {
	err := d.db.Delete(UserRole{}, "user_id = ?", userId).Error
	if err != nil {
		logError(err)
		return err
	}
	err = d.db.Delete(UserRule{}, "user_id = ?", userId).Error
	if err != nil {
		logError(err)
		return err
	}
	return nil
}

func (d *UsersDAO) Save(u *User) error {
	isNewUser := false
	if u.ID == 0 {
		isNewUser = true
	}

	err := d.db.Save(u).Error

	if isNewUser {
		d.LogAddUser(u)
	}

	logError(err)
	return err
}

func (d *UsersDAO) GetRoles(uid int) ([]int, error) {
	rows, err := d.db.Model(&UserRole{}).Select("role_id").Where(&UserRole{UserID: uid}).Rows()
	if err != nil {
		logError(err)
		return nil, err
	}

	return rowsToIntSlice(rows)
}

func (d *UsersDAO) AddRole(u *User, r int) {
	d.LogAddRole(u, r)
	err := d.db.Save(&UserRole{UserID: u.ID, RoleID: r}).Error
	logError(err)
}

func (d *UsersDAO) RemoveRole(u *User, r int) {
	if u.ID == 0 || r == 0 {
		return
	}

	d.LogRemoveRole(u, r)
	obj := UserRole{UserID: u.ID, RoleID: r}
	err := d.db.Where(&obj).Delete(&obj).Error
	logError(err)
}

func (d *UsersDAO) GetRules(uid int) ([]int, error) {
	rows, err := d.db.Model(&UserRule{}).Select("rule_id").Where(&UserRule{UserID: uid}).Rows()
	if err != nil {
		logError(err)
		return nil, err
	}

	return rowsToIntSlice(rows)
}

func (d *UsersDAO) AddRule(u *User, r int) {
	d.LogAddRule(u, r)
	err := d.db.Save(&UserRule{UserID: u.ID, RuleID: r}).Error
	logError(err)
}

func (d *UsersDAO) RemoveRule(u *User, r int) {
	if u.ID == 0 || r == 0 {
		return
	}

	d.LogRemoveRule(u, r)
	obj := UserRule{UserID: u.ID, RuleID: r}
	err := d.db.Where(&obj).Delete(&obj).Error
	logError(err)
}
