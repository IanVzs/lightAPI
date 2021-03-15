package chat

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/IanVzs/lightAPI/log"
	"github.com/IanVzs/lightAPI/rds"
)

type UpDetails struct {
	Name   string `json:"name"`
	Update int    `json:"update"`
}
type AlertResp struct {
	Code           string      `json:"code"`
	Status         int         `json:"status"`
	Msg            string      `json:"msg"`
	Update         int         `json:"update"`
	Update_details []UpDetails `json:"update_details"`
}

// TODO 得加入appid,否则会有问题,不过也就是多轮询一次,所以暂时不加
type GenKey struct {
	ChatId string
	Uid    string
	Action string
}

func (keys *GenKey) getKeyUid() string {
	return keys.Action + "_" + keys.ChatId + "_" + keys.Uid
}
func (keys *GenKey) getKeyNoUid() string {
	return keys.Action + "_" + keys.ChatId
}
func (keys *GenKey) getKeyNoChatID() string {
	return keys.Action + "_" + "_" + keys.Uid
}
func (keys *GenKey) getKey() string {
	actionKeyMap := map[string]bool{"get_msgs": true, "get_realtime": false, "had_msgs": true, "get_recommends": false, "get_collect_data": false}
	actionKeyMap2 := map[string]bool{"had_msgs": false}
	v, ok := actionKeyMap[keys.Action]
	if ok && v {
		v2, ok2 := actionKeyMap2[keys.Action]
		if ok2 && !v2 {
			return keys.getKeyNoChatID()
		} else {
			return keys.getKeyUid()
		}
	} else {
		return keys.getKeyNoUid()
	}
}

// http Api:获取chat提醒
func ChatAlert(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var rst AlertResp
	var chatID, uid, userType, token string
	req.ParseForm()
	chatIDs, cErr := req.Form["chat_id"]
	uids, uErr := req.Form["uid"]
	actions, aErr := req.Form["action"]
	tokens, tErr := req.Form["token"]
	var updateList []UpDetails
	if !cErr && !uErr && !tErr {
		log.Logger.Info("!cErr && !uErr")
		rst = AlertResp{"200", 1, "请求失败", 0, []UpDetails{}}
	} else if !aErr {
		log.Logger.Info("!aErr")
		rst = AlertResp{"200", 0, "请求成功", 0, []UpDetails{}}
	} else {
		if tokens != nil {
			log.Logger.Info("had token", tokens)
			token = tokens[0]
		} else {
			token = ""
		}

		listAction := strings.Split(actions[0], ",")
		log.Logger.Info("get token: ", token)
		if token != "" {
			deToken := Decrypt2String(*log.KeyAPI, token)
			log.Logger.Info("get deToken", deToken)
			listToken := strings.Split(deToken, *log.KeySplit)
			// listToken := strings.Split("hahaha_._heihei_._wwawa", "_._")
			if len(listToken) < 2 {
				chatID, uid, userType = "", "", ""
				log.Logger.Warn("decrpt failed token get:", listToken)
			} else {
				chatID, uid, userType = listToken[0], listToken[1], listToken[2]
				log.Logger.Info("had token get:", chatID, uid, userType)
			}
		} else {
			if chatIDs != nil {
				chatID = chatIDs[0]
			} else {
				chatID = ""
			}
			if uids != nil {
				uid = uids[0]
			} else {
				uid = ""
			}
			log.Logger.Info("no token get:", chatID, uid)
		}

		for _, action := range listAction {
			log.Logger.Info("chat.chatAlert get action: ", action)

			genkey := GenKey{ChatId: chatID, Uid: uid, Action: action}
			key := genkey.getKey()
			update := rds.GetDelByKey(key)
			updateList = append(updateList, UpDetails{Name: action, Update: update})
			rst = AlertResp{"200", 0, "请求成功", update, updateList}
		}
	}

	log.Logger.Info("Get rst from rds: ", rst)
	if err := json.NewEncoder(w).Encode(rst); err != nil {
		log.Logger.Fatal(err)
	}
}
