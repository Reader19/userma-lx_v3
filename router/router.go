package router

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"userma-lx/config"
	"userma-lx/protocol"
	"userma-lx/service"
	"userma-lx/utils"
)

var Pwd, _ = os.Getwd()
var templatelogin *template.Template
var templateprofile *template.Template

func init() {
	templatelogin = template.Must(template.ParseFiles(Pwd + "/template/login.html"))
	templateprofile = template.Must(template.ParseFiles(Pwd + "/template/profile.html"))
}

func index(w http.ResponseWriter, r *http.Request) {

}

////////////////
//login
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req, err := loginPreVali(w, r)
		if err != nil {
			log.Println()
			return
		}
		resq, token, err := service.DoLogin(req)
		if err != nil {
			templatelogin.Execute(w, "login again")
			return
		}
		fmt.Println(resq)
		setCache(w, resq, token)
		templateprofile.Execute(w, resq)
		//w.Write([]byte("hello"))
	} else {
		templatelogin.Execute(w, "login again")
	}
}

//Pre-validation of user accounts
func loginPreVali(w http.ResponseWriter, r *http.Request) (protocol.ReqLogin, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	resq := protocol.ReqLogin{
		username,
		password,
	}
	if username == "" || password == "" {
		templatelogin.Execute(w, "login again")
		return resq, errors.New("username or password cannot be empty")
	}
	return resq, nil
}

func setCache(w http.ResponseWriter, resq protocol.ResqLogin, token string) {
	//cookies := make([]http.Cookie, 3)
	//cookies[0] = http.Cookie{Name: "username", Value: resq.UserName, MaxAge: config.MaxExTime}
	//cookies[1] = http.Cookie{Name: "nickname", Value: resq.NickName, MaxAge: config.MaxExTime}
	//cookies[2] = http.Cookie{Name: "picname", Value: resq.PicName, MaxAge: config.MaxExTime}
	//http.SetCookie(w, &cookies[0])
	//http.SetCookie(w, &cookies[1])
	//http.SetCookie(w, &cookies[2])
	cookie := http.Cookie{Name: "username", Value: resq.UserName, MaxAge: config.MaxExTime}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "token", Value: token, MaxAge: config.MaxExTime}
	http.SetCookie(w, &cookie)
	log.Println("success setting cookies")
}

////////////////
//profile
func getprofile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("one GetProfile request")
		ok, username := verifylogin(w, r)
		if !ok {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		resp, err := service.GetUserInfoCache(username)
		if err != nil {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		templateprofile.Execute(w, resp)
	}
}

func verifylogin(w http.ResponseWriter, r *http.Request) (bool, string) {
	username, err := r.Cookie("username")
	if err != nil {
		return false, ""
	}
	token, _ := r.Cookie("token")
	ok := service.VerifyToken(username.Value, token.Value)
	return ok, username.Value
}

//////////////////
//signout
func signOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "username", Value: "", MaxAge: -1}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "token", Value: "", MaxAge: -1}
	http.SetCookie(w, &cookie)
	templatelogin.Execute(w, "Success and Please login")

}

////////////////////
//signup
func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req, err := signUpPreVali(w, r)
		if err != nil {
			log.Println("incorrect username or password")
			return
		}
		log.Println(req)
		newErr := service.DoSignUp(req)
		if newErr != nil {
			fmt.Println("fail SignUp")
			templatelogin.Execute(w, "SignUp again")
		} else {
			fmt.Println("success")
			templatelogin.Execute(w, "Success and Please login")
		}
	}
}

func signUpPreVali(w http.ResponseWriter, r *http.Request) (protocol.ReqLogin, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		templatelogin.Execute(w, "SignUp again")
		return protocol.ReqLogin{}, errors.New("incorrect username or password")
	}
	req := protocol.ReqLogin{
		username,
		password,
	}
	return req, nil
}

func updateNickName(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ok, _ := verifylogin(w, r)
		if !ok {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		nickname := r.FormValue("nickname")
		username := r.FormValue("username")
		resp, _ := service.UpdateNickname(username, nickname)
		templateprofile.Execute(w, resp)
	}
	//if err != nil {
	//	templateprofile.Execute(w, nil)
	//}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ok, _ := verifylogin(w, r)
		if !ok {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		username := r.FormValue("username")
		file, header, err := r.FormFile("image")
		log.Println(header.Filename)
		if err != nil {
			log.Println("fail to get the img")
			resp := service.GetUserInfonByname(username)
			templateprofile.Execute(w, resp)
			return
		}
		defer file.Close()
		// check legal of file
		newName, isLegal := utils.CheckAndCreateFileName(header.Filename)
		//username := r.FormValue("username")
		if !isLegal {
			log.Println("img is no legal")
			resp := service.GetUserInfonByname(username)
			templateprofile.Execute(w, resp)
			return
		}
		filePath := config.ResourceImg + newName
		fileName := newName
		log.Println(filePath)
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		defer dstFile.Close()

		_, err = io.Copy(dstFile, file)
		if err != nil {
			log.Println("fail to copy imgFile: ", err)
			resp := service.GetUserInfonByname(username)
			templateprofile.Execute(w, resp)
			return
		}

		resp, _ := service.UploadFile(username, fileName)
		templateprofile.Execute(w, resp)
	}
}

func Router() {
	http.HandleFunc("/", getprofile)
	http.HandleFunc("/profile", getprofile)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signOut", signOut)
	http.HandleFunc("/signUp", signUp)
	http.HandleFunc("/updateNickName", updateNickName)
	http.HandleFunc("/uploadFile", uploadFile)
}
