package client

import (
        "net"
        "strconv"
        "encoding/json"
        u   "../../utils"
        tcpc "../command"
    )


type TcpClient struct {
    Conn        net.Conn
    Dest        string
    Port        int
    SessionID   string
}


func CreateTcpClient(dest string, port int) *TcpClient {
    c := new(TcpClient)
    c.Dest = dest
    c.Port = port

    return c
}


func (c *TcpClient) Connect() (err error){

    c.Conn, err = net.Dial("tcp4", c.Dest + ":" + strconv.Itoa(c.Port))

   if u.OnErrorResume(err, "Connect"){
        c.Conn = nil
        return err
    }
    
    return nil
}


func (c *TcpClient)Send(cmd *tcpc.Command) (hasClosed bool, res *tcpc.Response, err error){

    res = tcpc.CreateResponse(map[string]string{"command":"none", "message": "invalid command"}, "")


    //client check
    if c.Conn == nil {
        return true, res, err
    }

    if cmd.Method == "" {
        return false, res, err
    }

    //send
    b, _ := json.Marshal(cmd)
    c.Conn.Write(b)


    //response
    response := make([]byte, 1024)
    readLength := 0
    readLength, err = c.Conn.Read(response)

    if  u.OnErrorResume(err, "disconnect") {
        c.Conn.Close()
        return true, res, err

    } else if readLength == 0 {
        u.Info("no data" ,"response")
        c.Conn.Close()
        return true, res, err
    }

    err = json.Unmarshal(response[:readLength], res)

    if u.OnErrorResume(err, "UnmarshalResponse"){
        c.Conn.Close()
        return true, res, err        
    }

    return false, res, err        
    
}
