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

// ENV VARS
var (
	PORT        = ":8080"
	PASS        string
	MAIL        string
	OKUL_SUFFIX string
)

// STRUCTS
type Student struct {
	Name     string `json:"ogr_name"`
	ID       string `json:"ogr_id"`
	ClientID string `json:"client_id"`
}

func main() {
	// ENV VARS
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("!!!Error loading .env file!!!")
		return
	}
	PASS = os.Getenv("PASS")
	MAIL = os.Getenv("MAIL")
	OKUL_SUFFIX = os.Getenv("OKUL_SUFFIX")
	// MAIN
	// devmode env logging
	if os.Args[len(os.Args)-1] == "--dev" {
		fmt.Printf(`
		MAIL_SUFFIX: %v
		MAIL: %v
		PASS: %v`, OKUL_SUFFIX, MAIL, PASS)
	}
	// POST /verify -> generate verify code, sending and priting it
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		//read body jsonstring and convert to struct
		// Read the request body and convert it to a Student struct
		body, err := io.ReadAll(r.Body)
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the request body into a Student struct
		var student Student
		err = json.Unmarshal(body, &student)
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//save all the requests to the log.json in a client_id:[{ogr_id:number,ogr_name:name}] format
		// if there is no log.json file, create it
		if _, err := os.Stat("./log.json"); os.IsNotExist(err) {
			// create a new log.json file with {} in it
			os.Create("./log.json")
		}
		logData, err := os.ReadFile("./log.json")
		if err != nil {
			http.Error(w, err.Error()+"1", http.StatusInternalServerError)
			return
		}
		//if log.json is not a json file, make it a json file
		if !strings.HasPrefix(string(logData), "{") {
			logData = []byte("{}")
		}
		var log map[string][]map[string]string
		err = json.Unmarshal(logData, &log)
		if err != nil {
			http.Error(w, err.Error()+"2", http.StatusInternalServerError)
			return
		}
		log[student.ClientID] = append(log[student.ClientID], map[string]string{"ogr_id": student.ID, "ogr_name": student.Name})
		logData, err = json.Marshal(log)
		if err != nil {
			http.Error(w, err.Error()+"3", http.StatusInternalServerError)
			return
		}
		err = os.WriteFile("./log.json", logData, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error()+"4", http.StatusInternalServerError)
			return
		}
		// Generate a verify code
		var code = GenerateVerifyCode()
		//check is there codes.json file exist
		if _, err := os.Stat("./codes.json"); os.IsNotExist(err) {
			os.Create("./codes.json")
		}
		// Read the codes.json file
		codesData, err := os.ReadFile("./codes.json")
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//if codes.json is not a json file, make it a json file
		if !strings.HasPrefix(string(codesData), "{") {
			codesData = []byte("{}")
		}

		// Unmarshal the codes.json file into a map of strings
		var codes map[string]string
		err = json.Unmarshal(codesData, &codes)
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add the code to the map with the student's ID as the key
		codes[student.ID] = code

		// Marshal the codes map back into JSON
		codesData, err = json.Marshal(codes)
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the updated codes.json file
		err = os.WriteFile("./codes.json", codesData, os.ModePerm)
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send a verify mail to the student with the code
		SendVerifyMail(student.ID+OKUL_SUFFIX, student.ID, code, student.Name)
	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		// Get the code and ID from the query parameters
		code := r.URL.Query().Get("code")
		id := r.URL.Query().Get("id")

		// Read the codes.json file
		codesData, err := os.ReadFile("./codes.json")
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the codes.json file into a map of strings
		var codes map[string]string
		err = json.Unmarshal(codesData, &codes)
		if err != nil {
			// In case of an error, return an internal server error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the content type of the response to JSON
		w.Header().Set("Content-Type", "application/json")

		// Check if the code matches the student's ID
		if codes[id] == code {
			// Return a JSON response indicating success
			fmt.Fprint(w, `{"msg":"OK"}`)
		} else {
			// Return a JSON response indicating failure
			fmt.Fprint(w, `{"msg":"FAIL"}`)
		}
	})
	http.HandleFunc("/odev", func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST requests
		if r.Method != "POST" {
			http.Error(w, "Sadece POST istekleri kabul edilir.", http.StatusMethodNotAllowed)
			return
		}
		// Read and parse the multipart form data
		err := r.ParseMultipartForm(64 * 1024 * 1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Print the form data
		fmt.Println(r.Form) // Form verilerini yazdır
		// Get the odev files from the request
		odevFiles := r.MultipartForm.File["odev_files"]
		ogrID := r.FormValue("ogr_id")
		ogrName := r.FormValue("ogr_name")
		dersName := r.FormValue("ders_name")
		verifyCode := r.FormValue("verify_code")
		// Read the codes.json file
		codesData, err := os.ReadFile("./codes.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Unmarshal the codes.json file into a map of strings
		var codes map[string]string
		err = json.Unmarshal(codesData, &codes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Check if the verify code is correct
		if codes[ogrID] != verifyCode {
			http.Error(w, "Kod yanlis!", http.StatusUnauthorized)
			return
		} else {
			fmt.Println("Kod dogru!")
		}
		// Remove the code from the codes.json file
		delete(codes, ogrID)
		// Marshal the updated codes map back into JSON
		codesData, err = json.Marshal(codes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Write the updated codes.json file
		err = os.WriteFile("./codes.json", codesData, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write files to ./odevler/ders_name_ogr_name_ogr_id/
		for i := 0; i < len(odevFiles); i++ {
			// Get the file from the request
			file, err := odevFiles[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Get the file name from the request
			defer file.Close()
			// Get the file name from the request
			err = os.MkdirAll("./odevler/"+dersName+"_"+ogrName+"_"+ogrID, os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Get the file name from the request
			targetFile, err := os.Create("./odevler/" + dersName + "_" + ogrName + "_" + ogrID + "/" + odevFiles[i].Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Get the file name from the request
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

	// Enable CORS
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
	// Serve Website
	http.Handle("/", cors(http.FileServer(http.Dir("./dist"))))
	// Get local IPv4 address
	var link string
	//get local ipv4 adress 192.168.1.33 etc.
	var ip net.IP
	ifaces, err := net.Interfaces()
	for _, i := range ifaces {
		if i.Name == "Ethernet" {
			addrs, err := i.Addrs()
			if err != nil {
				println(err)
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ip = ipnet.IP
					link = "http://" + ip.String() + PORT
					break
				}
			}
			break
		}
	}
	// If --global flag is set, create a tunnel and get the global URL
	if os.Args[len(os.Args)-1] == "--global" {
		go func() {
			// A tunnel service to create a tunnel to the local server
			// We use the 'ssh' command to create a tunnel to serveo.net
			// -R flag specifies the remote port forwarding
			// 80:localhost:8080 maps port 80 on serveo.net to port 8080 on the local machine
			// serveo.net is a free tunnel service
			cmd := exec.Command("ssh", "-R", "80:localhost:8080", "serveo.net")
			// Discard the standard error output of the command
			cmd.Stderr = io.Discard
			// Get the standard output of the command
			stdout, _ := cmd.StdoutPipe()
			// Start the command
			cmd.Start()
			// Create a scanner to read the output line by line
			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			// Loop through the output lines
			for scanner.Scan() {
				// Replace the prefix "Forwarding HTTP traffic from " with an empty string
				// to get the global URL
				var globalUrl = strings.Replace(scanner.Text(), "Forwarding HTTP traffic from ", "", 1)
				// Print the global URL
				fmt.Printf("🔗 Global URL: %v \n", globalUrl)
				//clearing all cli styling
				fmt.Printf("\033[0m") // ANSI renk kodlarını sıfırla
				fmt.Printf("\r")      // Satırı temizle
				break
			}
			// Wait for the command to finish
			cmd.Wait()
		}()
	}
	//description
	fmt.Printf(`
✨Ödevinatör✨ by Haume

🚀 Hazırız, aşağıdaki bağlantıyı öğrenciler
ile paylaşabilirsiniz.

⚠️ UYARI: AYNI WI-FI AĞINA BAĞLI OLMALISINIZ❗

🔗 %v
`, link)
	err = http.ListenAndServe(":8080", cors(http.DefaultServeMux))
	if err != nil {
		panic("ListenAndServeTLS: " + err.Error())
	}
}
