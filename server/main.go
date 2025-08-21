package main

import (
    "fmt"
    "log"

    "nadhi.dev/sarvar/fun/server" 
    "nadhi.dev/sarvar/fun/routes"
)

func webserver(port int) {
    log.Fatal(server.Route.Listen(fmt.Sprintf(":%d", port)))
}

func main() {
    routes.Register() // Register all routes (they use server.Route)
    webserver(317)
}