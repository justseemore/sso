package models

import (
	"encoding/json"
)

type Application struct {
	Base
	Name             string          `gorm:"size:100;not null;unique" json:"name"`
	Description      string          `gorm:"size:255" json:"description"`
	ClientID         string          `gorm:"size:100;not null;unique" json:"client_id"`
	ClientSecret     string          `gorm:"size:100;not null" json:"-"`
	RedirectURIs     string          `gorm:"type:text" json:"-"`
	AllowedScopes    string          `gorm:"type:text" json:"-"`
	Active           bool            `gorm:"default:true" json:"active"`
	ThemeID          *uint           `json:"theme_id"`
	Theme            *Theme          `gorm:"foreignKey:ThemeID" json:"theme,omitempty"`
	Settings         json.RawMessage `gorm:"type:json" json:"settings"`
}

// GetRedirectURIs 获取重定向URI列表
func (a *Application) GetRedirectURIs() ([]string, error) {
	var uris []string
	err := json.Unmarshal([]byte(a.RedirectURIs), &uris)
	if err != nil {
		return []string{}, err
	}
	return uris, nil
}

// SetRedirectURIs 设置重定向URI列表
func (a *Application) SetRedirectURIs(uris []string) error {
	jsonData, err := json.Marshal(uris)
	if err != nil {
		return err
	}
	a.RedirectURIs = string(jsonData)
	return nil
}

// GetAllowedScopes 获取允许的作用域列表
func (a *Application) GetAllowedScopes() ([]string, error) {
	var scopes []string
	err := json.Unmarshal([]byte(a.AllowedScopes), &scopes)
	if err != nil {
		return []string{}, err
	}
	return scopes, nil
}

// SetAllowedScopes 设置允许的作用域列表
func (a *Application) SetAllowedScopes(scopes []string) error {
	jsonData, err := json.Marshal(scopes)
	if err != nil {
		return err
	}
	a.AllowedScopes = string(jsonData)
	return nil
}

// GetSettings 获取应用设置
func (a *Application) GetSettings() (map[string]interface{}, error) {
	if a.Settings == nil || len(a.Settings) == 0 {
		return make(map[string]interface{}), nil
	}

	var settings map[string]interface{}
	err := json.Unmarshal(a.Settings, &settings)
	return settings, err
}

// SetSettings 设置应用设置
func (a *Application) SetSettings(settings map[string]interface{}) error {
	jsonData, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	a.Settings = jsonData
	return nil
}