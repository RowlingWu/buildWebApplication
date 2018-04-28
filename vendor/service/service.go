package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"manager"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// InitRoutes 返回一个经过初始化的 *mux.Router
func InitRoutes() *mux.Router {
	mx := mux.NewRouter()
	mx.HandleFunc("/{id}", deleteHandler).Methods("DELETE")
	mx.HandleFunc("/", postHandler).Methods("POST")
	return mx
}

// glManager 是全局的任务管理器，负责定时任务的创建、删除
var glManager = manager.MyManager

// postHandler 处理创建定时任务的handler
func postHandler(w http.ResponseWriter, r *http.Request) {
	sbody, _ := ioutil.ReadAll(r.Body) // 读取http的json参数
	body, _ := url.QueryUnescape(string(sbody))
	defer r.Body.Close()

	var task manager.Task
	err := json.Unmarshal([]byte(body), &task) // 从json中解析Task的内容

	if err != nil { // json解析参数出错
		fmt.Println("json error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
	} else { // json解析无误，判断请求是POST还是DELETE
		ok, err := glManager.Create(&task) // 创建定时任务

		switch {
		case ok: // 任务创建成功，返回200
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"ok\":true, \"id\":\"" + task.ID + "\"}\n"))
		case !ok: // 已存在该任务，任务创建失败，返回409
			w.WriteHeader(409)
			w.Write([]byte("{\"ok\":false, \"error\":\"The task " + task.ID + " already exists.\"}\n"))
		case err != nil: // 创建任务时发生其他错误
			fmt.Println("manager.create() error:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error() + "\n"))
		}
	}
}

// deleteHandler 处理删除定时任务的handler
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // vars 获取url中的taskID

	if err := glManager.Destroy(vars["id"]); err != nil { // 没有该任务，删除该任务时出错，返回404
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"ok\":false, \"error\":\"The task " + vars["id"] + " is not found.\"}\n"))
	} else { // 删除该任务无误，返回200
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"ok\":true, \"id\":\"" + vars["id"] + "\"}\n"))
	}
}
