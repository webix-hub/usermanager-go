package data

import 	"github.com/jinzhu/gorm"

type Link [2]int
type DBLink struct {
	A int
	B int
}

type Links struct {
	RoleRule []Link
	UserRole []Link
	UserRule []Link
}

type MetaDAO struct {
	dao *DAO
	db  *gorm.DB
}

func NewMetaDAO(dao *DAO, db *gorm.DB) MetaDAO {
	return MetaDAO{dao, db}
}

func (d *MetaDAO) GetLinks() (*Links, error) {
	var sql string
	var err error

	sql = `select role_id as a, rule_id as b from role_rules`
	rr := make([]DBLink, 0)
	err = d.db.Raw(sql).Scan(&rr).Error
	if err != nil {
		return nil, err
	}

	sql = `select user_id as a, role_id as b from user_roles`
	uro := make([]DBLink, 0)
	err = d.db.Raw(sql).Scan(&uro).Error
	if err != nil {
		return nil, err
	}

	sql = `select user_id as a, rule_id as b from user_rules`
	uru := make([]DBLink, 0)
	err = d.db.Raw(sql).Scan(&uru).Error
	if err != nil {
		return nil, err
	}

	return &Links{relink(rr), relink(uro), relink(uru)}, nil
}

func relink(source []DBLink) []Link {
	t := make([]Link, len(source))
	for i := range source {
		t[i][0] = source[i].A
		t[i][1] = source[i].B
	}

	return t
}
