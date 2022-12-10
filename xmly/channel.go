package xmly

type Channel struct {
	ChannelName    string       `json:"channelName,omitempty"`
	SubChannels    []SubChannel `json:"subChannels,omitempty"`
	VipRate        int          `json:"vipRate,omitempty"`
	EndState       int          `json:"endState,omitempty"`
	SubChannelSize int          `json:"subChannelSize,omitempty"`
	AlbumCount     int64        `json:"albumCount,omitempty"`
}

type SubChannel struct {
	ChannelName string   `json:"channelName,omitempty"`
	VipRate     int      `json:"vipRate,omitempty"`
	EndState    int      `json:"endState,omitempty"`
	AlbumCount  int64    `json:"albumCount,omitempty"`
	ShowTop3    []*Album `json:"showTop3,omitempty"`
	AudioTop3   []*Item  `json:"audioTop3,omitempty"`
}

func GetInitialChannels() map[int]*Channel {
	channels := make(map[int]*Channel, 26)

	//There is no better way for the time being,
	//it is too heavy to introduce a third party to directly analyze the webpage
	channels[7] = &Channel{ChannelName: "小说"}
	channels[11] = &Channel{ChannelName: "儿童"}
	channels[9] = &Channel{ChannelName: "相声小品"}
	channels[10] = &Channel{ChannelName: "评书"}
	channels[13] = &Channel{ChannelName: "娱乐"}
	channels[14] = &Channel{ChannelName: "悬疑"}
	channels[17] = &Channel{ChannelName: "人文"}
	channels[18] = &Channel{ChannelName: "国学"}
	channels[24] = &Channel{ChannelName: "头条"}
	channels[19] = &Channel{ChannelName: "音乐"}
	channels[16] = &Channel{ChannelName: "历史"}
	channels[20] = &Channel{ChannelName: "情感"}
	channels[26] = &Channel{ChannelName: "投资理财"}
	channels[31] = &Channel{ChannelName: "个人提升"}
	channels[22] = &Channel{ChannelName: "健康"}
	channels[21] = &Channel{ChannelName: "生活"}
	channels[15] = &Channel{ChannelName: "影视"}
	channels[27] = &Channel{ChannelName: "商业管理"}
	channels[29] = &Channel{ChannelName: "英语"}
	channels[12] = &Channel{ChannelName: "少儿素养"}
	channels[28] = &Channel{ChannelName: "科技"}
	channels[32] = &Channel{ChannelName: "教育考试"}
	channels[25] = &Channel{ChannelName: "体育"}
	channels[30] = &Channel{ChannelName: "小语种"}
	channels[8] = &Channel{ChannelName: "广播剧"}
	channels[23] = &Channel{ChannelName: "汽车"}

	return channels
}
