/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/zulridzwan/admstool/util"
	"github.com/spf13/cobra"
)

// fixCmd represents the fix command
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("fix called")
		var id []int
		var invalid bool

		if len(args)==0{
			fmt.Println("Please enter the ID number of the Profile to fix")
			os.Exit(1)
		} else if len(args)>2 {
			fmt.Println("Too many parameters")
			os.Exit(1)
		} else if len(args)==1 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid number %s %s\n",args[0],err)
				os.Exit(1)
			}
			id = append(id,i)
		} else {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid number %s %s\n",args[0],err)
				invalid=true
			} else {
				id = append(id,i)
			}
			
			i, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Printf("Invalid number %s %s\n",args[1],err)
				invalid=true
			} else {
				id = append(id,i)
			}
			
			if invalid {
				os.Exit(1)
			}			
		}
		
		util.Status(false)
		
		//Storing in map for searching by ID
		tblMap:= make(map[int]util.UserProfile)
		for _,entry:= range util.ListUserProfile{
			tblMap[entry.ID]=entry
		}
		//check if already exist
		if len(id)==1{
			id = append(id,id[0])
		}
		for i:=id[0];i<=id[1];i++ {
			obj,exist:=tblMap[i]
			if exist {
				//fmt.Println(obj)
				if strings.Contains(obj.Key,".bak"){
					fmt.Printf("\nWarning! Skipping key %s. It's name contains .bak\n",obj.Key)
					continue
				} else {
					if obj.Exist=="found"{
						fmt.Printf("\nThe Profile ID %d already exist in ADMS entry\n",i)
						continue
					} else {
						t:=time.Now()
						tf:=t.Format("01/02/2006 15:04:05")
						value:="true|"+tf
						util.CreateProperty(obj.Key,value)
						fmt.Println("Done. New property and value created.")	
					}
				}
			}
		}
		fmt.Println("Refreshing...")
		util.Status(true)

	},
}

func init() {
	rootCmd.AddCommand(fixCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fixCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fixCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
