package tcpclient

import (
    "net/http"
    "encoding/json"
    "io/ioutil"
    s "strings"
    "github.com/julienschmidt/httprouter"
    repo "../../repo/tcpclient"
    tcpc "../../tcp/command"
    u   "../../utils"

  )

func TcpConnect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    r.ParseForm()
    u.InfoJson(r.Form, "TcpConnect:params")

    //createClient -> connect
    client, err := repo.Regist("localhost", 8080)

    if u.OnErrorResume(err, "TcpConnect"){
        sendError(w, err)
        return
    }

    //hello
    cmd := tcpc.CreateCommand("hello", []string{}, "")

    execCommandAndRespond(w, client, cmd)

}


func Command(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    r.ParseForm()
    u.InfoJson(r.Form, "Command:params")

    //find client
    client, err := repo.Get(r.FormValue("id"))

    if u.OnErrorResume(err, "Get"){
        sendError(w, err)
        return
    }

    //create command
    commandText := s.TrimSpace(r.FormValue("method")) 
    args := s.Split(commandText, " ")
    
    u.InfoJson(args, "Split")

    cmd := tcpc.CreateCommand(args[0], args[1:], "")

    //exec
    execCommandAndRespond(w, client, cmd)

}



func ReadLogFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    r.ParseForm()

    fileType := r.FormValue("method")

    filePath := ""
        
    switch fileType {
    case "WebAPIServer" :
        filePath = "_webapiserver.log"
    case "TcpServer" :
        filePath = "_tcpserver.log"
    }

    log, err := ioutil.ReadFile(filePath)

    if u.OnErrorResume(err, "ReadLogFile"){
        sendError(w, err)
        return
    }

    //response
    result := map[string]string{"fileType":fileType, "log" : string(log)}
    res := tcpc.CreateResponse(result, "")

    sendResponse(w, res)

}



func sendResponse(w http.ResponseWriter, res interface{}) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)

    err := json.NewEncoder(w).Encode(res)
    u.OnErrorResume(err, "SendResponse")

}

func sendError(w http.ResponseWriter, err error) {
    result := map[string]string{"command":"error", "message" : err.Error(), "cid":"" }
    res := tcpc.CreateResponse(result, "")

    u.InfoJson(res ,"Response")

    sendResponse(w, res)

}

func execCommandAndRespond(w http.ResponseWriter, client *repo.TcpClientForApi, cmd *tcpc.Command) {

    hasClosed, res, err := client.Send(cmd)

    if hasClosed {
        repo.Delete(client.ClientID)
    }

    if u.OnErrorResume(err, "Send"){
        sendError(w, err)

    } else {
        res.Result["cid"] = client.ClientID
    
        u.InfoJson(res ,"Response")

        sendResponse(w, res)
    }

}

/*

func TcpExecCommand(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    

}


func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
  fmt.Fprintf(w, "%q", html.EscapeString(ps.ByName("id")))
}


func CreateTcpClient(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    c := Client{Name: "insertTest", Completed: false}

    c = createClient(c)
    log.Print(c)

    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(clients); err != nil {
        log.Print(err);
    }
    return
}


func JTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    

    if err := json.NewEncoder(w).Encode(clients); err != nil {
        log.Fatal(err)
    }

}

*/