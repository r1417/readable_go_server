package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    a "./webapi/tcpclient"
    u   "./utils"

  )

func main() {
    logfile := u.InitLogMultiWriter("./_webapiserver.log")
    defer logfile.Close()

    router := httprouter.New()
    router.GET("/TcpConnect", a.TcpConnect)
    router.GET("/Command", a.Command)
    router.GET("/ReadLogFile", a.ReadLogFile)


    err := http.ListenAndServe(":8070", router)
    u.OnErrorTerminate(err, "ListenAndServe")
}

