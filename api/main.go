package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

var (
	PORT        = ":8080"
	PASS        string
	MAIL        string
	OKUL_SUFFIX string
)

type Student struct {
	Name string `json:"ogr_name"`
	ID   string `json:"ogr_id"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("!!!Error loading .env file!!!")
		return
	}
	PASS = os.Getenv("PASS")
	MAIL = os.Getenv("MAIL")
	OKUL_SUFFIX = os.Getenv("OKUL_SUFFIX")
	if os.Args[len(os.Args)-1] == "--dev" {
		fmt.Printf(`
		MAIL_SUFFIX: %v
		MAIL: %v
		PASS: %v`, OKUL_SUFFIX, MAIL, PASS)
	}

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		//read body jsonstring and convert to struct

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
		SendVerifyMail(student.ID+OKUL_SUFFIX, code, student.Name)
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

		//write files to ./odevler/ders_name_ogr_name_ogr_id/
		for i := 0; i < len(odevFiles); i++ {
			file, err := odevFiles[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()
			err = os.MkdirAll("./odevler/"+dersName+"_"+ogrName+"_"+ogrID, os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			targetFile, err := os.Create("./odevler/" + dersName + "_" + ogrName + "_" + ogrID + "/" + odevFiles[i].Filename)
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
	http.Handle("/", cors(http.FileServer(http.Dir("./dist"))))
	var link string
	//get local ipv4 adress 192.168.1.33 etc.
	var ip net.IP
	ifaces, err := net.Interfaces()
	for idxx, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			println(err)
		}
		if idxx == 0 {
			ip = addrs[1].(*net.IPNet).IP
			link = "http://" + ip.String() + PORT
			// println(ip.String())
		}
	}
	if os.Args[len(os.Args)-1] == "--global" {
		go func() {
			cmd := exec.Command("ssh", "-R", "80:localhost:8080", "serveo.net")
			cmd.Stderr = io.Discard
			stdout, _ := cmd.StdoutPipe()
			cmd.Start()

			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				var globalUrl = strings.Replace(scanner.Text(), "Forwarding HTTP traffic from ", "", 1)
				fmt.Printf(`
	🔗 Global URL: %v`, globalUrl)
			}
			cmd.Wait()
		}()
	}
	//description
	fmt.Printf(`
	✨Ödevinatör✨ by Haume

	🚀 Hazırız, aşağıdaki bağlantıyı öğrenciler
	ile paylaşabilirsiniz.

	⚠️ UYARI: AYNI WI-FI AĞINA BAĞLI OLMALISINIZ❗

	🔗 %v`, link)
	err = http.ListenAndServe(":8080", cors(http.DefaultServeMux))
	if err != nil {
		panic("ListenAndServeTLS: " + err.Error())
	}
}
