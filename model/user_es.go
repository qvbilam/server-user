package model

type UserES struct {
	ID int64 `json:"id"`

	Code         int64 `json:"code"`
	Level        int64 `json:"level"`
	FansCount    int64 `json:"fans_count"`
	FollowCount  int64 `json:"follow_count"`
	GetLikeCount int64 `json:"get_like_count"`

	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`

	IsVisible bool `json:"isVisible"`

	//Score float64 `json:"score"`
}

func (UserES) GetIndexName() string {
	return "user"
}

func (UserES) GetMapping() string {
	userMapping := `{
    "mappings":{
        "properties":{
            "code":{
                "type":"integer"
            },
            "level":{
                "type":"integer"
            },
            "fans_count":{
                "type":"integer"
            },
            "follow_count":{
                "type":"integer"
            },
            "get_like_count":{
                "type":"integer"
            },
            "is_visible":{
                "type":"boolean"
            },
            "mobile":{
                "type":"text",
                "analyzer":"ik_max_word"
            },
            "nickname":{
                "type":"text",
                "analyzer":"ik_max_word"
            },
			"avatar":{
                "type":"text",
                "analyzer":"ik_max_word"
            },
			"gender":{
                "type":"text",
                "analyzer":"ik_max_word"
            }
        }
    }
}`

	return userMapping
}
