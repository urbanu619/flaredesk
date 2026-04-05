package service

import (
	"go_server/service/app"
	"go_server/service/login"
	"go_server/service/system"
)

// Define all services

var RealizationLayer = new(Group)

type Group struct {
	LoginServiceGroup  login.ServiceGroup
	SystemServiceGroup system.ServiceGroup
	AppServiceGroup    app.ServiceGroup
}
