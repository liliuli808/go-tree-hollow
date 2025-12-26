-- Insert initial tags based on Category enum
INSERT INTO tags (name) VALUES 
    ('恋爱'),
    ('游戏'),
    ('音乐'),
    ('电影'),
    ('交友'),
    ('此刻'),
    ('表白'),
    ('吐槽')
ON CONFLICT(name) DO NOTHING;
