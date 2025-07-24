package main

import (
	"hope_blog/router"
)

func main() {
	route := router.Router()
	route.Run(":9999")
}
