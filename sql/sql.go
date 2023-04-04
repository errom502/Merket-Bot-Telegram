package sql

import (
	"Market-Bot/models"
	"context"
	"fmt"
	"github.com/yanzay/tbot/v2"
	"strconv"
)

func GetContext() (context.Context, context.CancelFunc) {
	if Db.Config().ConnectTimeout > 0 {
		return context.WithTimeout(context.Background(), Db.Config().ConnectTimeout)
	}
	return context.Background(), nil
}

func ClientShowCart(m *tbot.Message, client *tbot.Client, bot *tbot.Server) {
	var cart models.Cart
	cartSl := make([]models.Cart, 0)
	rows, err := Db.Query(context.Background(), `SELECT * FROM shopping_cart_table WHERE id_user = $1`, m.From.Username)
	for rows.Next() {
		err = rows.Scan(&cart.Id_cart, &cart.Id_user, &cart.Id_product)
		cartSl = append(cartSl, cart)
	}
	models.CheckError(err)
	if len(cartSl) == 0 {
		client.SendMessage(m.Chat.ID, "Ваша корзина пуста")
		return
	}
	var p models.Product
	productSl := make([]models.Product, 0)
	for _, v := range cartSl {
		rows, err = Db.Query(context.Background(), `SELECT * FROM product_table WHERE id_product = $1`, v.Id_product)
		for rows.Next() {
			err = rows.Scan(&p.Id_product, &p.Id_seller, &p.Product_name, &p.Product_category, &p.Product_description, &p.Product_image, &p.Product_cost, &p.Product_availability)
			productSl = append(productSl, p)
		}
		models.CheckError(err)
	}
	for i, v := range productSl {
		_, err := client.SendPhotoFile(m.Chat.ID, v.Product_image, tbot.OptFoursquareID(strconv.Itoa(v.Id_product)),
			tbot.OptCaption("Номер в корзине: "+strconv.Itoa(cartSl[i].Id_cart)+"\nНомер продукта: "+strconv.Itoa(v.Id_product)+"\nЛогин продавца: "+v.Id_seller+"\nНазвание продукта: "+v.Product_name+"\nКатегория продукта: "+v.Product_category+"\nОписание продукта: "+v.Product_description+"\nЦена продукта: "+strconv.Itoa(v.Product_cost)),
			tbot.OptInlineKeyboardMarkup(MakeButtonsCartProduct()))
		if err != nil {
			fmt.Println(err)
		}
	}
	models.CheckError(err)
	client.SendMessage(m.Chat.ID, "Это все ваши товары из корзины")

}

func ClientShowFavor(m *tbot.Message, client *tbot.Client, bot *tbot.Server) {
	var cart models.Cart
	cartSl := make([]models.Cart, 0)
	rows, err := Db.Query(context.Background(), `SELECT * FROM favour_table WHERE id_user = $1`, m.From.Username)
	for rows.Next() {
		err = rows.Scan(&cart.Id_cart, &cart.Id_user, &cart.Id_product)
		cartSl = append(cartSl, cart)
	}
	models.CheckError(err)
	if len(cartSl) == 0 {
		client.SendMessage(m.Chat.ID, "У вас нету избранных товаров")
		return
	}
	var p models.Product
	productSl := make([]models.Product, 0)
	for _, v := range cartSl {
		rows, err = Db.Query(context.Background(), `SELECT * FROM product_table WHERE id_product = $1`, v.Id_product)
		for rows.Next() {
			err = rows.Scan(&p.Id_product, &p.Id_seller, &p.Product_name, &p.Product_category, &p.Product_description, &p.Product_image, &p.Product_cost, &p.Product_availability)
			productSl = append(productSl, p)
		}
		models.CheckError(err)
	}
	for i, v := range productSl {
		_, err := client.SendPhotoFile(m.Chat.ID, v.Product_image, tbot.OptFoursquareID(strconv.Itoa(v.Id_product)),
			tbot.OptCaption("Номер в избранном: "+strconv.Itoa(cartSl[i].Id_cart)+"\nНомер продукта: "+strconv.Itoa(v.Id_product)+"\nЛогин продавца: "+v.Id_seller+"\nНазвание продукта: "+v.Product_name+"\nКатегория продукта: "+v.Product_category+"\nОписание продукта: "+v.Product_description+"\nЦена продукта: "+strconv.Itoa(v.Product_cost)),
			tbot.OptInlineKeyboardMarkup(MakeButtonsFavorProduct()))
		if err != nil {
			fmt.Println(err)
		}
	}
	models.CheckError(err)
	client.SendMessage(m.Chat.ID, "Это все ваши избранные товары")

}

