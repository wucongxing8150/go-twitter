package Model

type CrawlContent struct {
	ContentId    int64  `gorm:"column:content_id;primary_key" json:"content_id"`          // 主键
	HId          int    `gorm:"column:h_id;default:0;NOT NULL" json:"h_id"`               // 分类表id
	KeyWordId    int64  `gorm:"column:key_word_id;default:0;NOT NULL" json:"key_word_id"` // 关键词id
	UniqueId     string `gorm:"column:unique_id" json:"unique_id"`                        // 唯一值
	Language     int    `gorm:"column:language;default:1" json:"language"`                // 语言 1:英语 2:阿拉伯语
	ContentTitle string `gorm:"column:content_title;NOT NULL" json:"content_title"`       // 标题
	Content      string `gorm:"column:content" json:"content"`                            // 内容
	ContentImgs  string `gorm:"column:content_imgs;NOT NULL" json:"content_imgs"`         // 图片-逗号分隔
	ContentVideo string `gorm:"column:content_video;NOT NULL" json:"content_video"`       // 视频
	Created      int    `gorm:"column:created" json:"created"`                            // 创建时间
	Updated      int    `gorm:"column:updated" json:"updated"`                            // 更新时间
	Status       string `gorm:"column:status;default:0" json:"status"`                    // 1通过0未通过
}

func (m *CrawlContent) TableName() string {
	return "kt_crawl_content"
}
