// (C) Ryan Warren 2015
// Authserver
//
// 

package main

import (
    "net/http"
    "os"
    //"os/exec"
    "net"
    "fmt"
    "strings"
)

func init() {
}

func getPath(w http.ResponseWriter, r *http.Request) {
}

func setPath(w http.ResponseWriter, r *http.Request) {
}

func errorer(w http.ResponseWriter, r *http.Request) {
}

func main() {
        //ifaces, ipError := net.Interfaces()
        //fmt.Println(net.Interfaces().InterfaceByName("eth0"))
        //fmt.Println(ifaces.InterfaceByName("eth0"))
        ief, err0 := net.InterfaceByName("eth0")
        if err0 !=nil{
                //log.Fatal(err)
        }
        addrs, err1 := ief.Addrs()
        if err1 !=nil{
                //log.Fatal(err)
        }
        //fmt.Println("HERE:")
        //fmt.Println(addrs)
        //fmt.Println(addrs[0])
        ipAddr := ""
        if addrs != nil {
            theIP := fmt.Sprintf("%s", addrs[0])
            ipAddr = fmt.Sprintf("%s", strings.Split(theIP, "/")[0])
        } else {
            ipAddr = "localhost"
        }
        fmt.Println(ipAddr)

        //fmt.Println(strings.Split(theIP, "/"))
        //fmt.Println(strings.Split(addrs[0].ToString(), "/"))
        //if ipError != nil {
        //  fmt.Println(ipError)
        //}
        //for _, i := range ifaces {
        //    addrs, err := i.Addrs()
        //    if err != nil {
        //      fmt.Println(err)
        //    }
        //    //fmt.Println(addrs.get("eth0"))
        //    for _, addr := range addrs {
        //        //fmt.Println(_)
        //        //fmt.Println(addr)
        //        switch v := addr.(type) {
        //            case *net.IPAddr:
        //            // process IP address
        //            fmt.Println(v)
        //            //fmt.Println(*net.IPAddr)
        //        }
        //    }
        //}
        //addrs, err3 := net.InterfaceAddrs()

        // if err3 != nil {
        //         fmt.Println(err3)
        //         os.Exit(1)
        // }
        // fmt.Println(addrs[*net.IPNet])
        // //fmt.Println(addrs.(*net.IPNet))

        // for _, address := range addrs {

        //       // check the address type and if it is not a loopback the display it
        //       if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
        //          if ipnet.IP.To4() != nil {
        //             fmt.Println(ipnet.IP.String())
        //          }

        //       }
        // }
        //test := "ifconfig"
        ////test := "ifconfig | grep -A 2 \"eth0\" | grep \"inet addr\" | cut -d: -f2 | awk '{ printf $$1}"
        //out, err4 := exec.Command(test).Output()
        //if err4 != nil {
        //              //fmt.Println("error occured")
        //              fmt.Printf("%s", err4)
        //}
        //fmt.Printf("%s", out)
	http.HandleFunc("/get", getPath)
	http.HandleFunc("/set", setPath)
	http.HandleFunc("/", errorer)
	err := http.ListenAndServe(":9090", nil)
	//err := http.ListenAndServe(portString, nil)
	if err != nil {
		//log.Errorf("Server Failed: %s", err)
		os.Exit(1)
	}
}