func ClientOrderShow(id_user string) (bool, []models.Product) {
	productSl := make([]models.Product, 0)
	var baseExist bool
	row := Db.QueryRow(
		context.Background(),
		"select exists(SELECT FROM ordered_products_table WHERE id_user = $1)", id_user)
	err := row.Scan(&baseExist)
	models.CheckError(err)
	if baseExist == false {
		return false, productSl
	}

	var o models.Order
	orderSl := make([]models.Order, 0)
	rows, err := Db.Query(context.Background(), `SELECT * FROM ordered_products_table WHERE id_user = $1`, id_user)
	for rows.Next() {
		err = rows.Scan(&o.Id_Order, &o.Id_user, &o.Id_product)
		orderSl = append(orderSl, o)
	}
	models.CheckError(err)
	var p models.Product

	for _, v := range orderSl {
		rows, err = Db.Query(context.Background(), `SELECT * FROM product_table WHERE id_product = $1`, v.Id_product)
		for rows.Next() {
			err = rows.Scan(&p.Id_product, &p.Id_seller, &p.Product_name, &p.Product_category, &p.Product_description, &p.Product_image, &p.Product_cost, &p.Product_availability)
			productSl = append(productSl, p)
		}
		models.CheckError(err)
	}
	return true, productSl
}

func ClientBuy(id_user string) bool {
	var baseExist bool
	row := Db.QueryRow(
		context.Background(),
		"select exists(SELECT FROM shopping_cart_table WHERE id_user = $1)", id_user)
	err := row.Scan(&baseExist)
	models.CheckError(err)
	if baseExist == false {
		return false
	}
	var id_p int
	cartSl := make([]int, 0)
	rows, err := Db.Query(context.Background(), `SELECT id_product FROM shopping_cart_table WHERE id_user = $1`, id_user)
	for rows.Next() {
		err = rows.Scan(&id_p)
		cartSl = append(cartSl, id_p)
	}
	for _, v := range cartSl {
		_, err = Db.Exec(context.Background(), "insert into ordered_products_table (id_user,id_product) values ($1,$2)", id_user, v)
		models.CheckError(err)
		_, err = Db.Exec(context.Background(), "update product_table set product_availability = numDown(product_availability) where id_product = $1 ", v)
		models.CheckError(err)
	}
	return true
}

//_, err = Db.Exec(context.Background(), "update product_table set product_availability = numUp(product_availability) where id_product = $1 ", id_p)
//models.CheckError(err)
func DeleteOneProduct(from, id_p, user_id string) {
	switch from {
	case "Корзина":
		_, err := Db.Exec(context.Background(), "delete from shopping_cart_table WHERE id_cart = $1 and id_user = $2;", id_p, user_id)
		models.CheckError(err)
	case "Избранное":
		_, err := Db.Exec(context.Background(), "delete from favour_table WHERE id_favour = $1 and id_user = $2;", id_p, user_id)
		models.CheckError(err)
	case "Объявление":
		_, err := Db.Exec(context.Background(), "delete from product_table WHERE id_product = $1", id_p)
		models.CheckError(err)
	}
}
func DeleteAllProducts(from string, user_id string) bool {
	var Exist bool
	switch from {
	case "Корзина":
		row := Db.QueryRow(
			context.Background(),
			"select exists(SELECT FROM shopping_cart_table WHERE id_user = $1)", user_id)
		err := row.Scan(&Exist)
		models.CheckError(err)
		if Exist == false {
			return false
		}
		_, err = Db.Exec(context.Background(), "delete from shopping_cart_table where id_user = $1;", user_id)
		models.CheckError(err)
	case "Избранное":
		row := Db.QueryRow(
			context.Background(),
			"select exists(SELECT FROM favour_table WHERE id_user = $1)", user_id)
		err := row.Scan(&Exist)
		models.CheckError(err)
		if Exist == false {
			return false
		}
		_, err = Db.Exec(context.Background(), "delete from favour_table where id_user = $1;", user_id)
		models.CheckError(err)
	}
	return true
}
func ShowSellerAd(user string, client *tbot.Client, m *tbot.Message) {
	var p models.Product
	productSl := make([]models.Product, 0)
	rows, err := Db.Query(context.Background(), `SELECT * FROM product_table WHERE id_seller = $1`, user)
	for rows.Next() {
		err = rows.Scan(&p.Id_product, &p.Id_seller, &p.Product_name, &p.Product_category, &p.Product_description, &p.Product_image, &p.Product_cost, &p.Product_availability)
		productSl = append(productSl, p)
	}
	models.CheckError(err)
	if len(productSl) == 0 {
		client.SendMessage(m.Chat.ID, "У вас нету объявлений")
		return
	}
	for _, v := range productSl {
		_, err := client.SendPhotoFile(m.Chat.ID, v.Product_image, tbot.OptFoursquareID(strconv.Itoa(v.Id_product)),
			tbot.OptCaption("Номер продукта: "+strconv.Itoa(v.Id_product)+"\nЛогин продавца: "+v.Id_seller+"\nНазвание продукта: "+v.Product_name+"\nКатегория продукта: "+v.Product_category+"\nОписание продукта: "+v.Product_description+"\nЦена продукта: "+strconv.Itoa(v.Product_cost)+"\nВ наличии: "+strconv.Itoa(v.Product_availability)),
			tbot.OptInlineKeyboardMarkup(MakeButtonsAD()))
		if err != nil {
			fmt.Println(err)
		}
	}
}
func MoveToCartFavor(object string, id string, id_user string) {
	if object == "Корзина" {
		var id_p int
		row := Db.QueryRow(
			context.Background(),
			"SELECT id_product FROM favour_table WHERE id_favour = $1", id)
		err := row.Scan(&id_p)
		models.CheckError(err)
		_, err = Db.Exec(context.Background(), "insert into shopping_cart_table(id_user,id_product) values ($1,$2);", id_user, id_p)
		models.CheckError(err)
		_, err = Db.Exec(context.Background(), "delete from favour_table where id_favour=$1", id)
		models.CheckError(err)
		return
	}
	var id_p int
	row := Db.QueryRow(
		context.Background(),
		"SELECT id_product FROM shopping_cart_table WHERE id_cart = $1", id)
	err := row.Scan(&id_p)
	models.CheckError(err)
	_, err = Db.Exec(context.Background(), "insert into favour_table(id_user,id_product) values ($1,$2);", id_user, id_p)
	models.CheckError(err)
	_, err = Db.Exec(context.Background(), "delete from shopping_cart_table where id_cart=$1", id)
	models.CheckError(err)
	return

}

