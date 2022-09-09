module main

replace file_utils => ./utils/file_utils

replace time_utils => ./utils/time_utils

replace system_utils => ./utils/system_utils

replace http_utils => ./utils/http_utils

go 1.18

require (
	file_utils v0.0.0-00010101000000-000000000000
	http_utils v0.0.0-00010101000000-000000000000
	system_utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	time_utils v0.0.0-00010101000000-000000000000 // indirect
)
