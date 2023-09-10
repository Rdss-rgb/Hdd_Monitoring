package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/mbdeguzman/hdd_monitoring/models"
	"github.com/mbdeguzman/hdd_monitoring/views"
	"github.com/uadmin/uadmin"
	"golang.org/x/crypto/ssh"
)

func main() {
	uadmin.Register(
		models.Servers{},
		models.Storage{},
	)
	http.HandleFunc("/", uadmin.Handler(views.IndexHandler))
	Services()
	uadmin.RootURL = "admin/"
	uadmin.SiteName = "Storage Monitoring"
	uadmin.Port = 1122
	uadmin.StartServer()
}

func Services() {
	go func() {
		for {
			servers := []models.Servers{}
			uadmin.All(&servers)
			for i := range servers {
				uadmin.Preload(&servers[i])

				username := servers[i].Username
				password := servers[i].Password
				ipaddress := servers[i].IPAddress
				port := servers[i].Port

				// SSH configuration
				sshConfig := &ssh.ClientConfig{
					User: username,
					Auth: []ssh.AuthMethod{
						ssh.Password(password),
					},
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				}
				remoteHost := ipaddress + ":" + fmt.Sprint(port)

				client, err := ssh.Dial("tcp", remoteHost, sshConfig)
				if err != nil {
					uadmin.Trail(uadmin.ERROR, "Failed to connect: %v", err)
					continue
				}
				defer client.Close()

				session, err := client.NewSession()
				if err != nil {
					uadmin.Trail(uadmin.ERROR, "Failed to create session: %v", err)
					continue

				}
				defer session.Close()

				output, err := session.CombinedOutput("df -BG /")
				if err != nil {
					uadmin.Trail(uadmin.ERROR, "Failed to run command: %v", err)
					continue

				}

				storageOutput := string(output)
				pattern := `(\d+G)\s+(\d+G)\s+(\d+G)\s+(\d+%)\s+\/`
				re := regexp.MustCompile(pattern)
				matches := re.FindStringSubmatch(storageOutput)

				if len(matches) < 5 {
					uadmin.Trail(uadmin.ERROR, "Error: No matches found")
					continue

				}

				oneGBlocks := matches[1]
				used := matches[2]
				available := matches[3]
				usedPercent := matches[4]

				uadmin.Trail(uadmin.INFO, "SERVER:%v", servers[i].ServerName)
				uadmin.Trail(uadmin.INFO, "1G-blocks: %s\n", oneGBlocks)
				uadmin.Trail(uadmin.INFO, "Used: %s\n", used)
				uadmin.Trail(uadmin.INFO, "Available: %s\n", available)
				uadmin.Trail(uadmin.INFO, "Used %%: %s\n\n", usedPercent)

				storage := models.Storage{}

				if uadmin.Count(&storage, "server_id = ?", servers[i].ID) > 0 {
					uadmin.Get(&storage, "server_id = ?", servers[i].ID)
					storage.Server = servers[i]
					storage.TotalStorage = oneGBlocks
					storage.UsedStorage = used
					storage.AvailableStorage = available
					storage.UsedStoragePercentage = usedPercent
					uadmin.Save(&storage)
				} else {

					storage.Server = servers[i]
					storage.TotalStorage = oneGBlocks
					storage.UsedStorage = used
					storage.AvailableStorage = available
					storage.UsedStoragePercentage = usedPercent
					uadmin.Save(&storage)
				}
			}

			time.Sleep(time.Minute * 1)
		}
	}()
}
