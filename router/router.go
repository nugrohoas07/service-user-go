package router

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB) {
	fmt.Println("Masuk ke init router")
}
