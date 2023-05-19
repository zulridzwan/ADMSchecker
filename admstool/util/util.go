package util

import (
	"os"
	"golang.org/x/sys/windows/registry"
	"strings"
	"fmt"
)

type UserProfile struct{
	ID	int
	Key	string
	Val	string
	Exist	string
}

type UserProfileChanged struct{
	ID	string
	Key	string
}

var (
	ListUserProfile []UserProfile
	ListUserProfileChanged []UserProfileChanged
	ListSearchResult []UserProfile
)

func Status(printstat bool) {
	if ListUserProfile!=nil{
		ListUserProfile=nil
	}
	if ListUserProfileChanged!=nil{
		ListUserProfileChanged=nil
	}
	
	if printstat{
		fmt.Println("Reading HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList...")
	}
	k1, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		fmt.Println("HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList error:",err)
		os.Exit(1)
	}
	defer k1.Close()
	profiles, err := k1.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("k1 ReadSubKeyNames error:",err)
		os.Exit(1)
	}
	if printstat{
		fmt.Println("Result:")
		fmt.Println(profiles)
	}
	
	
	if printstat{
		fmt.Println("\n\nReading HKLM\\SOFTWARE\\WOW6432Node\\Microsoft\\ADMS-WMT...")
	}
	k2, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\ADMS-WMT`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		fmt.Println("HKLM\\SOFTWARE\\WOW6432Node\\Microsoft\\ADMS-WMT error:",err)
		os.Exit(1)
	}
	defer k2.Close()
	changedprofiles, err := k2.ReadValueNames(-1)
	if err != nil {
		fmt.Println("k2 ReadValueNames error:",err)
		os.Exit(1)
	}
	
	f:= "WmiUserProfileChanged-"
	l:= len(f)
	for _,c:= range changedprofiles{
		if strings.Contains(c,f){
			rn:= []rune(c)
			id:= string(rn[l:len(c)])
			item:= UserProfileChanged{id,c}
			ListUserProfileChanged = append(ListUserProfileChanged,item)
		}
	}
	if printstat{
		fmt.Println("Result:")
		fmt.Printf("%s\n\n",ListUserProfileChanged)
	}
	//Storing in map for searching by ID
	cpMap:= make(map[string]string)
	for _,cp:= range ListUserProfileChanged{
		cpMap[cp.ID]=cp.Key
	}
	
	for i,p := range profiles{
		subkey:= "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList\\"+p
		k3, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey, registry.QUERY_VALUE)
		if err != nil {
			fmt.Printf("%s error: HKLM\\%s\n",subkey,err)
			os.Exit(1)
		}
		user,_,err :=k3.GetStringValue("ProfileImagePath")
		if err != nil {
			fmt.Printf("Warning! Skipping HKLM\\%s\n",subkey)
			fmt.Printf("ProfileImagePath: %s\n",err)
			continue
			//os.Exit(1)
		}
		if strings.Contains(user,":\\Users\\"){
			var result string
			_,exist:=cpMap[p]
			if exist {
				result="found"
			} else {
				result="missing"
			}
			item:=UserProfile{i,p,user,result}
			ListUserProfile = append(ListUserProfile,item)
		}
		k3.Close()
	}
	if printstat{
		//fmt.Println(ListUserProfile)
		fmt.Printf("\n\n%-6s%-50s%-50s%-20s\n","ID","Key","User","ADMS Entry")
		for _,u:= range ListUserProfile{
			fmt.Printf("%-6d%-50s%-50s%-20s\n",u.ID,u.Key,u.Val,u.Exist)
		}
	}
}

func CreateProperty(key string, val string) {
	fmt.Println("Accessing key HKLM\\SOFTWARE\\WOW6432Node\\Microsoft\\ADMS-WMT...")
	f:= "WmiUserProfileChanged-"
	k1, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\ADMS-WMT`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		fmt.Println("k1 registry.OpenKey error:",err)
		os.Exit(1)
	}
	defer k1.Close()
	fmt.Printf("Creating property %s...\n",f+key)
	err=k1.SetStringValue(f+key,val)
	if err != nil {
		fmt.Println("k1 SetStringValue error:",err)
		os.Exit(1)
	}
	
}


func DeleteProperty(key string) {
	fmt.Println("Accessing key HKLM\\SOFTWARE\\WOW6432Node\\Microsoft\\ADMS-WMT...")
	f:= "WmiUserProfileChanged-"
	k1, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\ADMS-WMT`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		fmt.Println("k1 registry.OpenKey error:",err)
		os.Exit(1)
	}
	defer k1.Close()
	fmt.Printf("Deleting property %s...\n",f+key)
	err=k1.DeleteValue(f+key)
	if err != nil {
		fmt.Println("k1 DeleteValue error:",err)
		os.Exit(1)
	}
	fmt.Println("Done.")
	
}

func Search(keyword string) []UserProfile{
	ListSearchResult = nil
	Status(false)
	for _,p:= range ListUserProfile{
		if strings.Contains(strings.ToLower(p.Val),keyword){
			ListSearchResult = append(ListSearchResult,p)
		}
	}

	return ListSearchResult
}

