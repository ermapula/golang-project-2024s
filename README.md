# Games project
A game store like project. It contains games and publishers. Later when users are integrated, they can buy, wishlist them

## Game REST API

```
GET /games/:id
POST /games
PUT /games/:id
DELETE /games/:id
```

## DB Structure

```
Table publishers {
  id integer [primary key]
  name text
  headquarters text
  website text
}

Table games {
  id integer [primary key]
  title text
  genre text
  price double
  release_date timestamp
  publisher_id integer
}

Ref: games.publisher_id > publishers.id
```

## Project Team
Tolymbekov Ermek Beisenuly, 20B030635