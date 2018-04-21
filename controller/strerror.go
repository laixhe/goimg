package controller

func showMain() []byte {
	show := "<html><div>Goimg 轻量级的图片服务器</div>" +
		"<div><a href=\"/test\">开始吧</a></div></html>"

	return []byte(show)
}

func show404(str string) []byte {
	show := "<html><div>" + str + " not found</div>" +
		"<div><a href=\"/test\">开始吧</a></div></html>"

	return []byte(show)
}
