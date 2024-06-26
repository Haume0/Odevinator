package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Auth struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func AuthCheck(id string, code string, w http.ResponseWriter, r *http.Request) bool {
	if _, err := os.Stat("./auth.json"); os.IsNotExist(err) {
		//create file with {} in it
		file, err := os.Create("./auth.json")
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err.Error())
		}
		defer file.Close()
		file.WriteString("[]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	authData, err := os.ReadFile("./auth.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	var auth []Auth
	err = json.Unmarshal(authData, &auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	var authResult bool = false
	//check id and code is correct
	for _, v := range auth {
		if v.ID == id && v.Code == code {
			// verified
			authResult = true
			break
		}
	}
	return authResult
}
func Login(w http.ResponseWriter, r *http.Request) {
	//Getting query parameters from the request */login?id=2314716027
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	if id == "" || name == "" {
		fmt.Fprintf(w, "unknown")
		return
	}
	var code = GenerateVerifyCode()
	//reading file
	authData, err := os.ReadFile("./auth.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err.Error())
	}
	var auth []Auth
	//unmarshal the json data
	err = json.Unmarshal(authData, &auth)
	if err != nil {
		fmt.Printf("Error at unmarshall json: %v\n", err.Error())
	}
	//if the id is already in the file, return
	for _, v := range auth {
		if v.ID == id {
			fmt.Printf("\r")
			fmt.Printf("\n💯 %v • %v: 🔑 %v 🔑 💢Zaten kayıtlı!", id, name, v.Code)
			fmt.Fprintf(w, "exists")
			return
		}
	}
	auth = append(auth, Auth{ID: id, Code: code}) // Assign the result of append to 'auth'
	//marshal the data
	authData, err = json.Marshal(auth)
	if err != nil {
		fmt.Printf("Error at marshall json: %v\n", err.Error())
	}
	//write the data to the file
	err = os.WriteFile("./auth.json", authData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err.Error())
	}
	//send the code to the user
	SendVerifyMail(id+OKUL_SUFFIX, id, code, name)
	fmt.Fprintf(w, "done")
}
func Verify(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	code := r.URL.Query().Get("code")
	if !AuthCheck(id, code, w, r) {
		http.Error(w, "not verified", http.StatusUnauthorized)
		return
	} else {
		fmt.Fprint(w, "verified")
	}
}
func New(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(64 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	odevFiles := r.MultipartForm.File["files"]
	ogrID := r.FormValue("ogr_id")
	ogrCode := r.FormValue("ogr_code")
	homeworkLesson := r.FormValue("homework_lesson")
	homeworkName := r.FormValue("homework_name")

	if !AuthCheck(ogrID, ogrCode, w, r) {
		http.Error(w, "not verified", http.StatusUnauthorized)
		return
	} else {
		fmt.Fprintf(w, "verified")
	}
	dirName := fmt.Sprintf("./Odevler/%v_%v_%v", ogrID, homeworkLesson, homeworkName)
	//check Odevler directory exists
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		//create directory
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Save the files to the directory
	for _, fileHeader := range odevFiles {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		out, err := os.Create(fmt.Sprintf("%v/%v", dirName, fileHeader.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func Odevler(w http.ResponseWriter, r *http.Request) {
	ogrID := r.URL.Query().Get("ogr_id")
	ogrCode := r.URL.Query().Get("ogr_code")

	if !AuthCheck(ogrID, ogrCode, w, r) {
		http.Error(w, "not verified", http.StatusUnauthorized)
		return
	} else {
	}
	//header json
	w.Header().Set("Content-Type", "application/json")
	// List all files in the "Odevler" directory and its subdirectories
	fileList := []string{}
	err := filepath.Walk("./Odevler", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, strings.ReplaceAll(strings.ReplaceAll(path, "\\", "/"), "Odevler/", ""))
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var filtered []string
	var odevList []string
	for _, v := range fileList {
		if strings.Contains(v, ogrID) {
			filtered = append(filtered, v)
		}
		odevList = append(odevList, strings.Split(v, "/")[0])
		// remove duplicates
		odevList = removeDuplicates(odevList)
	}
	type Odev struct {
		Lesson string   `json:"lesson"`
		Name   string   `json:"name"`
		Files  []string `json:"files"`
	}
	var datajson []Odev
	for _, v := range odevList {
		odev := Odev{}
		odev.Lesson = strings.Split(v, "_")[1]
		odev.Name = strings.Split(v, "_")[2]
		for _, f := range filtered {
			if strings.Contains(f, v) {
				odev.Files = append(odev.Files, strings.ReplaceAll(f, v+"/", ""))
			}
		}
		datajson = append(datajson, odev)
	}
	result, err := json.Marshal(datajson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(result))
}
func Edit(w http.ResponseWriter, r *http.Request) {
	// Edit an existing homework
	err := r.ParseMultipartForm(64 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	odevFiles := r.MultipartForm.File["files"]
	ogrID := r.FormValue("ogr_id")
	ogrName := r.FormValue("ogr_name")
	ogrCode := r.FormValue("ogr_code")
	homeworkLesson := r.FormValue("homework_lesson")
	homeworkName := r.FormValue("homework_name")
	oldHomeworkLesson := r.FormValue("homework_old_lesson")
	oldHomeworkName := r.FormValue("homework_old_name")
	removeFiles := r.FormValue("remove_files")
	//ESCAPEING
	ogrID = strings.ReplaceAll(ogrID, "..", "")
	ogrName = strings.ReplaceAll(ogrName, "..", "")
	homeworkLesson = strings.ReplaceAll(homeworkLesson, "..", "")
	homeworkName = strings.ReplaceAll(homeworkName, "..", "")
	removeFiles = strings.ReplaceAll(removeFiles, "..", "")
	ogrID = strings.ReplaceAll(ogrID, "/", "")
	ogrName = strings.ReplaceAll(ogrName, "/", "")
	homeworkLesson = strings.ReplaceAll(homeworkLesson, "/", "")
	homeworkName = strings.ReplaceAll(homeworkName, "/", "")
	removeFiles = strings.ReplaceAll(removeFiles, "/", "")

	if !AuthCheck(ogrID, ogrCode, w, r) {
		http.Error(w, "not verified", http.StatusUnauthorized)
		return
	} else {
		fmt.Fprintf(w, "verified")
	}
	//remove files from the directory if removeFiles is not empty
	oldDirName := fmt.Sprintf("./Odevler/%v_%v_%v", ogrID, oldHomeworkLesson, oldHomeworkName)
	dirName := fmt.Sprintf("./Odevler/%v_%v_%v", ogrID, homeworkLesson, homeworkName)
	//change the directory name
	// fmt.Println(oldDirName)
	// fmt.Println(dirName)
	err = os.Rename(oldDirName, dirName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//remove files from the directory
	var removeFilesList []string
	err = json.Unmarshal([]byte(removeFiles), &removeFilesList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// if dirname exists and removeFiles is not empty
	if _, err := os.Stat(dirName); err == nil && removeFiles != "" {
		for _, v := range removeFilesList {
			err := os.Remove(fmt.Sprintf("%v/%v", dirName, v))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	// Save new files to the directory
	for _, fileHeader := range odevFiles {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		out, err := os.Create(fmt.Sprintf("%v/%v", dirName, fileHeader.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// fmt.Println(odevFiles, ogrID, ogrName, ogrCode, homeworkLesson, homeworkName, removeFiles)
}
func removeDuplicates(input []string) []string {
	// Create a map to track unique elements
	uniqueMap := make(map[string]bool)
	var result []string

	// Iterate over the input slice
	for _, item := range input {
		// Check if the item is already in the map
		if !uniqueMap[item] {
			// If not, add it to the map and result slice
			uniqueMap[item] = true
			result = append(result, item)
		}
	}

	return result
}
