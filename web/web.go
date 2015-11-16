package web

import (
	"fmt"

	"gopkg.in/macaron.v1"

	"github.com/containerops/dockyard/backend"
	"github.com/containerops/dockyard/middleware"
	"github.com/containerops/dockyard/router"
	"github.com/containerops/wrench/db"
	"github.com/containerops/wrench/setting"
)

func SetDockyardMacaron(m *macaron.Macaron) {
	//Setting Database
	if err := db.InitDB(setting.DBURI, setting.DBPasswd, setting.DBDB); err != nil {
		fmt.Printf("Connect Database error %s", err.Error())
	}

	if err := backend.InitBackend(); err != nil {
		fmt.Printf("Init backend error %s", err.Error())
	}

	if err := middleware.Initfunc(); err != nil {
		fmt.Printf("Init middleware error %s", err.Error())
	}

	//Setting Middleware
	middleware.SetMiddlewares(m)

	//Setting Router
	router.SetRouters(m)
}
