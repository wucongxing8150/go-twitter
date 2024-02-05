package Model

type CrawlSetting struct {
	SettingId       int64 `gorm:"column:setting_id;primary_key" json:"setting_id"`                    // 主键
	SettingPlatform int   `gorm:"column:setting_platform;default:1;NOT NULL" json:"setting_platform"` // 爬取平台 1:twitter 2:facebook
	Created         int   `gorm:"column:created" json:"created"`                                      // 创建时间
	Updated         int   `gorm:"column:updated" json:"updated"`                                      // 更新时间
}

func (m *CrawlSetting) TableName() string {
	return "kt_crawl_setting"
}
