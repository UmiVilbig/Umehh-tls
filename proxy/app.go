package proxy

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fatih/color"

	"github.com/andybalholm/brotli"

	"crypto/tls"
	"net/http"
)

var chrome_86 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"
var chrome_89 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36"
var chrome_90 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"

func log(t string, request RequestData, message string) {
	if request.Debug == true {
		now := time.Now().Format("2006-01-02 15:04:05.000000")
		magenta := color.New(color.FgMagenta).SprintFunc()
		if t == "info" {
			cyan := color.New(color.FgCyan).SprintFunc()
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("[%s] [%s] > %s.\n", magenta(now), yellow(request.ID), cyan(message))
		}
		if t == "error" {
			red := color.New(color.FgRed).SprintFunc()
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("[%s] [%s] > %s.\n", magenta(now), yellow(request.ID), red(message))
		}
		if t == "success" {
			green := color.New(color.FgGreen).SprintFunc()
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("[%s] [%s] > %s.\n", magenta(now), yellow(request.ID), green(message))
		}
	}
}

func ProxyRequest(data []byte) ResponseData {

	var requestData RequestData
	// var allowRedirect = true
	// var timeout = 20
	// var err error

	// if err = json.Unmarshal(data, &requestData); err != nil {
	// 	return ResponseData{
	// 		Success: false,
	// 		Message: "Error: Invalid Request Data",
	// 	}
	// }

	var responseData = ResponseData{
		ID:     requestData.ID,
		Method: requestData.Method,
		URL:    requestData.URL,
	}

	// log("info", requestData, "Starting...")

	// if requestData.URL == "" {
	// 	log("error", requestData, "Missing Request Url")
	// 	return ResponseData{
	// 		ID:      requestData.ID,
	// 		Success: false,
	// 		Message: "Error: Missing Request Url",
	// 	}
	// }

	// if requestData.Method == "" {
	// 	log("error", requestData, "Missing Request Method")
	// 	return ResponseData{
	// 		ID:      requestData.ID,
	// 		Success: false,
	// 		Message: "Error: Missing Request Method",
	// 	}
	// }

	// if len(requestData.Headers) == 0 {
	// 	log("error", requestData, "Missing Request Headers")
	// 	return ResponseData{
	// 		ID:      requestData.ID,
	// 		Success: false,
	// 		Message: "Error: Missing Request Headers",
	// 	}
	// }

	// if requestData.Headers["User-Agent"] == "" {
	// 	log("error", requestData, "Missing UserAgent")
	// 	return ResponseData{
	// 		ID:      requestData.ID,
	// 		Success: false,
	// 		Message: "Error: Missing UserAgent",
	// 	}
	// }

	// if requestData.Redirect == false {
	// 	allowRedirect = false
	// }

	// if requestData.Timeout != "" {
	// 	timeout, err = strconv.Atoi(requestData.Timeout)
	// 	if err != nil {
	// 		timeout = 20
	// 	}
	// }

	// if timeout > 60 {
	// 	log("error", requestData, "Timeout Cannot Be Longer Than 60 Seconds")
	// 	return ResponseData{
	// 		ID:      requestData.ID,
	// 		Success: false,
	// 		Message: "Error: Timeout Cannot Be Longer Than 60 Seconds",
	// 	}
	// }

	// var tlsClient tls.ClientHelloID

	// if strings.Contains(strings.ToLower(requestData.Headers["User-Agent"]), "chrome") {
	// 	tlsClient = tls.HelloChrome_Auto
	// } else if strings.Contains(strings.ToLower(requestData.Headers["User-Agent"]), "firefox") {
	// 	tlsClient = tls.HelloFirefox_Auto
	// } else {
	// 	tlsClient = tls.HelloIOS_Auto
	// }

	// client, err := cclient.NewClient(tlsClient, requestData.Proxy, allowRedirect, time.Duration(timeout))

	client := &http.Client{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			// NextProtos: []string{"h2"},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
			},
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		},
	}

	client.Transport = tr

	var req *http.Request

	// req, err = http.NewRequest(requestData.Method, requestData.URL, bytes.NewBuffer([]byte(requestData.Body)))
	req, _ = http.NewRequest(requestData.Method, requestData.URL, nil)

	req.Header.Set("User-Agent", chrome_90)

	// Set the Accept headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)

	if err != nil {
		log("error", requestData, "Failed Request "+err.Error())
		return ResponseData{
			ID:      requestData.ID,
			Success: false,
			Message: "Error: Failed Request " + err.Error(),
		}
	}

	defer resp.Body.Close()

	responseData.Headers = map[string][]interface{}{}

	for k, v := range resp.Header {
		if k != "Content-Length" && k != "Content-Encoding" {
			for _, kv := range v {
				responseData.Headers[k] = append(responseData.Headers[k], kv)
			}
		}
	}

	responseData.StatusCode = resp.StatusCode

	encoding := resp.Header["Content-Encoding"]

	body, err := ioutil.ReadAll(resp.Body)

	finalres := ""

	if err != nil {
		log("error", requestData, "Failed Request - Getting Content")
		return ResponseData{
			ID:      requestData.ID,
			Success: false,
			Message: "Error: Failed Request - Getting Content",
		}
	}

	finalres = string(body)

	if len(encoding) > 0 {
		if encoding[0] == "gzip" {
			unz, err := gUnzipData(body)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else if encoding[0] == "deflate" {
			unz, err := enflateData(body)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else if encoding[0] == "br" {
			unz, err := unBrotliData(body)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else {
			finalres = string(body)
		}
	} else {
		finalres = string(body)
	}

	responseData.Success = true
	responseData.Body = finalres

	log("success", requestData, "Request Successfully Proxied")

	return responseData
}

func gUnzipData(data []byte) (resData []byte, err error) {
	gz, _ := gzip.NewReader(bytes.NewReader(data))
	defer gz.Close()
	respBody, err := ioutil.ReadAll(gz)
	return respBody, err
}
func enflateData(data []byte) (resData []byte, err error) {
	zr, _ := zlib.NewReader(bytes.NewReader(data))
	defer zr.Close()
	enflated, err := ioutil.ReadAll(zr)
	return enflated, err
}
func unBrotliData(data []byte) (resData []byte, err error) {
	br := brotli.NewReader(bytes.NewReader(data))
	respBody, err := ioutil.ReadAll(br)
	return respBody, err
}
