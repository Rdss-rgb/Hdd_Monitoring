package models

import "github.com/uadmin/uadmin"

type Servers struct {
	uadmin.Model
	ServerName string
	IPAddress  string
	Username   string `uadmin:"list_exclude;"`
	Password   string `uadmin:"list_exclude;"`
	Port       int    `uadmin:"list_exclude"`
}

func (s *Servers) String() string {
	return s.ServerName
}