/////////////////////////////
func CategoryProductShow(category string, m *tbot.Message, client *tbot.Client) {

	var p models.Product

	productSl := make([]models.Product, 0)
	rows, err := Db.Query(context.Background(), `SELECT * FROM product_table WHERE product_category = $1 ORDER BY id_product DESC LIMIT 20;`, category)
	for rows.Next() {
		err = rows.Scan(&p.Id_product, &p.Id_seller, &p.Product_name, &p.Product_category, &p.Product_description, &p.Product_image, &p.Product_cost, &p.Product_availability)
		productSl = append(productSl, p)
	}
	models.CheckError(err)
	if len(productSl) == 0 {
		client.SendMessage(m.Chat.ID, "В данной категорри еще нет товаров")
		return
	}

	for _, v := range productSl {
		_, err := client.SendPhotoFile(m.Chat.ID, v.Product_image, tbot.OptFoursquareID(strconv.Itoa(v.Id_product)),
			tbot.OptCaption("Номер продукта: "+strconv.Itoa(v.Id_product)+"\nЛогин продавца: "+v.Id_seller+"\nНазвание продукта: "+v.Product_name+"\nКатегория продукта: "+v.Product_category+"\nОписание продукта: "+v.Product_description+"\nЦена продукта: "+strconv.Itoa(v.Product_cost)+"\nВ наличии: "+strconv.Itoa(v.Product_availability)),
			tbot.OptInlineKeyboardMarkup(MakeButtonsAddProduct()))
		if err != nil {
			fmt.Println(err)
		}
	}

}
func UpdateMessage(id string) string {
	var p models.Product
	id_p, _ := strconv.Atoi(id)
	rows, err := Db.Query(context.Background(), `SELECT * FROM product_table WHERE id_product = $1`, id_p)
	for rows.Next() {
		err = rows.Scan(&p.Id_product, &p.Id_seller, &p.Product_name, &p.Product_category, &p.Product_description, &p.Product_image, &p.Product_cost, &p.Product_availability)
	}
	models.CheckError(err)
	return "Номер продукта: " + strconv.Itoa(p.Id_product) + "\nЛогин продавца: " + p.Id_seller + "\nНазвание продукта: " + p.Product_name + "\nКатегория продукта: " + p.Product_category + "\nОписание продукта: " + p.Product_description + "\nЦена продукта: " + strconv.Itoa(p.Product_cost) + "\nВ наличии: " + strconv.Itoa(p.Product_availability)
}

func GetIdProductFromMessageCart(message string) (string, string) {
	und := "Номер в корзине: "
	id := ""
	word := ""
	for _, v := range message {
		if word != und {
			word += string(v)
		} else if string(v) == "\n" {
			break
		} else {
			id += string(v)
		}
	}

	if id == "" {
		return "", ""
	}
	var exist int
	row := Db.QueryRow(
		context.Background(),
		"SELECT id_product FROM shopping_cart_table WHERE id_cart = $1", id)
	err := row.Scan(&exist)
	models.CheckError(err)
	if exist == 0 {
		return "", ""
	}
	return id, "Корзина"
}

