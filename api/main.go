package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var SITE = "http://0.0.0.0:80"
var API = "http://0.0.0.0:8080"
var PORT = ":8080"
var KEY = ""
var CERT = ""

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("!!!Error loading .env file!!!")
		return
	}
	SITE = os.Getenv("SITE")
	API = os.Getenv("API")
	PORT = os.Getenv("PORT")
	KEY = os.Getenv("KEY")
	CERT = os.Getenv("CERT")
	fmt.Printf("API: %s\n SITE: %s\n PORT: %s\n", API, SITE, PORT)

	http.HandleFunc("/gallery", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir("./public/gallery")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var folders map[string][]string
		folders = make(map[string][]string)

		for _, file := range files {
			if file.IsDir() {
				folderName := file.Name()
				folderFiles, err := os.ReadDir("./public/gallery/" + folderName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				var fileList []string
				for _, folderFile := range folderFiles {
					fileList = append(fileList, folderFile.Name())
				}
				folders[folderName] = fileList
			}
		}

		jsonData, err := json.Marshal(folders)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonData))
	})
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		//read body jsonstring and convert to struct
		type Student struct {
			Name string `json:"ogr_name"`
			ID   string `json:"ogr_id"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var student Student
		err = json.Unmarshal(body, &student)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var code = GenerateVerifyCode()
		//read ./codes.json
		codesData, err := os.ReadFile("./codes.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var codes map[string]string
		err = json.Unmarshal(codesData, &codes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		codes[student.ID] = code
		//write ./codes.json
		codesData, err = json.Marshal(codes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = os.WriteFile("./codes.json", codesData, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//target student.id+@mehmetakif.edu.tr
		SendVerifyMail(student.ID+"@ogr.mehmetakif.edu.tr", code, student.Name)
		fmt.Println("Sending mail to "+student.ID, " with code: "+code, " and name: "+student.Name)

	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		//get query params
		code := r.URL.Query().Get("code")
		id := r.URL.Query().Get("id")
		//read ./codes.json
		codesData, err := os.ReadFile("./codes.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var codes map[string]string
		err = json.Unmarshal(codesData, &codes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if codes[id] == code {
			fmt.Fprint(w, `{"msg":"OK"}`)
		} else {
			fmt.Fprint(w, `{"msg":"FAIL"}`)
		}
	})
	http.HandleFunc("/odev", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Sadece POST istekleri kabul edilir.", http.StatusMethodNotAllowed)
			return
		} else {
			fmt.Println("POST istegi geldi.")
		}

		err := r.ParseMultipartForm(64 * 1024 * 1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(r.Form) // Form verilerini yazdır
		odevFiles := r.MultipartForm.File["odev_files"]
		ogrID := r.FormValue("ogr_id")
		ogrName := r.FormValue("ogr_name")
		dersName := r.FormValue("ders_name")
		verifyCode := r.FormValue("verify_code")
		codesData, err := os.ReadFile("./codes.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var codes map[string]string
		err = json.Unmarshal(codesData, &codes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// fmt.Fprint(w, "json:"+codes[ogrID]+" form:"+verifyCode)
		if codes[ogrID] != verifyCode {
			http.Error(w, "Kod yanlis!", http.StatusUnauthorized)
			return
		} else {
			fmt.Println("Kod dogru!")
		}
		//remove code from ./codes.json
		delete(codes, ogrID)
		codesData, err = json.Marshal(codes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = os.WriteFile("./codes.json", codesData, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//write files to ./public/ders_name_ogr_name_ogr_id/
		for i := 0; i < len(odevFiles); i++ {
			file, err := odevFiles[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()
			err = os.MkdirAll("./public/"+dersName+"_"+ogrName+"_"+ogrID, os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			targetFile, err := os.Create("./public/" + dersName + "_" + ogrName + "_" + ogrID + "/" + odevFiles[i].Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer targetFile.Close()
			_, err = file.Seek(0, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = file.Seek(0, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = io.Copy(targetFile, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	})

	// Enable CORS	})

	cors := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	http.Handle("/", cors(http.FileServer(http.Dir("../dist"))))
	//file is exist
	if _, err := os.Stat(CERT); os.IsNotExist(err) {
		fmt.Println("CERT not fount! Server started at http://localhost:8080")
		http.ListenAndServe(":8080", cors(http.DefaultServeMux))
		return
	}

	err = http.ListenAndServeTLS(":8443", CERT, KEY, cors(http.DefaultServeMux))
	if err != nil {
		panic("ListenAndServeTLS: " + err.Error())
	}
}
