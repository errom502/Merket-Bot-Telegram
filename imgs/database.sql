drop table user_table;

create or replace function downgrade_value(value int)
	RETURNS integer
	LANGUAGE plpgsql
	STRICT
AS $function$
declare
output int := 0;
begin
output := value - 1;
RETURN (output);
END;
$function$;


--остальное:
create table user_table(login varchar(100), password varchar(100), role varchar(100) DEFAULT "customer");
create table category(id_category int PRAMARY KEY, name_category text);

--для продавца:
create table product_table(id_product int PRIMARY KEY, id_seller varchar(100),  product_name text,product_category text, product_description text not null default '-'::text,
                           product_image bytea not null default "C:\Users\Egore\GolandProjects\Market-Bot\default_image.jpg",
                           product_cost float4 not null default 0.0, product_availability int,FOREIGN KEY (id_seller) REFERENCES user_table(login),
                           FOREIGN KEY (product_category) REFERENCES category(name_category));


--для покупателя:
create table shopping_cart_table(id_user varchar(100),id_product int,FOREIGN KEY (id_user) REFERENCES user_table(login),FOREIGN KEY (id_product) REFERENCES product_table(id_product));




--insert into scientist (id, firstname, lastname) values (1, 'albert', 'einstein');
--insert into scientist (id, firstname, lastname) values (2, 'isaac', 'newton');
--insert into scientist (id, firstname, lastname) values (3, 'marie', 'curie');
--select * from scientist;