package order_item

var Schema = `
	CREATE TABLE IF NOT EXISTS order_items (
    id INT GENERATED ALWAYS AS IDENTITY,
	order_id INT,
	item_name VARCHAR(255),
	quantity INT,    
	
	PRIMARY KEY (id),
	CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id)	
);
`
