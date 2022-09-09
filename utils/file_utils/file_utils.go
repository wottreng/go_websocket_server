package file_utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time_utils"
)

//function to check if file exists
func Does_file_exist(file_path string) bool {
	info, err := os.Stat(file_path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//function to check if filder exists
func Does_folder_exist(folder_path string) bool {
	info, err := os.Stat(folder_path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

//function to list all files in directory
func List_all_files_in_directory(path string) []string {
	var files []string
	var err error
	//
	files_list, err := ioutil.ReadDir(path)
	if err != nil {
		Log_error_to_file(err, "List_all_files_in_directory")
		return files
	}
	//
	for _, file := range files_list {
		files = append(files, file.Name())
	}
	return files
}

//function to create folder
func CreateFolder(folderPath string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, 0755)
		if err != nil {
			Log_error_to_file(err, "CreateFolder")
		}
	}
}

//function to create file
func CreateFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			Log_error_to_file(err, "CreateFile")
		}
		defer file.Close()
	}
}

func Read_string_from_file(path string, file_name string) string {
	var string_data string
	var err error
	var byte_data []byte
	var message string
	//
	absolute_path := path + "/" + file_name
	number_of_attempts := 5
	for i := 0; i < number_of_attempts; i++ {
		byte_data, err = os.ReadFile(absolute_path)
		if err != nil {
			Log_error_to_file(err, "Read_string_from_file")
			message = "data read error\n"
			return message
		}
		if len(byte_data) > 0 {
			break
		}
	}
	if len(byte_data) < 1 {
		return "data read error\n"
	}
	//
	string_data = string(byte_data)
	return string_data
}

func Write_string_to_file(data_string string, path string, file_name string) bool {
	var err error
	absolute_path := path + "/" + file_name
	//
	if !Does_folder_exist(path) {
		CreateFolder(path)
	}
	//
	f, err := os.OpenFile(absolute_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		Log_error_to_file(err, "Write_string_to_file")
		return false
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			Log_error_to_file(err, "Write_string_to_file")
		}
	}(f)
	//
	logger := log.New(f, "", 0) // used to handle concurrent writes
	err = logger.Output(2, data_string)
	if err != nil {
		Log_error_to_file(err, "Write_string_to_file")
		return false
	}
	//
	return true
}

func Write_data_to_file(data []byte, path string, file_name string) bool {
	var err error
	absolute_path := path + "/" + file_name
	file, err := os.Create(absolute_path)
	if err != nil {
		Log_error_to_file(err, "Write_data_to_file")
	}
	//
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Log_error_to_file(err, "Write_data_to_file")
		}
	}(file)
	//
	_, err = file.Write(data)
	if err != nil {
		Log_error_to_file(err, "Write_data_to_file")
	}
	//
	return true
}

// function for writing an error log
func Log_error_to_file(err error, custom_message ...string) {
	var error_string string
	if len(custom_message) > 0 {
		error_string = fmt.Sprintf("[-->] Error: %v - %v", custom_message[0], err)
	} else {
		error_string = fmt.Sprintf("[-->] Error: %v", err)
	}
	cwd, _ := os.Getwd()
	error_log_path := cwd + "/logs"
	error_log_file := "error_log.txt"
	//
	f, err := os.OpenFile(error_log_path+"/"+error_log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0655)
	logger := log.New(f, "", log.LstdFlags)
	logger.Output(2, error_string)
	//Write_string_to_file(error_string, error_log_path, error_log_file)
}

// function for writing an error log
func Log_to_file(log_type string, custom_message string) {
	var data_string = fmt.Sprintf("[%v]: %v", log_type, custom_message)
	cwd, _ := os.Getwd()
	error_log_path := cwd + "/logs"
	error_log_file := "run_log.txt"
	//
	f, err := os.OpenFile(error_log_path+"/"+error_log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		Log_error_to_file(err, "[file_utils.go][Log_to_file]")
	}
	logger := log.New(f, "", log.LstdFlags)
	logger.Output(2, data_string)
}

// function to build file name with topic and current date
func Build_file_name(topic string) string {
	filename := topic + "_" + time_utils.Return_current_date() + ".txt"
	return filename
}
