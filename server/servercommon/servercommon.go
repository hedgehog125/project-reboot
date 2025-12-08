package servercommon

import (
	"github.com/NicoClack/cryptic-stash/common"
	"github.com/gin-gonic/gin"
)

type ServerApp struct {
	*common.App
	Router          *gin.Engine
	AdminMiddleware gin.HandlerFunc
}

type Group struct {
	*gin.RouterGroup
	App *ServerApp
}

func (group *Group) Group(relativePath string) *Group {
	return &Group{
		RouterGroup: group.RouterGroup.Group(relativePath),
		App:         group.App,
	}
}
