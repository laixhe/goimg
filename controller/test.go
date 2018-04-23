package controller

import (
	"net/http"
)

// 测试上传
func Test(w http.ResponseWriter, r *http.Request) {

	html := `
<html>
    <head>
        <title>测试上传</title>
    </head>
    <body>
        <form action="/" method="post" enctype="multipart/form-data">
            <input type="file" name="userfile" /> 
            <input type="file" name="userfile" /> 
            <input type="submit" name="submit" />
        </form>
    </body>
</html>`

	w.Write([]byte(html))

}
