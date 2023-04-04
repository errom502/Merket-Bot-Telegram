package customer

import (
	"Market-Bot/clientGo"
	"github.com/yanzay/tbot/v2"

	"Market-Bot/models"
	"Market-Bot/sql"
	"context"
	"fmt"
	"strconv"
)

//{(strconv.Itoa(c.Number)+". "+c.Name+"("+strconv.Itoa(c.Ads_score)+")"),false}
func GetCategory(client string) []models.Chosen_Category {
	var c models.Category
	categorySl := make([]models.Chosen_Category, 0)
	rows, err := sql.Db.Query(context.Background(), `select * from category`)
	for rows.Next() {
		err = rows.Scan(&c.Number, &c.Name, &c.Ads_score)
		if clientGo.UserState[client] == "CLIENT_INTERFACE" {
			categorySl = append(categorySl, models.Chosen_Category{
				Category_info: (c.Name + "(" + strconv.Itoa(c.Ads_score) + ")"),
				Chosen:        c.Name,
			})
		} else if clientGo.UserState[client] == "SELLER_AD_CREATE" {
			categorySl = append(categorySl, models.Chosen_Category{
				Category_info: (c.Name),
				Chosen:        c.Name,
			})
		}
	}
	if err != nil {
		panic(err)
	}
	return categorySl
}

func ClientShowOrderProduct(m *tbot.Message, client *tbot.Client, bot *tbot.Server) {
	check, productSl := sql.ClientOrderShow(m.Chat.Username)
	if check {
		for _, v := range productSl {
			_, err := client.SendPhotoFile(m.Chat.ID, v.Product_image, tbot.OptFoursquareID(strconv.Itoa(v.Id_product)),
				tbot.OptCaption("Номер продукта: "+strconv.Itoa(v.Id_product)+"\nЛогин продавца: "+v.Id_seller+"\nНазвание продукта: "+v.Product_name+"\nКатегория продукта: "+v.Product_category+"\nОписание продукта: "+v.Product_description+"\nЦена продукта: "+strconv.Itoa(v.Product_cost)))
			models.CheckError(err)
		}
		return
	}
	client.SendMessage(m.Chat.ID, "Вы еще ничего не покупали")
}

func ClientShowCategory(m *tbot.Message, client *tbot.Client) {
	categorySl := GetCategory(m.From.Username)
	fmt.Println(categorySl)

	_, err := client.SendMessage(m.Chat.ID, "Категории товаров:", tbot.OptInlineKeyboardMarkup(makeCategoryButtons(categorySl)))
	if err != nil {
		panic(err)
	}
	//ChosenID := fmt.Sprintf("%s:%d", m.Chat.ID, msg.MessageID)

}

func CallBackDataOn(UserState map[string]string, client *tbot.Client, bot *tbot.Server) {
	app := &Application{client: client}
	bot.HandleCallback(app.CallbackCategoryHandler)
}

type Application struct {
	client    *tbot.Client
	UserState map[string]string
}

func ClientBuyAllProduct(m *tbot.Message, client *tbot.Client, bot *tbot.Server) {
	check := sql.ClientBuy(m.From.Username)
	if check {
		sql.DeleteAllProducts("Корзина", m.From.Username)
		client.SendSticker(m.Chat.ID, "CAACAgIAAxkBAAEGqr1ji5Nax0hzev6QZMy6X0-1AlPJ0gACrRgAAnmFiUi3haSlMSLa5SsE")
		client.SendMessage(m.Chat.ID, "Товары куплены")
		return
	}
	client.SendMessage(m.Chat.ID, "Ваша корзина пуста")
	return
}

func ClientDeleteAllProductsFromCart(m *tbot.Message, client *tbot.Client, bot *tbot.Server) {

	check := sql.DeleteAllProducts("Корзина", m.From.Username)
	if check {
		client.SendAnimationFile(m.Chat.ID, "./imgs/ryan-gosling.gif", tbot.OptCaption("Товары удалены из корзины"))
		return
	}
	client.SendAnimationFile(m.Chat.ID, "./imgs/tne.gif", tbot.OptCaption("У вас товаров то нет"))
}
func ClientDeleteAllProductsFromFavor(m *tbot.Message, client *tbot.Client, bot *tbot.Server) {
	check := sql.DeleteAllProducts("Избранное", m.From.Username)
	if check {
		client.SendAnimationFile(m.Chat.ID, "./imgs/26be37b8-7610-4af0-b50b-eec01e51275e.gif", tbot.OptCaption("Товары удалены из избранного"))
		return
	}
	client.SendAnimationFile(m.Chat.ID, "./imgs/tne.gif", tbot.OptCaption("У вас товаров то нет"))
}

