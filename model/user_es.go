package model

type UserES struct {
	ID int64 `json:"id"`

	Code        int64  `json:"code"`
	Level       int64  `json:"level"`
	FansCount   int64  `json:"fans_count"`
	FollowCount int64  `json:"follow_count"`
	Nickname    string `json:"nickname"`
	Introduce   string `json:"introduce"`
	Gender      string `json:"gender"`
	IsVisible   bool   `json:"is_visible"`
}

func (UserES) GetIndexName() string {
	return "user"
}

func (UserES) GetMapping() string {
	userMapping := `{
    "mappings":{
        "properties":{
            "level":{
                "type":"integer"
            },
            "fans_count":{
                "type":"integer"
            },
            "follow_count":{
                "type":"integer"
            },
            "is_visible":{
                "type":"boolean"
            },
            "code":{
                "type":"text"
            },
            "nickname":{
                "type":"text",
                "analyzer":"ik_max_word"
            },
            "introduce":{
                "type":"text",
                "analyzer":"ik_max_word"
            },
			"gender":{
                "type":"text"
            }
        }
    }
}`

	return userMapping
}
