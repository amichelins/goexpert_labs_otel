package presenters

import "encoding/json"

func ToJson(Message any) string {

    RespMessage, _ := json.Marshal(Message)

    return string(RespMessage)
}
