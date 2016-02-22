package main

import (
	"log"
)

type KarmaPlugin struct {
	Name string
}

type Karma struct {
	Id     int    `db:"id, primarykey, autoincrement"`
	User   string `db:"name, size:500"`
	Points int    `db:"points"`
}

func (kp KarmaPlugin) Register() (err error) {
	kp.Name = "Karma"
	return nil
}

func (kp KarmaPlugin) Parse(sender, channel, input string, conn *Connection) (err error) {
	if !Match(input, "[0-9a-zA-Z._-]+(\\+|-){2,}") {
		return nil
	}

	change := 0
	user := MatchAndPull(input, "[0-9a-zA-Z._-]+\\+\\+", "([0-9a-zA-Z._-]+)\\+\\+")
	if user != "" {
		change = 1
		conn.SendChan(user + " has one more karma I guess! ¯\\_(ツ)_/¯")
	} else {
		user := MatchAndPull(input, "[0-9a-zA-Z._-]+--", "([0-9a-zA-Z._-]+)--")
		if user != "" {
			change = -1
			conn.SendChan(user + " has one less karma I guess! ¯\\_(ツ)_/¯")
		}
	}
	if change != 0 {
		k := Karma{User: user}
		err = k.FindOrCreate()
		if err != nil {
			log.Println("ERROR: Unable to find or create karma entry:", err)
			return
		}
		k.Points = k.Points + change
		err = k.Update()
		if err != nil {
			log.Println("ERROR: Unable to update karma entry:", err)
			return
		}
	}
	return nil
}

func (kp KarmaPlugin) Help() (helpText string) {
	return "<name>++ or <name>--"
}

func (k *Karma) FindOrCreate() (err error) {
	err = database.Dbm.SelectOne(&k, "select * from karma where user=?", k.User)
	return err
}

func (k *Karma) Update() (err error) {
	return nil
}
