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
	"strings"
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

// func main() {
//     // Uygulamanın çalışma dizinini alın
//     cwd, err := os.Getwd()
//     if err != nil {
//         fmt.Println("Çalışma dizini alınamadı:", err)
//         os.Exit(1)
//     }
//     fmt.Println("Uygulama dizini:", cwd)

//     // Config dosyasının yolunu oluşturun
//     configPath := filepath.Join(cwd, "config.json")
//     fmt.Println("Config dosya yolu:", configPath)

//     // Config dosyasının varlığını kontrol edin
//     _, err = os.Stat(configPath)
//     if os.IsNotExist(err) {
//         fmt.Println("config.json dosyası bulunamadı.")
//         os.Exit(1)
//     }

//     // Config dosyasını okuyun
//     configJSON, err := os.ReadFile(configPath)
//     if err != nil {
//         fmt.Println("config.json dosyası okunamadı:", err)
//         os.Exit(1)
//     }

//     fmt.Println("config.json dosyası başarıyla okundu.")
//     // Gerekli işlemleri yapın
// }

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
	// MAIN
	fmt.Print(PASS, MAIL, OKUL_SUFFIX)
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
				fmt.Println("")
				fmt.Printf("🌐 Global URL: %v\n", globalUrl)
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

⚠️ UYARI: AYNI WI-FI AĞINA BAĞLI OLMALISINIZ 🛜
`)
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
				fmt.Printf("🔗 Local URL: http://%v:8080\n", ipnet.IP)
			}
		}
	}
	err = http.ListenAndServe(":8080", cors(http.DefaultServeMux))
	if err != nil {
		panic("ListenAndServeTLS: " + err.Error())
	}
}
