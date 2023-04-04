package seller

import (
	"Market-Bot/models"
	"Market-Bot/sql"
	"fmt"
	"github.com/yanzay/tbot/v2"
	"strconv"
	"strings"
)

func ShowCreatedProduct(id_product int, client *tbot.Client, m *tbot.Message) {
	p := sql.GetProductInfo(id_product)
	fmt.Println(p.Product_image)
	_, err := client.SendPhotoFile(m.Chat.ID, p.Product_image, tbot.OptFoursquareID(strconv.Itoa(p.Id_product)),
		tbot.OptCaption("Номер продукта: "+strconv.Itoa(p.Id_product)+"\nЛогин продавца: "+p.Id_seller+"\nНазвание продукта: "+p.Product_name+"\nКатегория продукта: "+p.Product_category+"\nОписание продукта: "+p.Product_description+"\nЦена продукта: "+strconv.Itoa(p.Product_cost)))
	models.CheckError(err)
}

func deleteLeftPart(s string) string {
	infoNow := false
	for i := 0; i < len(s); i++ {
		if infoNow == true {
			if string(s[i]) != " " {
				return s[i:]
			}
		}
		if string(s[i]) == ":" {
			infoNow = true
		}
	}
	return s
}

func GetInfoAnalys(str string) ([]string, bool) {
	SL := strings.Split(str, "\n")
	if len(SL) != 4 {
		return SL, false
	}
	var newSL []string
	for _, v := range SL {
		newSL = append(newSL, deleteLeftPart(v))
	}
	word := newSL[2]
	_, err := strconv.Atoi(newSL[2])
	if err != nil {
		return newSL, false
	}
	_, err = strconv.Atoi(newSL[3])
	if err != nil {
		return newSL, false
	}

	fmt.Println(word)
	return newSL, true
}
