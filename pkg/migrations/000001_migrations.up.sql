CREATE TABLE IF NOT EXISTS publishers (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    headquarters text NOT NULL,
    website text NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    genres text[] NOT NULL,
    price double precision NOT NULL DEFAULT 0,
    release_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    publisher_id bigserial NOT NULL REFERENCES publishers (id),
    version integer NOT NULL DEFAULT 1
);

INSERT INTO publishers (name, headquarters, website) 
	VALUES
		('Electronic Arts', 'Redwood City, California, USA', 'https://www.ea.com'),
		('Ubisoft', 'Montreuil, France', 'https://www.ubisoft.com'),
		('Nintendo', 'Kyoto, Japan', 'https://www.nintendo.com'),
		('Activision Blizzard', 'Santa Monica, California, USA', 'https://www.activisionblizzard.com'),
		('FromSoftware', 'Tokyo, Japan', 'https://www.fromsoftware.jp');
		
INSERT INTO games (title, genres, release_date, price, publisher_id) 
	VALUES 
		('Battlefield V', '{"fps", "action"}', '2018-11-20', 59.99, 1),
		('Assassin''s Creed Unity', '{"action-adventure", "stealth"}', '2014-11-11', 39.99, 2),
		('The Legend of Zelda: Breath of the Wild', '{"action-adventure", "open world", "puzzle"}', '2017-03-03', 59.99, 3),
		('Call of Duty: Warzone', '{"battle royale", "fps"}', '2020-03-10', 0, 4),
		('Elden Ring', '{"action-rpg", "souls-like"}', '2022-02-25', 59.99, 5),
		('Apex Legends', '{"fps", "battle royale"}', '2019-02-04', 0, 1),
		('Far Cry 6', '{"action-adventure", "fps", "open world"}', '2021-10-07', 59.99, 2),
		('Super Mario Odyssey', '{"platform", "action-adventure"}', '2017-10-27', 49.99, 3),
		('Call of Duty: Ghosts', '{"fps"}', '2014-03-25', 59.99, 4),
		('Dark Souls III', '{"action-rpg", "souls-like"}', '2019-03-22', 59.99, 5),
		('Super Mario Sunshine', '{"platform", "action-adventure"}', '2002-07-19', 19.99, 3),
		('Far Cry 5', '{"action-adventure", "fps", "open world"}', '2018-03-27', 59.99, 2);
