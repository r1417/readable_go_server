package main

import (
  t   "./tcp/server"
  u   "./utils"

  )

func main() {
    logfile := u.InitLogMultiWriter("./_tcpserver.log")
    defer logfile.Close()
        
    t.ServerListen(8080)
}



