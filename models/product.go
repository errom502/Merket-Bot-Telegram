package models

type Product struct {
	Id_product           int
	Id_seller            string
	Product_name         string
	Product_category     string
	Product_description  string
	Product_image        string
	Product_cost         int
	Product_availability int
}

type ProductId struct {
	Id_product string
	Id_user    string
}
