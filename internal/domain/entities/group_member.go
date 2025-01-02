package entities

type GroupRole string

const (
	Admin  GroupRole = "admin"
	Member GroupRole = "member"
)

type GroupMember struct {
	GroupId string    `json:"group_id" gorm:"primary_key"`
	UserId  string    `json:"user_id" gorm:"primary_key"`
	Role    GroupRole `json:"role" gorm:"type:varchar(10);check:role in ('admin','member')"`

	Group Group `json:"group" gorm:"foreignKey:GroupId"`
	User  User  `json:"user" gorm:"foreignKey:UserId"`
}
