package data

import (
	"fmt"
	"time"
	"webix/usermanager/rules"

	"github.com/jinzhu/gorm"
)

type LogType int

const (
	UnknownLogType LogType = iota
	LogLogin
	LogRoleAdded
	LogRoleChanged
	LogRoleDeleted
	LogUserAdded
	LogUserChanged
	LogUserDeleted
)

type LogsDAO struct {
	dao *DAO
	db  *gorm.DB
}

func NewLogsDAO(dao *DAO, db *gorm.DB) LogsDAO {
	return LogsDAO{dao, db}
}

type Log struct {
	ID       int       `gorm:"primary_key" json:"id"`
	UserID   int       `json:"user_id"`
	TargetID int       `json:"target_id"`
	Type     LogType   `json:"type"`
	Details  string    `json:"details"`
	Date     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
}

type LogTypeDescription struct {
	Name   string `json:"name"`
	Target string `json:"target"`
}

func (d *LogsDAO) GetBy(userId, targetId int, types []LogType) ([]Log, error) {
	query := d.db
	if targetId != 0 {
		query = query.Where("target_id = ?", targetId)
	}
	if userId != 0 {
		query = query.Where("user_id = ?", userId)
	}
	if types != nil && len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	t := make([]Log, 0)
	err := query.Find(&t).Error
	return t, err
}

func (d *LogsDAO) Add(l *Log) error {
	l.ID = 0
	l.Date = time.Now()

	err := d.db.Save(l).Error
	logError(err)
	return err
}

func (d *UsersDAO) LogAddUser(u *User) error {
	logRecord := Log{UserID: 0, TargetID: u.ID, Type: LogUserAdded, Details: u.Name}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *UsersDAO) LogRemoveUser(id int) error {
	u := User{}
	d.db.Where(&User{ID: id}).First(&u)
	logRecord := Log{UserID: 0, TargetID: u.ID, Type: LogUserDeleted, Details: u.Name}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *UsersDAO) LogAddEmail(u *User, email string) error {
	logDetails := fmt.Sprintf("email %s changed to %s", u.Email, email)
	logRecord := Log{UserID: 0, TargetID: u.ID, Type: LogUserChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *UsersDAO) LogAddRole(u *User, r int) error {
	role := Role{}
	d.db.Where(&Role{ID: r}).First(&role)
	logDetails := fmt.Sprintf("role %s added", role.Name)
	logRecord := Log{UserID: 0, TargetID: u.ID, Type: LogUserChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *UsersDAO) LogRemoveRole(u *User, r int) error {
	role := Role{}
	d.db.Where(&Role{ID: r}).First(&role)
	logDetails := fmt.Sprintf("role %s removed", role.Name)
	logRecord := Log{UserID: 0, TargetID: u.ID, Type: LogUserChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *UsersDAO) LogAddRule(u *User, ruleIdx int) error {
	ruleName := rules.GetRuleName(ruleIdx)
	logDetails := fmt.Sprintf("rule %s added", ruleName)
	logRecord := Log{UserID: 0, TargetID: u.ID, Type: LogUserChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *UsersDAO) LogRemoveRule(u *User, ruleIdx int) error {
	ruleName := rules.GetRuleName(ruleIdx)
	logDetails := fmt.Sprintf("rule %s removed", ruleName)
	logRecord := Log{UserID: 0, TargetID: u.ID, Type: LogUserChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *RolesDAO) LogRoleAdded(role *Role) error {
	logRecord := Log{UserID: 0, TargetID: role.ID, Type: LogRoleAdded, Details: role.Name}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *RolesDAO) LogRoleDeleted(id int) error {
	role := Role{}
	d.db.Where(&Role{ID: id}).First(&role)
	logRecord := Log{UserID: 0, TargetID: role.ID, Type: LogRoleDeleted, Details: role.Name}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *RolesDAO) LogRoleChanged(role *Role, newName string) error {
	logDetails := fmt.Sprintf("role %s changed to %s", role.Name, newName)
	logRecord := Log{UserID: 0, TargetID: role.ID, Type: LogRoleChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *RolesDAO) LogRuleAdded(role *Role, ruleId int) error {
	ruleName := rules.GetRuleName(ruleId)
	logDetails := fmt.Sprintf("rule %s added", ruleName)
	logRecord := Log{UserID: 0, TargetID: role.ID, Type: LogRoleChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func (d *RolesDAO) LogRuleRemoved(role *Role, ruleId int) error {
	ruleName := rules.GetRuleName(ruleId)
	logDetails := fmt.Sprintf("rule %s removed", ruleName)
	logRecord := Log{UserID: 0, TargetID: role.ID, Type: LogRoleChanged, Details: logDetails}
	err := d.dao.Logs.Add(&logRecord)
	logError(err)
	return err
}

func GetDescription() map[LogType]LogTypeDescription {
	meta := make(map[LogType]LogTypeDescription)

	meta[LogLogin] = LogTypeDescription{"User login", "user"}
	meta[LogRoleAdded] = LogTypeDescription{"Role added", "role"}
	meta[LogRoleChanged] = LogTypeDescription{"Role data changed", "role"}
	meta[LogRoleDeleted] = LogTypeDescription{"Role deleted", "role"}
	meta[LogUserAdded] = LogTypeDescription{"User added", "user"}
	meta[LogUserChanged] = LogTypeDescription{"User data changed", "user"}
	meta[LogUserDeleted] = LogTypeDescription{"User deleted", "user"}

	return meta
}
