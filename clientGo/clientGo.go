package clientGo

import (
	"github.com/yanzay/tbot/v2"

	"Market-Bot/sql"

	"context"
	"fmt"
)

//данный файл только для вноса данних юзера после регистрации
// вообще пока что это все примерно и временно

func DeleteAcc(m *tbot.Message) {
	_, err := sql.Db.Exec(context.Background(), `delete from user_table where login=$1`, m.From.Username)
	if err != nil {
		panic(err)
	}
}

func ChangePassword(m *tbot.Message) {
	_, err := sql.Db.Exec(context.Background(), `UPDATE user_table SET password=$1 WHERE login=$2;`, m.Text, m.From.Username)
	if err != nil {
		panic(err)
	}
}

func LoginCheck(m *tbot.Message) (bool, string) {
	var check bool
	if err := sql.Db.QueryRow(context.Background(), "select exists(select 1 from user_table where login = $1)",
		m.From.Username).Scan(&check); err != nil {
	}
	if check {
		return false, ""
	}

	return true, ".."
}

func ClientRegistration(m *tbot.Message) {
	_, err := sql.Db.Exec(context.Background(), `INSERT INTO user_table(login,password) values ($1,$2)`, m.From.Username, m.Text)
	if err != nil {
		panic(err)
	}
}

func ClientLogin(m *tbot.Message, password string) bool {
	pass_from_db := ""
	if err := sql.Db.QueryRow(context.Background(), "SELECT (password) from user_table where login = $1",
		m.From.Username).Scan(&pass_from_db); err != nil {
		if pass_from_db == "" {
			fmt.Println("Пароль пустой")
			return false
		}
	} else {
		fmt.Println("Нет такого логина")
	}

	if password == pass_from_db {
		fmt.Println("Зашел удачно")
		return true
	} else {
		fmt.Println("Пароли не совпадают")
		return false
	}
}

func CheckCorrectPass(str string) (bool, string) {
	str += " "
	dubl := 0 // еннн5667
	thpair := str[0]
	if len(str) < 6 {
		return false, "Пароль короткий("
	}
	for i := 1; i < len(str); i++ {
		thpair = str[i-1]
		if str[i] == thpair {
			dubl++
			thpair = str[i]
			//fmt.Println(str[i])
		}
		if string(str[i-1]) == " " {
			return false, "В пароле не должно быть пробелов..."
		}
	}
	if str == "" {
		return false, "Пароль пустой!"
	}

	return true, "Хороший пароль, лайк!\nЗаношу в бд..."
}

var UserState = map[string]string{}
var UserProductInfo = map[string][]string{}
var UserProductPhoto = map[string]string{}
var UserProductCategory = map[string]string{}
