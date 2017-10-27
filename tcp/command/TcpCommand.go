package command

import (
    "net"
    "time"
    "strings"
)


type Command struct {
  Jsonrpc   string    `json:"jsonrpc"`
  Method    string    `json:"method"`
  Params  []string    `json:"params"`
  ID        string    `json:"id"`
}

func CreateCommand(m string, p []string, id string) *Command {
    c := new(Command)
    c.Jsonrpc = "2.0"
    c.Method = m
    c.Params = p
    c.ID = id
    return c
}

type Response struct {
  Jsonrpc   string               `json:"jsonrpc"`
  Result    map[string]string    `json:"result"`
  ID        string               `json:"id"`
}

func CreateResponse(result map[string]string, id string) *Response {
    r := new(Response)
    r.Jsonrpc = "2.0"
    r.Result = result
    r.ID = id
    return r
}

type WillClosed bool

type CommandExecutor interface{
    Execute(*Command, net.Conn) (*Response, WillClosed)
}



//hello
type CmdHello struct {
}

func (exe *CmdHello) Execute(cmd *Command, conn net.Conn) (*Response, WillClosed) {

    r := make(map[string]string)
    r["command"] = cmd.Method
    r["message"] = "Ciao  " + conn.RemoteAddr().String() + " -> " + conn.LocalAddr().String()
        
    return CreateResponse(r, cmd.ID), false

}


//bye
type CmdBye struct {
}

func (exe *CmdBye) Execute(cmd *Command, conn net.Conn) (*Response, WillClosed) {

    r := make(map[string]string)
    r["command"] = cmd.Method
    r["message"] = "Ciao"
        
    return CreateResponse(r, cmd.ID), true

}

//now
type CmdNow struct {
}

func (exe *CmdNow) Execute(cmd *Command, conn net.Conn) (*Response, WillClosed) {

    r := make(map[string]string)
    r["command"] = cmd.Method
    r["message"] = time.Now().String()
        
    return CreateResponse(r, cmd.ID), false

}

//test2
type CmdTest2 struct {
}

func (exe *CmdTest2) Execute(cmd *Command, conn net.Conn) (*Response, WillClosed) {

    r := make(map[string]string)
    r["command"] = cmd.Method
    r["message"] = "test2です"
        
    return CreateResponse(r, cmd.ID), false

}

//error
type CmdError struct {
}

func (exe *CmdError) Execute(cmd *Command, conn net.Conn) (*Response, WillClosed) {

    r := make(map[string]string)
    r["command"] = cmd.Method
    r["message"] = "command not found"
        
    return CreateResponse(r, cmd.ID), false

}


var CommandExecs map[string]CommandExecutor = map[string]CommandExecutor{
    "hello":new(CmdHello),
    "bye":new(CmdBye),
    "now?":new(CmdNow),
    "test2":new(CmdTest2),
    "error":new(CmdError),
}

func GetCommandExecutor(c *Command) CommandExecutor {
    exec := CommandExecs[strings.ToLower(c.Method)]
    if exec == nil {
        return CommandExecs["error"]
    } else {
        return exec
    }
}
