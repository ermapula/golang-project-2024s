package model

import (
	"errors"
)

type Game struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	ReleaseDate string `json:"releaseDate"`
	Price       string `json:"price"`
	PublisherId string `json:"publisherId"`
}

var games = []Game{
	{
		Id:          "1",
		Title:       "Battlefield V",
		Genre:       "First Person Shooter",
		ReleaseDate: "2018-11-20",
		Price:       "$59.99",
		PublisherId: "1",
	},
	{
		Id:          "2",
		Title:       "Assassin's Creed Unity",
		Genre:       "Action-Adventure",
		ReleaseDate: "2014-11-11",
		Price:       "$39.99",
		PublisherId: "2",
	},
	{
		Id:          "3",
		Title:       "The Legend of Zelda: Breath of the Wild",
		Genre:       "Action-Adventure",
		ReleaseDate: "2017-03-03",
		Price:       "$59.99",
		PublisherId: "3",
	},
	{
		Id:          "4",
		Title:       "Call of Duty: Warzone",
		Genre:       "Battle Royale",
		ReleaseDate: "2020-03-10",
		Price:       "0",
		PublisherId: "4",
	},
	{
		Id:          "5",
		Title:       "Elden Ring",
		Genre:       "Action RPG",
		ReleaseDate: "2022-02-25",
		Price:       "$59.99",
		PublisherId: "5",
	},
	{
		Id:          "6",
		Title:       "Apex Legends",
		Genre:       "Battle Royale",
		ReleaseDate: "2019-02-04",
		Price:       "0",
		PublisherId: "1",
	},
	{
		Id:          "7",
		Title:       "Far Cry 6",
		Genre:       "First Person Shooter",
		ReleaseDate: "2021-10-07",
		Price:       "$59.99",
		PublisherId: "2",
	},
	{
		Id:          "8",
		Title:       "Super Mario Odyssey",
		Genre:       "Platformer",
		ReleaseDate: "2017-10-27",
		Price:       "$49.99",
		PublisherId: "3",
	},
	{
		Id:          "9",
		Title:       "Call of Duty: Ghosts",
		Genre:       "First Person Shooter",
		ReleaseDate: "2014-03-25",
		Price:       "$59.99",
		PublisherId: "4",
	},
	{
		Id:          "10",
		Title:       "Dark Souls III",
		Genre:       "Action RPG",
		ReleaseDate: "2019-03-22",
		Price:       "$59.99",
		PublisherId: "5",
	},
}

func GetGames() []Game {
	return games
}

func GetGame(id string) (*Game, error) {
	for _, g := range games {
		if g.Id == id {
			return &g, nil
		}
	}
	return nil, errors.New("Game not found")
}
func GetPublisherGames(pId string) []Game {
	publisherGames := make([]Game, 0)
	for _, g := range games {
		if g.PublisherId == pId {
			publisherGames = append(publisherGames, g)
		}
	}
	return publisherGames
}