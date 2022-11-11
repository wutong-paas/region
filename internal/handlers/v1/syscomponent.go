/*
Copyright 2022 The Wutong Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wutong-paas/region/internal/services"
)

func ListSysComponents(c *gin.Context) {
	sysComponents, err := services.DefaultServicer().ListSysComponents()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, sysComponents)
}

func InstallSysComponent(c *gin.Context) {
	var data struct {
		Name    string `json:"name" binding:"required"`
		Version string `json:"version"`
	}

	c.Bind(&data)
	err := services.DefaultServicer().InstallComponnet(data.Name, data.Version)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}
