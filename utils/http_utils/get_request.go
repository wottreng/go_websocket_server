package http_utils

import (
	"file_utils"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"system_utils"
)

// main function to handle GET requests
func Get_request_handler(w http.ResponseWriter, r *http.Request) {
	//
	args := r.URL.Query()
	if system_utils.VERBOSE == true {
		println("args received: <" + args.Encode() + ">")
	}
	// serve up html
	if args.Get("test") == "1" {
		echo_html := template.Must(template.ParseFiles("./www/templates/echo.html"))
		err := echo_html.Execute(w, "ws://"+r.Host)
		if err != nil {
			return
		}
		return
	}
	// send message
	if args.Get("topic") != "" {
		// check that topic is not in Hub_global.topics, add it if not
		topic_found := false
		for _, topic := range Hub_global.topics {
			if topic == args.Get("topic") {
				topic_found = true
				break
			}
		}
		if topic_found == false {
			Hub_global.topics = append(Hub_global.topics, args.Get("topic"))
		}
		serve_ws_request(Hub_global, w, r, args.Get("topic"))
		return
	}
	// return all topics
	if args.Get("list_topics") == "1" {
		return_all_topics(w)
		return
	}
	file_utils.Log_to_file("debug", "not an endpoint: "+r.URL.Path)
	fmt.Fprintf(w, "{not an endpoint}")

	return
}

func return_all_topics(w http.ResponseWriter) {
	topic_titles := make([]string, 0)
	for _, topic := range Hub_global.topics {
		topic_titles = append(topic_titles, topic)
	}
	json_data := "{\"topics\":[" + strings.Join(topic_titles, ",") + "]}"
	_, err := fmt.Fprintf(w, "%s\n", json_data)
	if err != nil {
		file_utils.Log_error_to_file(err, "Get_request_handler")
	}
	return
}
