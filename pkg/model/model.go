package model

import "time"

type Vehicle struct {
	VRM     string `json:"vrm"`
	Country string `json:"country"`
	Make    string `json:"make"`
}

type Driver struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	PostCode  string `json:"post_code"`
	City      string `json:"city"`
	Region    string `json:"region"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type Data struct {
	Id               string    `json:"id"`
	LocationDateTime time.Time `json:"location_datetime"`
	Location         string    `json:"location"`
	TotalAmount      float64   `json:"total_amount"`
	Currency         string    `json:"currency"`
	Vehicle          Vehicle   `json:"vehicle"`
	Driver           Driver    `json:"driver"`
}

type Response struct {
	Status string `json:"status"`
}
