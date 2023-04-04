package sql

import (
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"context"
	"fmt"
	"os"
)

var Db *pgx.Conn

func ConnectToDB() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var dbURL string = os.Getenv("URL")

	Db, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	//	defer db.Close(context.Background())
	err = Db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	println("Successfully connected to database!")
}

func CreateDataBase() {
	var baseExist bool
	row := Db.QueryRow(
		context.Background(),
		"SELECT EXISTS (SELECT FROM pg_database WHERE datname = 'market_bot')")

	err := row.Scan(&baseExist)
	if err != nil {
		panic(err)
	}
	//fmt.Println(baseExist)
	if baseExist {
		var dbURL string = os.Getenv("URLMARKET")
		Db, err = pgx.Connect(context.Background(), dbURL)
		if err != nil {
			panic(err)
		}
		fmt.Println("database exist!\nSuccessfully connected to MARKET database!")
		return
	}
	_, err = Db.Exec(context.Background(), "create database market_bot;")
	if err != nil {
		panic(err)
	}

	var dbURL string = os.Getenv("URLMARKET")
	Db, err = pgx.Connect(context.Background(), dbURL)

	//остальное:
	_, err = Db.Exec(context.Background(), "create table user_table(login varchar(100) PRIMARY KEY NOT NULL ,password varchar(100) NOT NULL,role varchar(100) DEFAULT 'customer');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "create table category(id_category serial,name_category text NOT NULL PRIMARY KEY,number_of_product_ads int default 0);")
	if err != nil {
		panic(err)
	}

	//для продавца:
	_, err = Db.Exec(context.Background(), "create table product_table(id_product serial PRIMARY KEY,id_seller varchar(100) REFERENCES user_table(login) NOT NULL,product_name text NOT NULL,product_category text REFERENCES category(name_category) NOT NULL,product_description text not null default '-'::text,product_image text not null default '/imgs/default_image.jpg'::text,product_cost int not null default 0,product_availability int NOT NULL);")
	if err != nil {
		panic(err)
	}

	//для покупателя:
	_, err = Db.Exec(context.Background(), "create table shopping_cart_table(id_cart serial PRIMARY KEY,id_user varchar(100) REFERENCES user_table(login) ON DELETE CASCADE,id_product int REFERENCES product_table(id_product) ON DELETE CASCADE);")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "create table favour_table(id_favour serial PRIMARY KEY,id_user varchar(100) REFERENCES user_table(login) ON DELETE CASCADE,id_product int REFERENCES product_table(id_product) ON DELETE CASCADE);")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "create table ordered_products_table(id_order serial PRIMARY KEY,id_user varchar(100) REFERENCES user_table(login) ON DELETE CASCADE,id_product int REFERENCES product_table(id_product) ON DELETE CASCADE);")
	if err != nil {
		panic(err)
	}

	//заполнение категорий товара:
	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Одежда и обувь');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Аксессуары к одежде');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Электроника');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Детские товары');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Товары для хобби');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Товары для дома и сада');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Бытовая химия');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Косметика');")
	if err != nil {
		panic(err)
	}

	//_, err = Db.Exec(context.Background(), "insert into category(name_category) values ('Остальные категории');")
	//if err != nil {
	//	panic(err)
	//}

	//триггеры и функции для таблицы category
	_, err = Db.Exec(context.Background(), `
		create or REPLACE function categoryUp()
		returns trigger AS
		$$
		begin
			update category
			set number_of_product_ads = number_of_product_ads + 1
			where name_category = new.product_category;
			return new;
		end;
		$$
		language 'plpgsql';
	`)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), `
	create or REPLACE function categoryDown()
    returns trigger AS
    $$
		begin
			update category
			set number_of_product_ads = number_of_product_ads - 1
			where name_category = new.product_category;
			return new;
		end;
		$$
		language 'plpgsql';	
	`)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), `
	CREATE OR REPLACE FUNCTION numDown(var integer)
	RETURNS integer
	LANGUAGE plpgsql
	STRICT
	AS $function$
	declare
		output int := 0;
	begin
		output := var - 1;
		RETURN (output);
		END;
	$function$;
	`)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), `
	CREATE OR REPLACE FUNCTION numUp(var integer)
	RETURNS integer
	LANGUAGE plpgsql
	STRICT
	AS $function$
	declare
		output int := 0;
	begin
		output := var + 1;
		RETURN (output);
		END;
	$function$;
	`)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), `
	create trigger category_after_insert
    	after insert on product_table
    	for each row
    	execute  procedure categoryUp();
	`)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), `
	create trigger category_after_delete
    	after delete on product_table
    	for each row
    	execute  procedure categoryDown();
	`)
	if err != nil {
		panic(err)
	}

	//несколько тестовых объявлений:
	_, err = Db.Exec(context.Background(), "insert into user_table(login,password) values ('1','123456');")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into product_table(id_seller,product_name,product_category,product_description,product_image,product_cost,product_availability) values ('1','testProduct1','Одежда и обувь','some description','C:/Users/Egore/GolandProjects/Market-Bot/imgs/default_image.jpg',19000,40);")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into product_table(id_seller,product_name,product_category,product_description,product_image,product_cost,product_availability) values ('1','testProduct2','Аксессуары к одежде','some description','C:/Users/Egore/GolandProjects/Market-Bot/imgs/image_2022-06-04_22-44-59.png',19000,40);")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec(context.Background(), "insert into product_table(id_seller,product_name,product_category,product_description,product_image,product_cost,product_availability) values ('1','testProduct462','Аксессуары к одежде','some description','C:/Users/Egore/GolandProjects/Market-Bot/imgs/IMG_0864.PNG',19000,40);")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created tables")
}
