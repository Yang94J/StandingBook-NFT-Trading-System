package routes

import (
	"Go-StandingbookServer/config"
	"Go-StandingbookServer/dbs"
	"Go-StandingbookServer/eths"
	"Go-StandingbookServer/utils"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

var session *sessions.CookieStore

func init() {
	session = sessions.NewCookieStore([]byte("secret"))
}

//resp数据响应
func ResponseData(c echo.Context, resp *utils.Resp) {
	resp.ErrMsg = utils.RecodeText(resp.Errno)
	c.JSON(http.StatusOK, resp)
}

func PingHandler(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer ResponseData(c, &resp)
	return nil
}

func Register(c echo.Context) error {

	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	// Defer the population of the context till the end of the function
	defer ResponseData(c, &resp)

	user := &dbs.User{}

	if err := c.Bind(user); err != nil {
		fmt.Println("Binding ERROR!")
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	fmt.Println(user)
	address, err := eths.NewAcc(user.Password, config.ETHConnStr)
	if err != nil {
		resp.Errno = utils.RECODE_IPCERR
		return err
	}

	user.Address = address

	err = user.Add()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = address
	return nil
}

func Login(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer ResponseData(c, &resp)

	user := &dbs.User{}
	if err := c.Bind(user); err != nil {
		fmt.Println("Binding ERROR!")
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	ok, err := user.Query()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	if !ok {
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	fmt.Println(user)

	sess, _ := session.Get(c.Request(), "session")
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["address"] = user.Address
	sess.Values["password"] = user.Password
	sess.Save(c.Request(), c.Response())

	return nil
}

func GetSession(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer ResponseData(c, &resp)

	sess, err := session.Get(c.Request(), "session")
	if err != nil {
		fmt.Println("failed to get Session")
		resp.Errno = utils.RECODE_SESSIONERR
		return err
	}
	address := sess.Values["address"]
	if address == nil {
		fmt.Println("Failed to get Session, address is nil")
		resp.Errno = utils.RECODE_SESSIONERR
		return err
	}
	return nil
}

func Upload(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer ResponseData(c, &resp)

	content := &dbs.Content{}

	h, err := c.FormFile("fileName")
	if err != nil {
		fmt.Println("failed to formFile ", err)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	src, err := h.Open()
	defer src.Close()
	content.ContentPath = "static/standingBook/" + h.Filename
	dst, err := os.Create(content.ContentPath)
	defer dst.Close()

	if err != nil {
		fmt.Println("Failed to create file, ", err, content.ContentPath)
		resp.Errno = utils.RECODE_IOERR
		return err
	}
	cData := make([]byte, h.Size)
	n, err := src.Read(cData)
	if err != nil || h.Size != int64(n) {
		resp.Errno = utils.RECODE_IOERR
		return err
	}
	content.ContentHash = fmt.Sprintf("%x", sha256.Sum256(cData))
	dst.Write(cData)

	sess, _ := session.Get(c.Request(), "session")
	address, ok := sess.Values["address"].(string)
	if address == "" || !ok {
		resp.Errno = utils.RECODE_SESSIONERR
		return errors.New("Invalid Session")
	}
	password, ok := sess.Values["password"].(string)

	content.Address = address
	token_id := utils.GetToken()
	content.TokenId = string(token_id)
	content.AddContent()

	go eths.UploadPic(address, password, token_id)

	return nil
}
