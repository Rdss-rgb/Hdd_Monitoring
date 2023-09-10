package models

import "github.com/uadmin/uadmin"

type Storage struct {
	uadmin.Model
	Server                Servers `uadmin:"list_exclude"`
	ServerID              uint
	TotalStorage          string
	AvailableStorage      string
	UsedStorage           string
	UsedStoragePercentage string
}

func (s *Storage) String() string {
	uadmin.Preload(s)
	return s.Server.ServerName
}
