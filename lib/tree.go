package lib

import "github.com/thejerf/suture"

// CreateSupervisor returns a Supervisor struct
func CreateSupervisor() *suture.Supervisor {
	supervisor := suture.NewSimple("Supervisor")
	supervisor.Add(createIodine(supervisor))
	// supervisor.Add(createSocksProxy())
	// supervisor.Add(createNetworkConfig())
	return supervisor
}
