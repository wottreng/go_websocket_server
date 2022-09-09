module http_utils

replace file_utils => ../file_utils

replace time_utils => ../time_utils

replace system_utils => ../system_utils

go 1.18

require (
	file_utils v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.0
	system_utils v0.0.0-00010101000000-000000000000
	time_utils v0.0.0-00010101000000-000000000000
)
