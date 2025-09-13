package xiaohongshu

import (
    "context"
    "time"

    "github.com/go-rod/rod"
    "github.com/pkg/errors"
)

type LoginAction struct {
	page *rod.Page
}

func NewLogin(page *rod.Page) *LoginAction {
	return &LoginAction{page: page}
}

func (a *LoginAction) CheckLoginStatus(ctx context.Context) (bool, error) {
    pp := a.page.Context(ctx)

    // 避免 Must* 引发 panic；对导航和等待进行错误处理
    if err := pp.Navigate("https://www.xiaohongshu.com/explore"); err != nil {
        return false, errors.Wrap(err, "navigate explore failed")
    }

    // 给页面一点时间完成初始渲染（SPA 页面不一定会触发 load 事件）
    time.Sleep(1500 * time.Millisecond)

    // 在限定时间内等待“已登录”标识元素出现
    selector := `.main-container .user .link-wrapper .channel`
    _, err := pp.Timeout(8 * time.Second).Element(selector)
    if err != nil {
        // 未找到视为未登录；仅当是其它错误时再包装返回
        return false, nil
    }
    return true, nil
}

func (a *LoginAction) Login(ctx context.Context) error {
    pp := a.page.Context(ctx)

    // 导航到小红书首页，这会触发二维码弹窗
    if err := pp.Navigate("https://www.xiaohongshu.com/explore"); err != nil {
        return errors.Wrap(err, "navigate explore failed")
    }

    // 短暂等待首屏稳定
    time.Sleep(2 * time.Second)

    // 若已登录则直接返回
    if _, err := pp.Timeout(5 * time.Second).Element(`.main-container .user .link-wrapper .channel`); err == nil {
        return nil
    }

    // 等待扫码登录成功：在较长时间内等待“已登录”标识元素出现
    if _, err := pp.Timeout(2 * time.Minute).Element(`.main-container .user .link-wrapper .channel`); err != nil {
        return errors.Wrap(err, "等待登录超时或失败")
    }

    return nil
}
