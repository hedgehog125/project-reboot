package twofactoractions

import (
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func ConfigureEndpoints(group *gin.RouterGroup, app *servercommon.ServerApp) {
	group.POST("/:id/confirm", Confirm(app))
}