//var p models.ProductId
func (a *Application) CallbackCategoryHandler(cq *tbot.CallbackQuery) {
	switch cq.Data {
	case "Добавить в избранное":
		id := sql.GetIdProductFromMessage(cq.Message.Caption)
		if id == "" {
			a.client.SendMessage(cq.Message.Chat.ID, "Данного товара нет в наличии.")
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("out of range"))
			return
		}
		fmt.Println("\nid", id)
		sql.AddToFavor(id, cq.From.Username)
		//text := sql.UpdateMessage(id)
		//a.client.EditMessageCaption(cq.Message.Chat.ID, cq.Message.MessageID, text, tbot.OptInlineKeyboardMarkup(sql.MakeButtonsAddProduct()))
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
		return
	case "Добавить в корзину":
		id := sql.GetIdProductFromMessage(cq.Message.Caption)
		if id == "" {
			a.client.SendMessage(cq.Message.Chat.ID, "Данного товара нет в наличии.")
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("out of range"))
			return
		}
		fmt.Println("\nid", id)
		sql.AddToCart(id, cq.From.Username)
		//text := sql.UpdateMessage(id)
		//a.client.EditMessageCaption(cq.Message.Chat.ID, cq.Message.MessageID, text, tbot.OptInlineKeyboardMarkup(sql.MakeButtonsAddProduct()))
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
		return
	case "Удалить товар":
		id, from := sql.GetIdProductFromMessageCart(cq.Message.Caption)
		if id == "" && from == "" {
			id, from = sql.GetIdProductFromMessageFavor(cq.Message.Caption)

		}
		if id == "" {
			a.client.SendMessage(cq.Message.Chat.ID, "Данного товара нет в наличии.")
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("out of range"))
			return
		}
		sql.DeleteOneProduct(from, id, cq.From.Username)
		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
	case "Удалить объявление":
		id, from := sql.GetIdProductFromAD(cq.Message.Caption)
		if id == "" && from == "" {
			id, from = sql.GetIdProductFromMessageFavor(cq.Message.Caption)

		}
		if id == "" {
			a.client.SendMessage(cq.Message.Chat.ID, "Данного товара нет в наличии.")
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("out of range"))
			return
		}
		sql.DeleteOneProduct(from, id, cq.From.Username)
		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
	case "Переместить в корзину":
		id, _ := sql.GetIdProductFromMessageFavor(cq.Message.Caption)
		if id == "" {
			a.client.SendMessage(cq.Message.Chat.ID, "Данного товара нет в наличии.")
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("out of range"))
			return
		}
		sql.MoveToCartFavor("Корзина", id, cq.From.Username)
		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
	case "Переместить в избранное":
		id, _ := sql.GetIdProductFromMessageCart(cq.Message.Caption)
		if id == "" {
			a.client.SendMessage(cq.Message.Chat.ID, "Данного товара нет в наличии.")
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("out of range"))
			return
		}
		sql.MoveToCartFavor("Избранное", id, cq.From.Username)
		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))

	default:
		if clientGo.UserState[cq.From.Username] == "CLIENT_INTERFACE" {
			sql.CategoryProductShow(cq.Data, cq.Message, a.client)
			a.client.EditMessageReplyMarkup(cq.Message.Chat.ID, cq.Message.MessageID, tbot.OptInlineKeyboardMarkup(deleteCategoryButtons()))
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
		}
		if clientGo.UserState[cq.From.Username] == "SELLER_AD_CREATE" {
			a.client.SendMessage(cq.Message.Chat.ID, "Пришлите инфу как по след примеру:", tbot.OptReplyKeyboardRemove)
			a.client.SendMessage(cq.Message.Chat.ID, "\nНазвание продукта: Название вашего продукта\nОписание продукта: ваше описание продукта\nЦена продукта: ваша стоимость продукта за 1 шт\nВ наличии: количество продукта доступное для покупки")
			a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
			a.client.SendMessage(cq.Message.Chat.ID, "Скопируйте текст следующего сообщения для удобства заполнения", tbot.OptReplyKeyboardRemove)
			a.client.SendMessage(cq.Message.Chat.ID, "\nНазвание продукта: \nОписание продукта: \nЦена продукта: \nВ наличии: ")

			clientGo.UserState[cq.From.Username] = "SELLER_AD_GET_INFO"
			clientGo.UserProductCategory[cq.From.Username] = cq.Data
			// функция для создарния
		}
	}
}

func makeCategoryButtons(variable []models.Chosen_Category) *tbot.InlineKeyboardMarkup {
	buttons := make([]tbot.InlineKeyboardButton, len(variable))
	for i, v := range variable {
		buttons[i].Text = v.Category_info
		buttons[i].CallbackData = v.Chosen
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{buttons[0], buttons[1]},
			[]tbot.InlineKeyboardButton{buttons[2], buttons[3]},
			[]tbot.InlineKeyboardButton{buttons[4], buttons[5]},
			[]tbot.InlineKeyboardButton{buttons[6], buttons[7]},
		},
	}
}

func deleteCategoryButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("OK"),
		CallbackData: "",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1},
		},
	}
}
