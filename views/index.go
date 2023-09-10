package views

import (
	"net/http"
	"strings"

	"github.com/mbdeguzman/hdd_monitoring/models"
	"github.com/uadmin/uadmin"
)

type StorageResult struct {
	Server                models.Servers
	ServerID              uint
	TotalStorage          string
	AvailableStorage      string
	UsedStorage           string
	UsedStoragePercentage string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	context := map[string]interface{}{}
	context["Title"] = uadmin.SiteName
	storage := []StorageResult{}
	storages := []models.Storage{}
	uadmin.All(&storages)
	for i := range storages {
		uadmin.Preload(&storages[i])
		storage = append(storage, StorageResult{
			Server:                storages[i].Server,
			ServerID:              storages[i].ServerID,
			TotalStorage:          storages[i].TotalStorage,
			AvailableStorage:      storages[i].AvailableStorage,
			UsedStorage:           storages[i].UsedStorage,
			UsedStoragePercentage: strings.Trim(storages[i].UsedStoragePercentage, "%"),
		})
	}

	context["Storage"] = storage
	path := "templates/storage/index.html"
	uadmin.RenderHTML(w, r, path, context)
}
