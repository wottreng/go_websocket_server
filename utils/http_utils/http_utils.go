package http_utils

import (
	"file_utils"
	"fmt"
	"log"
	"net"
	"net/http"
	"system_utils"
)

// function to handle requests to `/` endpoint
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if system_utils.VERBOSE == true {
		println("\nHeaders: %+v\n", r.Header)
	}
	//check if topic arg is present
	if Check_if_args_are_present(w, r) == false {
		return
	}
	//check if api arg is present
	if Check_if_api_key_arg_is_present(r) == false {
		// return api key error
		w.Write([]byte("{\"error\":\"api key incorrect\"}"))
		return
	}
	// method handlers
	switch r.Method {
	case "GET":
		Get_request_handler(w, r)
	default:
		Method_not_allowed_handler(w)
	}
}

// function to check if topic arg is present
func Check_if_args_are_present(w http.ResponseWriter, r *http.Request) bool {
	args := r.URL.Query()
	// if help arg is present, return help message
	if args.Has("help") {
		_, err := fmt.Fprintf(w, "Usage: %s?topic=<topic>\n", r.URL.Path)
		if err != nil {
			file_utils.Log_error_to_file(err, "[http_utils.go][Check_if_args_are_present]")
			return false
		}
		return false
	}
	// if topic arg is not present, return error message
	if !args.Has("topic") && !args.Has("list_topics") && !args.Has("test") {
		_, err := fmt.Fprintf(w, "not an endpoint\n")
		if err != nil {
			file_utils.Log_error_to_file(err, "[http_utils.go][Check_if_args_are_present]")
			return false
		}
		file_utils.Log_to_file("debug", "topic arg not present")
		return false
	}
	return true
}

// function to check if api arg is present in request
func Check_if_api_key_arg_is_present(r *http.Request) bool {
	// if api key is not set then return true
	if system_utils.Api_key == "" {
		return true
	}
	// check if api key from request is correct
	args := r.URL.Query()
	// if api key not present return false
	if !args.Has("api_key") {
		return false
	}
	if args.Has("api_key") {
		api_key_sent := args.Get("api_key")
		if api_key_sent == system_utils.Api_key {
			return true
		}
		file_utils.Log_to_file("debug", "api key sent: "+api_key_sent)
		return false
	}
	file_utils.Log_to_file("debug", "api key not present")
	return false
}

// return response for methods not supported
func Method_not_allowed_handler(w http.ResponseWriter) {
	_, err := fmt.Fprintf(w, "Method not allowed\n")
	if err != nil {
		file_utils.Log_error_to_file(err, "[http_utils.go][Method_not_allowed_handler]")
		return
	}
}

// function to return host ip address and port
func Return_host_ip_address_and_port() string {
	only_bind_to_acceptable_addresses := false // make true if control over bind address is needed
	acceptable_addresses := []string{"192.168.", "10.42."}
	host_ip_address := ""
	host_port := "8080"
	//
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		file_utils.Log_error_to_file(err, "[http_utils.go][Return_host_ip_address_and_port]")
		log.Fatal("[Error] " + err.Error() + "\n")
	}
	//
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if only_bind_to_acceptable_addresses {
					println(ipnet.IP.String() + "\n")
					for _, acceptable_address := range acceptable_addresses {
						if ipnet.IP.String()[0:len(acceptable_address)] == acceptable_address {
							host_ip_address = ipnet.IP.String()
							break
						}
					}
				} else {
					host_ip_address = ipnet.IP.String()
					break
				}
			}
		}
	}
	if host_ip_address == "" {
		file_utils.Log_error_to_file(err, "[http_utils.go][Return_host_ip_address_and_port]")
		log.Fatal("[ERROR] Could not find host ip address")
	}
	//
	public_address_and_port := host_ip_address + ":" + host_port
	return public_address_and_port
}
