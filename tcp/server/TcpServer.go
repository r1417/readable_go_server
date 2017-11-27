package server

import (
    "net"
    "time"
    "strconv"
    "encoding/json"
    "io"
    "sync/atomic"
    u "../../utils"
    tcpc "../command"

)


var sessionID uint64 = 0


func ServerListen(port int) {

    tcpAddr, err := net.ResolveTCPAddr("tcp4", ":" + strconv.Itoa(port))
    u.OnErrorTerminate(err, "ResolveTCPAddr")

    //listen
    listener, err := net.ListenTCP("tcp", tcpAddr)
    u.OnErrorTerminate(err, "ListenTCP")

    u.InfoJson(*tcpAddr, "Listen")

	defer listener.Close()

    //accept -> fork
    for {
        conn, err := listener.Accept()

        if u.OnErrorResume(err, "Accept"){
            continue
        }
        
        sid := atomic.AddUint64(&sessionID, 1)

        u.Info("fork" ,"Accept " + strconv.FormatUint(sid, 16))

        go fork(strconv.FormatUint(sid, 16), conn)
    }

}


func fork(sid string, conn net.Conn) {

    defer conn.Close()

    conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) 


    for {
        request := make([]byte, 1024)

        u.Info("wait" ,"read " + sid)

        //receive
        readLength, err := conn.Read(request)    //conn.netFD.Raad -> poll.FD.Read ->  syscall.Read

        neterr, ok := err.(net.Error)
        if ok && neterr.Timeout() {
            u.Info("Timeout" ,"disconnect " + sid)
            break

        } else if err == io.EOF {
            u.Info("EOF" ,"disconnect " + sid)
            break

        } else if  u.OnErrorResume(err, "disconnect " + sid){
            break

        }else if readLength == 0 {
            u.Info("no data" ,"request " + sid)
            continue
        }

        u.Info(string(request) ,"request " + sid)

        //execute command
        cmd := new(tcpc.Command)
        err = json.Unmarshal(request[:readLength], cmd)

        if u.OnErrorResume(err, "request Unmarshal " + sid){
            break
        }
        
        cmd.ID = sid

        var c tcpc.CommandExecutor = tcpc.GetCommandExecutor(cmd)
        res, willClosed := c.Execute(cmd, conn)

        u.InfoJson(res ,"commandResult " + sid)

        //write response 
        b, _ := json.Marshal(res)
        conn.Write(b)
            
        if(willClosed){
            break
        }
    }
}

