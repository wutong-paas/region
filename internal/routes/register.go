package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/wutong-paas/region/internal/handlers/v1"
)

// Register registers the routes in the router.
func (r *Router) Register() *Router {
	r.v1()
	return r
}

// Router is a wrapper for gin.RouterGroup.
type Router struct {
	//*gin.Engine
	*gin.RouterGroup
}

// SubRouter register a sub-router group with a prefix and returns it.
func (r *Router) SubRouter(relativePath string, handlers ...gin.HandlerFunc) *Router {
	grouped := r.Group(relativePath, handlers...)
	//return &Router{r.Engine, grouped}
	return &Router{grouped}
}

// v1 registers the v1 routes in the router.
func (r *Router) v1() *Router {
	v1 := r.SubRouter("/v1")
	// add v1 sub routers here.
	v1.sysComponents()

	return v1
}

// sys-components registers the v1/sys-components routes in the router.
func (r *Router) sysComponents() *Router {
	sysComponents := r.SubRouter("/sys-components")
	sysComponents.GET("", handlers.ListSysComponents)
	sysComponents.POST("", handlers.InstallSysComponent)
	sysComponents.PUT("", handlers.UpgradeSysComponent)
	sysComponents.DELETE("", handlers.UninstallSysComponent)
	return sysComponents
}
