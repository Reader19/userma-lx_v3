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
	"usermaLX4/config"
	"usermaLX4/protocol"
	"usermaLX4/rpcService"
	"usermaLX4/utils"
)

var Pwd, _ = os.Getwd()
var templateLogin *template.Template
var templateProfile *template.Template
var client rpcService.TcpClient

func init() {
	templateLogin = template.Must(template.ParseFiles(Pwd + "/template/login.html"))
	templateProfile = template.Must(template.ParseFiles(Pwd + "/template/profile.html"))
}

func Router() {
	http.HandleFunc("/", getProfile)
	http.HandleFunc("/profile", getProfile)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signOut", signOut)
	http.HandleFunc("/signUp", signUp)
	http.HandleFunc("/updateNickName", updateNickName)
	http.HandleFunc("/uploadFile", uploadFile)

	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir(Pwd+"/resource"))))
}

func main() {
	tcpServer := rpcService.NewServer(config.TcpServerAddr)
	go tcpServer.Run()
	time.Sleep(1 * time.Second)

	var err error
	client, err = rpcService.NewClient(config.MaxNumConn, config.TcpServerAddr)
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

////////////////
//login
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("one login request")
		req, err := loginPreVali(w, r)
		if err != nil {
			return
		}
		var respAll protocol.RespLogin
		log.Println("before call: ", respAll)
		err = client.Call("DoLogin", req, &respAll)
		err = transError(err)
		log.Println("after call: ", respAll)
		if err != nil {
			log.Println("err: ", err)
			templateLogin.Execute(w, "login again")
			return
		}
		resp := protocol.RespProfile{
			respAll.UserName,
			respAll.NickName,
			respAll.PicName,
		}
		token := respAll.Token
		fmt.Println(resp)
		setCache(w, resp, token)
		templateProfile.Execute(w, resp)
		//w.Write([]byte("hello"))
	} else {
		templateLogin.Execute(w, "login again")
	}
}

//Pre-validation of user accounts
func loginPreVali(w http.ResponseWriter, r *http.Request) (protocol.ReqLogin, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	resp := protocol.ReqLogin{
		username,
		password,
	}
	if username == "" || password == "" {
		templateLogin.Execute(w, "login again")
		return resp, errors.New("username or password cannot be empty")
	}
	return resp, nil
}

func setCache(w http.ResponseWriter, resq protocol.RespProfile, token string) {
	cookie := http.Cookie{Name: "username", Value: resq.UserName, MaxAge: config.MaxExTime}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "token", Value: token, MaxAge: config.MaxExTime}
	http.SetCookie(w, &cookie)
	log.Println("success setting cookies")
}

////////////////
//profile
func getProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("one GetProfile request")
		ok, username := verifylogin(r)
		log.Println("first")
		if !ok {
			log.Println("Account expired")
			templateLogin.Execute(w, "Please login")
			log.Println("templateLogin over")
			return
		}
		var resp protocol.RespProfile
		log.Println(resp)
		err := client.Call("GetUserInfoCache", username, &resp)
		err = transError(err)
		log.Println(resp)
		if err != nil {
			log.Println("Account expired")
			templateLogin.Execute(w, "Please login")
			return
		}
		templateProfile.Execute(w, resp)
	}
}

func verifylogin(r *http.Request) (bool, string) {
	username, err := r.Cookie("username")
	if err != nil {
		return false, ""
	}
	token, err := r.Cookie("token")
	if err != nil {
		return false, ""
	}
	req := protocol.ReqVerifyToken{
		UserName: username.Value,
		Token:    token.Value,
	}
	var ok bool
	log.Println("req which verTok: ", req)
	_ = client.Call("VerifyToken", req, &ok)
	return ok, username.Value
}

//////////////////
//signout
func signOut(w http.ResponseWriter, r *http.Request) {
	log.Println("one signOut request")
	cookie := http.Cookie{Name: "username", Value: "", MaxAge: -1}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "token", Value: "", MaxAge: -1}
	http.SetCookie(w, &cookie)
	templateLogin.Execute(w, "Success and Please login")

}

////////////////////
//signup
func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("one signUp request")
		req, err := signUpPreVali(w, r)
		if err != nil {
			log.Println("incorrect username or password")
			return
		}
		log.Println(req)
		var ok bool
		err = client.Call("DoSignUp", req, &ok)
		log.Println("输出登陆ok值： ", ok)
		log.Println("输出call返回错误值： ", err)
		if !ok {
			fmt.Println("fail SignUp")
			templateLogin.Execute(w, "SignUp again")
		} else {
			fmt.Println("success")
			templateLogin.Execute(w, "Success and Please login")
		}
	}
}

func signUpPreVali(w http.ResponseWriter, r *http.Request) (protocol.ReqLogin, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		templateLogin.Execute(w, "SignUp again")
		return protocol.ReqLogin{}, errors.New("incorrect username or password")
	}
	req := protocol.ReqLogin{
		username,
		password,
	}
	return req, nil
}

///////////////////
//updateNickName
func updateNickName(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("one updateNickName request")
		ok, _ := verifylogin(r)
		if !ok {
			log.Println("Account expired")
			templateLogin.Execute(w, "Please login")
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
			templateProfile.Execute(w, resp)
			return
		}
		templateProfile.Execute(w, resp)
	} else if r.Method == "GET" {
		templateLogin.Execute(w, "Please login")
	}
	//if err != nil {
	//	templateprofile.Execute(w, nil)
	//}
}

///////////////////
//uploadFile
func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("one uploadFile request")
		ok, _ := verifylogin(r)
		log.Println("ok: ", ok)
		if !ok {
			log.Println("Account expired")
			templateLogin.Execute(w, "Please login")
			log.Println("over")
			return
		}
		username := r.FormValue("username")
		file, header, err := r.FormFile("image")
		//log.Println("err: ", err)
		//log.Println("header.Filename: ", header.Filename)
		var resp protocol.RespProfile
		if err != nil {
			log.Println("fail to get the img")
			_ = client.Call("GetUserInfoByName", username, &resp)
			templateProfile.Execute(w, resp)
			return
		}
		defer file.Close()
		// check legal of file
		newName, isLegal := utils.CheckAndCreateFileName(header.Filename)
		log.Println("3333")
		//username := r.FormValue("username")
		if !isLegal {
			log.Println("img is no legal")
			_ = client.Call("GetUserInfoByName", username, &resp)
			templateProfile.Execute(w, resp)
			return
		}
		//filePath := config.ResourceImg + newName
		filePath := Pwd + "/resource/img/" + newName
		fileName := newName
		log.Println("Generate new file name: ", filePath)
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Open file failed, error is ", err)
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, file)
		if err != nil {
			log.Println("fail to copy imgFile: ", err)
			_ = client.Call("GetUserInfoByName", username, &resp)
			templateProfile.Execute(w, resp)
			return
		}

		req := protocol.ReqUploadFile{
			UserName: username,
			FileName: fileName,
		}
		_ = client.Call("UploadFile", req, &resp)

		log.Println("更新完头像后的返回数据为：", resp)

		templateProfile.Execute(w, resp)
	} else if r.Method == "GET" {
		templateLogin.Execute(w, "Please login")
	}
}

func transError(er error) error {
	if er.Error() == "" {
		return nil
	}
	return er
}
