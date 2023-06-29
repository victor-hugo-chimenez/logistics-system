package item

var Schema = `
	CREATE TABLE IF NOT EXISTS order_status (
    id INT GENERATED ALWAYS AS IDENTITY,
	order_id INT,
	status VARCHAR(255),
	last_update_date TIMESTAMP NOT NULL DEFAULT NOW(),    
	
	CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id),    
	PRIMARY KEY (id)
);
