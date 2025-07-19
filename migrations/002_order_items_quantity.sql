DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='order_items' AND column_name='quantity'
    ) THEN
        ALTER TABLE order_items ADD COLUMN quantity INT NOT NULL DEFAULT 1;
    END IF;
END $$; 