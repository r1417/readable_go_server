package main

import (
          "os"
            "fmt"
            "bufio"
        s   "strings"
        t   "./tcp/client"
        tcpc "./tcp/command"
        u   "./utils"
    )

func main() {

    client := t.CreateTcpClient("localhost", 8080)

    //connect
    err := client.Connect()

    u.OnErrorTerminate(err, "Connect")

    defer client.Conn.Close()

    //stdin
    stdin := bufio.NewScanner(os.Stdin)
    for stdin.Scan() {

        //build command
        if err := stdin.Err(); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }

        commandText := stdin.Text()

        if commandText == "" {
            continue
        }

        args := s.Split(commandText, " ")
        
        u.InfoJson(args, "Split")

        cmd := tcpc.CreateCommand(args[0], args[1:], "")

        //send
        hasClosed, res, e := client.Send(cmd)


        u.OnErrorResume(e, "Send")

        if hasClosed {
            break
        }

        u.InfoJson(res ,"commandResult")


        //sandbox
        /*

        u.OnErrorTerminate(errors.New("hoge"), "Sandbox")

        use defer:  defer client.Conn.Close()
            → 2017/11/20 15:04:34 info: [disconnect 2] EOF
            use log.fatal
            → 2017/11/20 15:07:21 error: [disconnect 3] read tcp 127.0.0.1:8080->127.0.0.1:62818: wsarecv: An existing connection was forcibly closed by the remote host.
        no defer: 
            → 2017/11/20 15:07:21 error: [disconnect 3] read tcp 127.0.0.1:8080->127.0.0.1:62818: wsarecv: An existing connection was forcibly closed by the remote host.
        */
    }

}


