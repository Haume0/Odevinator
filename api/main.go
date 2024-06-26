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
	"regexp"
	"strconv"
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
	var input string
	//get input from user
	fmt.Println("\nSistemi internet üzerinden erişilir hale getirir.")
	fmt.Println("\n⚠️ İNTERNET BAĞLANTISI GEREKLİDİR.")
	fmt.Print("\n🌐 Global URL ister misiniz? [yes/no]: ")
	fmt.Scanln(&input)
	//if input is yes, set the --global flag
	if input == "yes" {
		os.Args = append(os.Args, "--global")
	}
	configJSON, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println("config.json dosyası bulunamadı.")
		os.Exit(1)
	}
	var config map[string]interface{}
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		fmt.Println("config.json dosyası okunamadı.")
		os.Exit(1)
	}
	PASS = config["PASS"].(string)
	MAIL = config["MAIL"].(string)
	OKUL_SUFFIX = config["OKUL_SUFFIX"].(string)
	//remove the auth.json file if it exists and create a new one
	if _, err := os.Stat("./auth.json"); os.IsNotExist(err) {
		//create file with {} in it
		file, err := os.Create("./auth.json")
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err.Error())
		}
		defer file.Close()
		file.WriteString("[]")
	} else {
		os.Remove("./auth.json")
		//create file with {} in it
		file, err := os.Create("./auth.json")
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err.Error())
		}
		defer file.Close()
		file.WriteString("[]")

	}
	// devmode env logging
	if os.Args[len(os.Args)-1] == "--dev" {
		fmt.Printf(`
		MAIL_SUFFIX: %v
		MAIL: %v
		PASS: %v`, OKUL_SUFFIX, MAIL, PASS)
	}
	// ROUTES
	http.HandleFunc("/auth", Login)
	http.HandleFunc("/verify", Verify)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/odevler", Odevler)
	// Enable CORS
	cors := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS,")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	// // Serve dist file
	// http.Handle("/", cors(http.FileServer(http.Dir("./dist"))))
	// // Serve client-side routing index.html
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./dist/index.html")
	// })

	fileServer := http.FileServer(http.Dir("./dist"))
	fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !fileMatcher.MatchString(r.URL.Path) {
			http.ServeFile(w, r, "./dist/index.html")
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	var allLinks = []string{}
	// If --global flag is set, create a tunnel and get the global URL
	if os.Args[len(os.Args)-1] == "--global" {
		go func() {
			// A tunnel service to create a tunnel to the local server
			// We use the 'ssh' command to create a tunnel to serveo.net
			// -R flag specifies the remote port forwarding
			// 80:localhost:8080 maps port 80 on serveo.net to port 8080 on the local machine
			// serveo.net is a free tunnel service
			cmd := exec.Command("ssh", "-R", "80:localhost:8080", "nokey@localhost.run")
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
				//regex for global url  https://4f84672b874728.lhr.life
				urlRegex := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.lhr\.life`)
				globalUrl := urlRegex.FindString(scanner.Text())
				if globalUrl != "" {
					// Print the global URL
					fmt.Println("")
					allLinks = append(allLinks, globalUrl)
					fmt.Printf("🌐 Global URL: %v\n", globalUrl)
					// Clearing all cli styling
					var qrLink = "http://localhost:8080/qr"
					for i, link := range allLinks {
						if i == 0 {
							qrLink = qrLink + "?" + strconv.Itoa(i) + "=" + link
						} else {
							qrLink = qrLink + "&" + strconv.Itoa(i) + "=" + link
						}
					}
					fmt.Println("")
					fmt.Printf("\r")
					fmt.Printf("\033[0m") // ANSI renk kodlarını sıfırla
					fmt.Printf("🏁 QR Kodlar: %v\n", qrLink)
					fmt.Printf("\033[0m") // ANSI renk kodlarını sıfırla
					fmt.Printf("\r")      // Satırı temizle
					break
				} else {
					fmt.Println("Hata: Global URL bulunamadı.")
				}
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

🛜  UYARI: Bu ikona sahip baglantılar için
aynı ağa bağlı olamlısınız!
🌐 Bu ikona sahip baglantılar ise internet üzerinden erişilebilir.
`)
	fmt.Println("")
	//get local router network ip adress ethernet and wifi and print them like 🔗 http://adress:8080
	// Get the local IP address of the machine
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	// Loop through the addresses
	for _, addr := range addrs {
		// Check if the address is an IP address
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// Check if the IP address is an IPv4 address
			if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
				// Print the local URL
				allLinks = append(allLinks, fmt.Sprintf("http://%v:8080", ipnet.IP))
				fmt.Printf("🛜  Local URL: http://%v:8080\n", ipnet.IP)
			}
		}
	}
	if os.Args[len(os.Args)-1] != "--global" {
		var qrLink = "http://localhost:8080/qr"
		for i, link := range allLinks {
			if i == 0 {
				qrLink = qrLink + "?" + strconv.Itoa(i) + "=" + link
			} else {
				qrLink = qrLink + "&" + strconv.Itoa(i) + "=" + link
			}
		}
		fmt.Println("")
		fmt.Printf("\r")
		fmt.Printf("\033[0m") // ANSI renk kodlarını sıfırla
		fmt.Printf("🏁 QR Kodlar: %v\n", qrLink)
	}

	err = http.ListenAndServe(":8080", cors(http.DefaultServeMux))
	if err != nil {
		panic("ListenAndServeTLS: " + err.Error())
	}
}
