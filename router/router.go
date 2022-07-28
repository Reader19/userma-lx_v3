//package router
package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"userma-lx/config"
	"userma-lx/protocol"
	"userma-lx/rpc"
	"userma-lx/utils"
)

var Pwd, _ = os.Getwd()
var templatelogin *template.Template
var templateprofile *template.Template
var client rpc.TcpClient

func init() {
	templatelogin = template.Must(template.ParseFiles(Pwd + "/template/login.html"))
	templateprofile = template.Must(template.ParseFiles(Pwd + "/template/profile.html"))
	//var err error
	//client, err = rpc.NewClient(1, config.TcpServerAddr)
	//if err != nil {
	//	log.Println(err)
	//}
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
		var resqAll protocol.RespLogin
		log.Println("before call: ", resqAll)
		err = client.Call("DoLogin", req, &resqAll)
		err = transError(err)
		log.Println("after call: ", resqAll)
		if err != nil {
			log.Println("err: ", err)
			templatelogin.Execute(w, "login again")
			return
		}
		resq := protocol.RespProfile{
			resqAll.UserName,
			resqAll.NickName,
			resqAll.PicName,
		}
		token := resqAll.Token
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

func setCache(w http.ResponseWriter, resq protocol.RespProfile, token string) {
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
		ok, username := verifylogin(r)
		log.Println("first")
		if !ok {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		var resp protocol.RespProfile
		log.Println(resp)
		err := client.Call("GetUserInfoCache", username, &resp)
		err = transError(err)
		log.Println(resp)
		if err != nil {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		templateprofile.Execute(w, resp)
	}
}

func verifylogin(r *http.Request) (bool, string) {
	username, err := r.Cookie("username")
	if err != nil {
		return false, ""
	}
	token, _ := r.Cookie("token")
	req := protocol.ReqVerifyToken{
		UserName: username.Value,
		Token:    token.Value,
	}
	var ok bool
	_ = client.Call("VerifyToken", req, &ok)
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
		var newErr error
		client.Call("DoSignUp", req, &newErr)
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
		ok, _ := verifylogin(r)
		if !ok {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		nickname := r.FormValue("nickname")
		username := r.FormValue("username")
		rep := protocol.ReqSetNickName{
			UserName: username,
			NickName: nickname,
		}
		var resp protocol.RespProfile
		err := client.Call("UpdateNickName", rep, &resp)
		err = transError(err)
		if err != nil {
			log.Println("fail to copy UpdateNickName: ", err)
			_ = client.Call("GetUserInfoByName", username, &resp)
			templateprofile.Execute(w, resp)
			return
		}
		templateprofile.Execute(w, resp)
	}
	//if err != nil {
	//	templateprofile.Execute(w, nil)
	//}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ok, _ := verifylogin(r)
		if !ok {
			log.Println("Account expired")
			templatelogin.Execute(w, "Please login")
			return
		}
		username := r.FormValue("username")
		file, header, err := r.FormFile("image")
		log.Println(header.Filename)
		var resp protocol.RespProfile
		if err != nil {
			log.Println("fail to get the img")
			_ = client.Call("GetUserInfoByName", username, &resp)
			templateprofile.Execute(w, resp)
			return
		}
		defer file.Close()
		// check legal of file
		newName, isLegal := utils.CheckAndCreateFileName(header.Filename)
		//username := r.FormValue("username")
		if !isLegal {
			log.Println("img is no legal")
			_ = client.Call("GetUserInfoByName", username, &resp)
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
			_ = client.Call("GetUserInfoByName", username, &resp)
			templateprofile.Execute(w, resp)
			return
		}

		req := protocol.ReqUploadFile{
			UserName: username,
			FileName: fileName,
		}
		_ = client.Call("UploadFile", req, &resp)
		templateprofile.Execute(w, resp)
	}
}

func transError(er error) error {
	if er.Error() == "" {
		return nil
	}
	return er
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

func main() {
	tcpServer := rpc.NewServer(config.TcpServerAddr)
	go tcpServer.Run()
	time.Sleep(1 * time.Second)

	var err error
	client, err = rpc.NewClient(2, config.TcpServerAddr)
	if err != nil {
		log.Println(err)
	}

	server := http.Server{
		Addr: "localhost:8080",
	}
	Router()
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
