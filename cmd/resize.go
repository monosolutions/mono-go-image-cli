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
	"strconv"
	"log"
	"io/ioutil"
	"github.com/spf13/cobra"
	"github.com/disintegration/imaging"
)


var file string
var directory string
var output string
var width string
var height string
var h int
var w int

var resizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize file or directory",
	Long: `Resize  a file or a directory of images`,
	Run: func(cmd *cobra.Command, args []string) {
		if(file=="" && directory == ""){
			fmt.Println("You have to either choose an image (-f) or a directory (-d) to resize")
			os.Exit(1)
		}
		if(height=="" && width == ""){
			fmt.Println("You have to set either height (-e 200) or width (-w 200)")
			os.Exit(1)
		}
		if(height!=""){
			h,_=strconv.Atoi(height)
		}else{
			h=0
		}

		if(width!=""){
			w,_=strconv.Atoi(width)
		}else{
			w=0
		}

		if(output==""){
			fmt.Println("You have to set an output either a folder or a file")
			os.Exit(1)
		}

		fmt.Printf("Resizing %s and storing it to %s \r\n", file, output)
		
		if(file!=""){
			src, err := imaging.Open(file)
			if err != nil {
				log.Fatalf("failed to open image: %v", err)
			}
			src = imaging.Resize(src, w, h, imaging.Lanczos)
			err = imaging.Save(src, output)
			if err != nil {
				log.Fatalf("failed to save image: %v", err)
			}
			os.Exit(0)
		}

		if(directory!=""){
			files, err := ioutil.ReadDir(directory)
			if err != nil {
				log.Fatal(err)
			}
			_ = os.Mkdir(output, 0777)
			for _, f := range files {
				dest:=fmt.Sprintf("%s/%s",output, f.Name())
				srcFile:=fmt.Sprintf("%s/%s", directory, f.Name())		

				src, err := imaging.Open(srcFile)
				if err != nil {
					log.Fatalf("failed to open image: %v", err)
				}
				src = imaging.Resize(src, w, h, imaging.Lanczos)
				err = imaging.Save(src, dest)
				if err != nil {
					log.Fatalf("failed to save image: %v", err)
				}				
			}			
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(resizeCmd)
	resizeCmd.Flags().StringVarP(&file, "file", "f", "", "File to resize")
	resizeCmd.Flags().StringVarP(&directory, "dir", "d", "", "Directory to resize")
	resizeCmd.Flags().StringVarP(&output, "output", "o", "", "Output folder or file")
	resizeCmd.Flags().StringVarP(&height, "height", "e", "", "To which height should it scale?")
	resizeCmd.Flags().StringVarP(&width, "width", "w", "", "To which width should it scale?")
}