func GetIdProductFromAD(message string) (string, string) {
	und := "Номер продукта: "
	id := ""
	word := ""
	for _, v := range message {
		if word != und {
			word += string(v)
		} else if string(v) == "\n" {
			break
		} else {
			id += string(v)
		}
	}

	if id == "" {
		return "", ""
	}
	return id, "Объявление"
}

func GetIdProductFromMessageFavor(message string) (string, string) {
	und := "Номер в избранном: "
	id := ""
	word := ""
	for _, v := range message {
		if word != und {
			word += string(v)
		} else if string(v) == "\n" {
			break
		} else {
			id += string(v)
		}
	}
	var exist int
	row := Db.QueryRow(
		context.Background(),
		"SELECT id_product FROM favour_table WHERE id_favour = $1", id)
	err := row.Scan(&exist)
	models.CheckError(err)
	if exist == 0 {
		return "", ""
	}
	return id, "Избранное"
}

func GetIdProductFromMessage(message string) string {
	und := "Номер продукта: "
	id := ""
	word := ""
	for _, v := range message {
		if word != und {
			word += string(v)
		} else if string(v) == "\n" {
			break
		} else {
			id += string(v)
		}
	}
	var exist int
	row := Db.QueryRow(
		context.Background(),
		"SELECT product_availability FROM product_table WHERE id_product = $1", id)
	err := row.Scan(&exist)
	models.CheckError(err)
	if exist == 0 {
		return ""
	}
	return id
}

//_, err = Db.Exec(context.Background(), "update product_table set product_availability = numDown(product_availability) where id_product = $1 ", id_p)
//models.CheckError(err)
func AddToCart(id_p, id_user string) {
	_, err := Db.Exec(context.Background(), "insert into shopping_cart_table(id_user,id_product) values ($1,$2);", id_user, id_p)
	models.CheckError(err)
}

func AddToFavor(id_p, id_user string) {
	_, err := Db.Exec(context.Background(), "insert into favour_table(id_user,id_product) values ($1,$2);", id_user, id_p)
	models.CheckError(err)
}
func GetProductInfo(id_product int) models.Product {
	var p models.Product
	id_p := id_product
	rows, err := Db.Query(context.Background(), `SELECT * FROM product_table WHERE id_product = $1`, id_p)
	for rows.Next() {
		err = rows.Scan(&p.Id_product, &p.Id_seller, &p.Product_name, &p.Product_category, &p.Product_description, &p.Product_image, &p.Product_cost, &p.Product_availability)
	}
	models.CheckError(err)
	return p
}
func CreateAd(user string, category string, slInfo map[string][]string, photo string) int {
	//_, err := Db.Exec(context.Background(), "insert into product_table(id_seller,product_name,product_category,product_description,product_image,product_cost,product_availability) values ($1,$2,$3,$4,$5,$6,$7) returning id_product;", user, slInfo[user][0], category, slInfo[user][1], "./imgs/default_image.jpg", slInfo[user][2], slInfo[user][3])
	//if err != nil {
	//	panic(err)
	//}
	var id_product int
	row := Db.QueryRow(
		context.Background(),
		"insert into product_table(id_seller,product_name,product_category,product_description,product_image,product_cost,product_availability) values ($1,$2,$3,$4,$5,$6,$7) returning id_product;", user, slInfo[user][0], category, slInfo[user][1], photo, slInfo[user][2], slInfo[user][3])
	err := row.Scan(&id_product)
	models.CheckError(err)
	fmt.Println(id_product)
	return id_product
}

func MakeButtonsAddProduct() *tbot.InlineKeyboardMarkup {

	button1 := tbot.InlineKeyboardButton{
		Text:         "Добавить в корзину",
		CallbackData: "Добавить в корзину",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         "Добавить в избранное",
		CallbackData: "Добавить в избранное",
	}
	return &tbot.InlineKeyboardMarkup{
		[][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2},
		},
	}
}

func MakeButtonsCartProduct() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         "Удалить товар",
		CallbackData: "Удалить товар",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         "Переместить в избранное",
		CallbackData: "Переместить в избранное",
	}
	return &tbot.InlineKeyboardMarkup{
		[][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2},
		},
	}
}

func MakeButtonsFavorProduct() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         "Удалить товар",
		CallbackData: "Удалить товар",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         "Переместить в корзину",
		CallbackData: "Переместить в корзину",
	}
	return &tbot.InlineKeyboardMarkup{
		[][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2},
		},
	}
}
func MakeButtonsAD() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         "Удалить объявление",
		CallbackData: "Удалить объявление",
	}

	return &tbot.InlineKeyboardMarkup{
		[][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1},
		},
	}
}
