package Model

type CrawlKeyWord struct {
	KeyWordId       int64  `gorm:"column:key_word_id;primary_key" json:"key_word_id"`                  // 主键
	HId             int    `gorm:"column:h_id;default:0;NOT NULL" json:"h_id"`                         // 分类表id
	ClassifyThemeId int64  `gorm:"column:classify_theme_id" json:"classify_theme_id"`                  // 主题id
	SettingPlatform int    `gorm:"column:setting_platform;default:1;NOT NULL" json:"setting_platform"` // 爬取平台 1:twitter 2:facebook
	KeyWordName     string `gorm:"column:key_word_name;NOT NULL" json:"key_word_name"`                 // 关键字名称
	Num             int    `gorm:"column:num;default:0" json:"num"`                                    // 已爬取数量
	CreatedAt       string `gorm:"column:created_at" json:"created_at"`                                // 最后爬取时间
	Created         int    `gorm:"column:created" json:"created"`                                      // 创建时间
	Updated         int    `gorm:"column:updated" json:"updated"`                                      // 更新时间
}

func (m *CrawlKeyWord) TableName() string {
	return "kt_crawl_key_word"
}
