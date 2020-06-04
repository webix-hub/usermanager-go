package data

import (
	"github.com/jinzhu/gorm"
)

type RolesDAO struct {
	dao *DAO
	db  *gorm.DB
}

func NewRolesDAO(dao *DAO, db *gorm.DB) RolesDAO {
	return RolesDAO{dao, db}
}

type Role struct {
	ID      int    `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Details string `json:"details"`
}

type RoleRule struct {
	ID     int `gorm:"primary_key"`
	RoleID int
	RuleID int
}

func (d *RolesDAO) GetAll() ([]Role, error) {
	t := make([]Role, 0)
	err := d.db.Find(&t).Error
	logError(err)

	return t, err
}

func (d *RolesDAO) GetOne(id int) (*Role, error) {
	t := Role{}
	err := d.db.Find(&t, id).Error
	logError(err)
	return &t, err
}

func (d *RolesDAO) Delete(id int) error {
	d.LogRoleDeleted(id)
	err := d.db.Delete(&Role{}, id).Error
	if err != nil {
		logError(err)
		return err
	}
	err = d.RemoveAllRoleRules(id)
	if err != nil {
		return err
	}
	return nil
}

func (d *RolesDAO) RemoveAllRoleRules(roleId int) error {
	err := d.db.Delete(RoleRule{}, "role_id = ?", roleId).Error
	logError(err)
	return err
}

func (d *RolesDAO) Save(u *Role) error {
	isNewRole := false
	if u.ID == 0 {
		isNewRole = true
	}

	err := d.db.Save(u).Error

	if isNewRole {
		d.LogRoleAdded(u)
	}

	logError(err)
	return err
}

func (d *RolesDAO) GetRules(rid int) ([]int, error) {
	rows, err := d.db.Model(&RoleRule{}).Select("rule_id").Where(&RoleRule{RoleID: rid}).Rows()
	if err != nil {
		logError(err)
		return nil, err
	}

	return rowsToIntSlice(rows)
}

func (d *RolesDAO) AddRule(u *Role, r int) {
	d.LogRuleAdded(u, r)
	err := d.db.Save(&RoleRule{RoleID: u.ID, RuleID: r}).Error
	logError(err)
}

func (d *RolesDAO) RemoveRule(u *Role, r int) {
	if u.ID == 0 || r == 0 {
		return
	}

	d.LogRuleRemoved(u, r)
	obj := RoleRule{RoleID: u.ID, RuleID: r}
	err := d.db.Where(&obj).Delete(&obj).Error
	logError(err)
}
