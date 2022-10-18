package authio

import "github.com/sonyamoonglade/authio/internal/cookie"

type SettingBuilder struct {
	settings map[string]*cookie.Setting
}

func NewSettingBuilder() *SettingBuilder {
	return &SettingBuilder{
		settings: make(map[string]*cookie.Setting),
	}
}

func (b *SettingBuilder) Append(label string, setting *cookie.Setting) *SettingBuilder {
	b.settings[label] = setting
	return b
}
