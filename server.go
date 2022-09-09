package main

import (
	"file_utils"
	"http_utils"
	"log"
	"net/http"
	"system_utils"
)

func main() {
	// check for cmd line args
	system_utils.Handle_cmd_line_args()
	//
	if system_utils.Mode == "dev" {
		system_utils.VERBOSE = true
	} else {
		system_utils.VERBOSE = false
		system_utils.Host_address_and_port = http_utils.Return_host_ip_address_and_port()
	}
	//
	if system_utils.Api_key != "" {
		println("[INFO] bind server to: " + system_utils.Host_address_and_port)
		println("[INFO] api key: " + system_utils.Api_key)
		println("[INFO] test server on address: http://" + system_utils.Host_address_and_port + "/?test=1&api_key=" + system_utils.Api_key)
		file_utils.Log_to_file("Info", "starting server on address: "+system_utils.Host_address_and_port)
	} else {
		println("[INFO] starting server on address: http://" + system_utils.Host_address_and_port + "/?test=1")
		file_utils.Log_to_file("Info", "starting server on address: "+system_utils.Host_address_and_port)
	}
	//
	go http_utils.Hub_global.Run()
	http.HandleFunc("/", http_utils.RootHandler)
	log.Fatal(http.ListenAndServe(system_utils.Host_address_and_port, nil))
}
