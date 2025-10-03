package main

import (
	"fmt"
	"log"

	webview "github.com/webview/webview_go"
	config "nadhi.dev/sarvar/fun/config"
	logg "nadhi.dev/sarvar/fun/logs"
	"nadhi.dev/sarvar/fun/routes"
	"nadhi.dev/sarvar/fun/server"
	sheet "nadhi.dev/sarvar/fun/sheets"
)



func init() {
    var err error
    var queue_dir string
    
    
    queue_dir_val := config.GetConfigValue("SHEET_QUEUE_DIR")
    if queue_dir_val != nil {
        queue_dir = queue_dir_val.(string)
    } else {
        queue_dir = "./storage/queue_data" // Fallback to default
    }
    
    sheet.GlobalSheetGenerator, err = sheet.NewSheetGenerator(nil, queue_dir, 2)
    if err != nil {
        logg.Error(fmt.Sprintf("Failed to initialize GlobalSheetGenerator: %v", err))
        logg.Exit()
    }
    logg.Success("GlobalSheetGenerator initialized successfully")
}



func webserver(port int) {
	log.Fatal(server.Route.Listen(fmt.Sprintf(":%d", port)))
}

func main() {
    routes.SetAssetsPath(ExtractedAssetsPath)
    routes.Register()
    go webserver(317)

    w := webview.New(true)
    defer w.Destroy()
    w.SetTitle("Vela by Nadhi.dev | Beta-synopsis")
    w.SetSize(1024, 768, webview.HintNone)
    w.Navigate("http://127.0.0.1:317")

    w.Run()
}