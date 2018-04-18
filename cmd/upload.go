// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)
type SessionResponse struct {
	Status struct {
		Code      int       `json:"code"`
		Text      string    `json:"text"`
		TimeStamp time.Time `json:"timeStamp"`
	} `json:"status"`
	Data []struct {
		SessionName string `json:"sessionName"`
		SessionID   string `json:"sessionId"`
	} `json:"data"`
}

type SiteInfo struct {
	Status struct {
		Code      int       `json:"code"`
		Text      string    `json:"text"`
		TimeStamp time.Time `json:"timeStamp"`
	} `json:"status"`
	Data []struct {
		Site struct {
			Subdomain       string        `json:"subdomain"`		
		} `json:"site"`
	} `json:"data"`
}

var token string
var siteID string
var upload string

func uploadToSite(sessionId string, fileToUpload string, Subdomain string){
	siteAPIURL := fmt.Sprintf("http://%s.mono.net/api.php/files", Subdomain)
	file := filepath.Base(fileToUpload)
	fmt.Println("Uploading " + file)
	resp, _ := resty.R().
		SetHeader("Cookie", "site_session=" + sessionId).
		SetFiles(map[string]string{
			"file_": fileToUpload,			
		}).
		SetFormData(map[string]string{
			"newFilename": file,			
		}).
		Post(siteAPIURL)
	
	if(resp.StatusCode() == 200){
		fmt.Println(file + " has been uploaded.")
	}else{
		fmt.Println("Upload failed, check that there is not already a file named " + file + "on the site.")
	}	
}

func getSession (token string, siteID string) *SessionResponse{
	body:=fmt.Sprintf("{\"userToken\": \"%s\",\"command\": \"apiLogin\",\"siteId\": %s}",token, siteID)
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&SessionResponse{}). 
		Post("https://hal.mono.net/api/v1/reseller/account/site")
	
		if(err!=nil){
		fmt.Println(err)
		os.Exit(1)
	}
	return  resp.Result().(*SessionResponse)		
}

func getSiteInfo (token string, siteID string) *SiteInfo{
	body := fmt.Sprintf("{\"userToken\": \"%s\",\"siteId\": %s,\"command\": \"getInfo\" }", token, siteID)
	siteInfoResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&SiteInfo{}). 
		Post("https://hal.mono.net/api/v1/reseller/account/site")
	
	if(err!=nil){
		fmt.Println(err)
		os.Exit(1)
	}
	return siteInfoResp.Result().(*SiteInfo)
}
// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a image or a directory of images",
	Long: `Upload a image or a directory of images`,
	Run: func(cmd *cobra.Command, args []string) {
		if(token==""){
			fmt.Println("You must supply a token")
			os.Exit(1)
		}
		if(upload==""){
			fmt.Println("You must supply a directory or file to upload")
			os.Exit(1)
		}

		if(upload==""){
			fmt.Println("You must supply a siteID")
			os.Exit(1)
		}
		
		auth := getSession(token, siteID)
		siteInfo := getSiteInfo(token, siteID)			
	
		fi, err := os.Stat(upload)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch mode := fi.Mode(); {
			case mode.IsDir():
				files, err := ioutil.ReadDir(upload)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}				
				for _, f := range files {			
					srcFile:=fmt.Sprintf("%s/%s", upload, f.Name())				
					uploadToSite(auth.Data[0].SessionID, srcFile, siteInfo.Data[0].Site.Subdomain)
				}
			case mode.IsRegular():
				uploadToSite(auth.Data[0].SessionID, upload, siteInfo.Data[0].Site.Subdomain)	
		}		

	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&token, "token", "t", "", "Your token")
	uploadCmd.Flags().StringVarP(&siteID, "siteid", "s", "", "Site id")
	uploadCmd.Flags().StringVarP(&upload, "upload", "u", "", "directory or file to upload")
}
