package main
import (
	"fmt"
	"os"
	"io"
	"strings"
	"os/exec"
	"bytes"
)
func main() {
	var FCIV_HOME 		string
	var fileName 		string
	var FCIV_TARGET 	string
	var targetFilePath 	string
	var targetFileName	string
	var isSuccess		bool
	var hashDir			string

	if len(os.Args) == 1 {
		fmt.Println(" FCIV yapılacak dosya ismini belirtiniz!")
		os.Exit(0)
	}

	isSuccess 	= false
	fileName 	= os.Args[1]
	FCIV_HOME 	= os.Getenv("FCIV_HOME")
	FCIV_TARGET = os.Getenv("FCIV_TARGET")

	if FCIV_HOME == "" {
		fmt.Println("Environment Variables altında FCIV_HOME tanımı bulunamadı!")
		os.Exit(0)
	}
	if FCIV_TARGET == "" {
		fmt.Println("Environment Variables altında FCIV_TARGET tanımı bulunamadı!\n" +
			"FCIV yapılacak dosyanın bulundu ve klasörün oluşturulacağı yol bu değişkenle tanımlanmalıdır.")
		os.Exit(0)
	}
	targetFilePath = FCIV_TARGET + "\\";
	if strings.ContainsAny(fileName, "\\") {
		arr := strings.Split(fileName, "\\")
		targetFileName = arr[len(arr)-1]
		targetFilePath += targetFileName
	} else {
		targetFileName = fileName
		targetFilePath += targetFileName
	}
	_, err := fileCopy(fileName, targetFilePath)
	if err != nil {
		fmt.Println("Dosya kopylamada hata oldu.:", err)
		os.Exit(0)
	}

	cmd := exec.Command("cmd", "/C", "fciv", fileName)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

	arrSonuc := strings.Split(out.String(),"\n")
	arrLen :=len(arrSonuc)

	for i:=0;i<arrLen;i++ {

		if len(arrSonuc[i])>10 && !strings.ContainsAny(arrSonuc[i],"//"){
			arrHash := strings.Split(arrSonuc[i]," ")

			hashDir = FCIV_TARGET+"\\"+arrHash[0]

			if _, err := os.Stat(hashDir); os.IsNotExist(err) {
				os.Mkdir(hashDir,os.ModeDir)
			}

			_, err := fileCopy(targetFilePath,hashDir+"\\"+targetFileName)
			if err != nil {
				fmt.Printf("%s kopyalanamadı! Hata: %s", targetFilePath,err)
				os.Exit(0)
			}else{
				isSuccess  =true
			}
		}

	}
	if isSuccess {
		os.Remove(targetFilePath)
		fmt.Println("FCIV islemi tamamlandı.")
		fmt.Println(hashDir+"\\"+targetFileName)
		os.Exit(0)
	}
}
func fileCopy(source,target string)(int64, error){
	src_file, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return 0, err
	}

	if !src_file_stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", source)
	}

	dst_file, err := os.Create(target)
	if err != nil {
		return 0, err
	}
	defer dst_file.Close()
	return io.Copy(dst_file, src_file)
}