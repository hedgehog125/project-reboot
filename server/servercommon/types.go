package servercommon

import (
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
)

type ServerApp struct {
	// TODO: struct embed App?
	App             *common.App
	Router          *gin.Engine
	AdminMiddleware gin.HandlerFunc
}
