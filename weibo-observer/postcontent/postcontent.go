package postcontent

import "time"

type ContentType int

const (
    Invalid ContentType = iota
    Sad
    Happy
    Injured
    Ordinary
    Wrong
)

// 定义微博类，分析属性，post content 类来描述我们微博也描述我们的评论

type PostContent struct {
    Id       int    // weibo id
    Content  string // weibo content
    PostMan  string // poster of post
    PostTime time.Time
    Type     ContentType
    To       string // comment to
}
