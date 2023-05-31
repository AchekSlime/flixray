package service

import (
	"bytes"
	crypto "crypto/rand"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"math/big"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"

	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/flixray/auth/service/dto"
)

func (auth *AuthService) RegMail(context *gin.Context) {
	// get body from request.
	var regMailRequest dto.RegMailRequest
	if err := context.ShouldBindJSON(&regMailRequest); err != nil {
		rest_api_utils.BindBadRequest(context, err)
		return
	}

	// Положить нового пользователя в мапу.
	generatedCode := NewCryptoRand()
	auth.unconfirmedMail[generatedCode] = regMailRequest.Email

	//отправить письмо
	text := fmt.Sprintf("%d", generatedCode)
	ok, err := sendEmail(text, regMailRequest.Email, "flixray.company@gmail.com", "cqqrbbqmcgcrigbz",
		"smtp.gmail.com", "25")
	if err != nil || !ok {
		rest_api_utils.BindInternalError(context, fmt.Errorf("error while sending confirm email, err: %s", err))
		fmt.Printf("error while sending email: %s", err)
		return
	} else {
		logrus.Info("email sanded")
	}

	//response := dto.RegMailResponse{
	//	Code: generatedCode,
	//}
	context.JSON(http.StatusOK, nil)
}

func NewCryptoRand() int64 {
	safeNum, err := crypto.Int(crypto.Reader, big.NewInt(100234))
	if err != nil {
		panic(err)
	}
	return safeNum.Int64()
}

type tmplt struct {
	Code string
}

func sendEmail(text, to, from, password, addr, port string) (bool, error) {
	auth := smtp.PlainAuth("", from, password, addr)

	templateData := tmplt{
		Code: text,
	}

	body, err := parseTemplate("confirm.html", templateData)
	if err != nil {
		fmt.Printf("error generating html: %s", err)
	}

	// Вынести
	subject := "Subject: Confirm your email address\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + body)

	if err := smtp.SendMail(addr+":"+port, auth, from, []string{to}, msg); err != nil {
		return false, err
	}

	return true, nil
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
