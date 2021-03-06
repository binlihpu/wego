package media

import (
	"encoding/json"
	"net/url"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/media/menu"
	"github.com/zozowind/wego/util"
)

const (
	menuCreateURL            = core.WxAPIURL + "/cgi-bin/menu/create"
	menuGetURL               = core.WxAPIURL + "/cgi-bin/menu/get"
	menuDeleteURL            = core.WxAPIURL + "/cgi-bin/menu/delete"
	menuAddConditionalURL    = core.WxAPIURL + "/cgi-bin/menu/addconditional"
	menuDeleteConditionalURL = core.WxAPIURL + "/cgi-bin/menu/delconditional"
	menuTryMatchURL          = core.WxAPIURL + "/cgi-bin/menu/trymatch"
	menuSelfMenuInfoURL      = core.WxAPIURL + "/cgi-bin/get_current_selfmenu_info"
)

//SetMenu 设置按钮
func (wm *WeMediaClient) SetMenu(buttons []*menu.Button) (err error) {
	token, err := wm.TokenServer.Token()
	if nil != err {
		return
	}
	params := url.Values{}
	params.Set("access_token", token)

	req := &menu.SetReq{
		Button: buttons,
	}

	data, err := util.HTTPJsonPost(nil, menuCreateURL+"?"+params.Encode(), req)
	if nil != err {
		return
	}
	rsp := &core.WxErrorResponse{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	return
}

// GetMenu 获取菜单配置
func (wm *WeMediaClient) GetMenu() (menus *menu.MenusRes, err error) {
	token, err := wm.TokenServer.Token()
	if nil != err {
		return
	}
	params := url.Values{}
	params.Set("access_token", token)

	data, err := util.HTTPGet(nil, menuGetURL+"?"+params.Encode())
	rsp := &core.WxErrorResponse{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	if nil != err {
		return
	}
	menus = &menu.MenusRes{}
	err = json.Unmarshal(data, menus)
	return
}

//DeleteMenu 删除菜单
func (wm *WeMediaClient) DeleteMenu() (err error) {
	token, err := wm.TokenServer.Token()
	if nil != err {
		return
	}
	params := url.Values{}
	params.Set("access_token", token)

	data, err := util.HTTPGet(nil, menuDeleteURL+"?"+params.Encode())
	rsp := &core.WxErrorResponse{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	return
}

//AddConditional 添加个性化菜单
func (wm *WeMediaClient) AddConditional(buttons []*menu.Button, rule *menu.MatchRule) (err error) {
	token, err := wm.TokenServer.Token()
	if nil != err {
		return
	}
	params := url.Values{}
	params.Set("access_token", token)

	req := &menu.SetReq{
		Button:    buttons,
		MatchRule: rule,
	}

	data, err := util.HTTPJsonPost(nil, menuAddConditionalURL+"?"+params.Encode(), req)
	if nil != err {
		return
	}
	rsp := &core.WxErrorResponse{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	return
}

//DeleteConditional 删除个性化菜单
func (wm *WeMediaClient) DeleteConditional(menuID int64) (err error) {
	token, err := wm.TokenServer.Token()
	if nil != err {
		return
	}
	params := url.Values{}
	params.Set("access_token", token)

	req := &menu.DeleteConditionalReq{
		MenuID: menuID,
	}

	data, err := util.HTTPJsonPost(nil, menuDeleteConditionalURL+"?"+params.Encode(), req)
	if nil != err {
		return
	}
	rsp := &core.WxErrorResponse{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	return
}

// //MenuTryMatch 菜单匹配
// func (wm *WeMediaClient) MenuTryMatch(userID string) (buttons []Button, err error) {

// }

// //GetCurrentSelfMenuInfo 获取自定义菜单配置接口
// func (wm *WeMediaClient) GetCurrentSelfMenuInfo() (resSelfMenuInfo ResSelfMenuInfo, err error) {

// }
