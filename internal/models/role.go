package models

type Role struct {
	Base
	Name        string      `gorm:"size:50;not null;unique" json:"name"`
	Description string      `gorm:"size:255" json:"description"`
	UserRoles   []UserRole  `gorm:"foreignKey:RoleID" json:"user_roles,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	Base
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Resource    string `gorm:"size:50;not null" json:"resource"`
	Action      string `gorm:"size:50;not null" json:"action"`
	Roles       []Role `gorm:"many2many:role_permissions;" json:"-"`
}

type UserRole struct {
	Base
	UserID   uint   `gorm:"not null" json:"user_id"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	User     User   `gorm:"foreignKey:UserID" json:"-"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}