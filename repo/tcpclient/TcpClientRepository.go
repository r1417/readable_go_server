package tcpclient

import (
            "sync"
            "sync/atomic"
            "strconv"
            "errors"
        u   "../../utils"
        t   "../../tcp/client"
        )


type TcpClientForApi struct{
    ClientID string
    *t.TcpClient
}

type clients map[uint64]*TcpClientForApi


var(
    currentClientID uint64 = 0
    mu *sync.RWMutex = new(sync.RWMutex)
    cts = make(map[uint64]*TcpClientForApi)
)


func Regist(dest string, port int)(c *TcpClientForApi, err error) {

    client := t.CreateTcpClient(dest, port)

    err = client.Connect()

    if err != nil {
        return nil, err
    }
    
    //generate client ID
    icid := atomic.AddUint64(&currentClientID, 1)    
    cid := strconv.FormatUint(icid, 10)
    
    //regist to map
    c = &TcpClientForApi{cid, client}
    
    mu.Lock()
    defer mu.Unlock()
    cts[icid] = c
    
    return c, err
}


func Get(cid string) (c *TcpClientForApi, err error){

    errNotFound := errors.New("no such connection existed [" + cid + "]")

    if(cid == ""){
        return nil, errNotFound
    }

    var icid uint64
    icid, err = strconv.ParseUint(cid, 10, 64)

    if u.OnErrorResume(err, "repo.Get"){
        return c, err
    }

    mu.RLock()
    defer mu.RUnlock()

    if c, ok := cts[icid]; ok {
        return c, nil
    } else {
        return nil, errNotFound
    }
}


func Delete(cid string){

    icid, err := strconv.ParseUint(cid, 10, 64)

    if u.OnErrorResume(err, "repo.Delete"){
        return
    }

    //delete
    mu.Lock()
    defer mu.Unlock()

    if c, ok := cts[icid]; ok {
        c.Conn.Close()
        delete(cts, icid)

    }
}


/*


func ListAll() []TcpClientForApi{
    
}

*/