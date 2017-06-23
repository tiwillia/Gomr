package gomr

import (
	"database/sql"
	"errors"
	"regexp"
	"strconv"

	"github.com/go-gorp/gorp"
)

type KarmaPlugin struct {
	Db   *gorp.DbMap
	Nick string
}

type Karma struct {
	Id     int    `db:"id, primarykey, autoincrement"`
	User   string `db:"user, size:500"`
	Points int    `db:"points"`
}

func (kp KarmaPlugin) Register() (err error) {
	return nil
}

func (kp KarmaPlugin) Parse(sender, channel, input string, conn *Connection) (err error) {
	if Match(input, `(?i)`+kp.Nick+`[\S]?\s+rank`) {
		if Match(input, `(?i)`+kp.Nick+`[\S]?\s+rank\s+[\S]+`) {
			urgx, _ := regexp.Compile(`\s+rank\s+([\S]+)`)
			umatch := urgx.FindStringSubmatch(input)
			if umatch != nil && len(umatch) > 1 {
				user := umatch[1]
				rank, points, err := kp.FindRank(user)
				if err != nil {
					if err == sql.ErrNoRows {
						conn.SendTo(channel, user+" has never had karma modified.")
						return nil
					}
					return err
				}
				conn.SendTo(channel, user+" is "+rank+" with "+strconv.Itoa(points)+" points of karma")
			}
		} else {
			klist, err := kp.GetKarmaByPoints()
			if err != nil {
				return err
			}

			if len(klist) == 0 {
				return nil
			}

			for i, k := range klist {
				rank := addSuffix(i + 1)
				conn.SendTo(sender, rank+") "+k.User+" with "+strconv.Itoa(k.Points)+" points")
				if i > 9 {
					break
				}
			}
		}
		return nil
	}

	if !Match(input, `\S+(\+|-){2,}`) {
		return nil
	}

	if channel == sender {
		conn.SendTo(sender, "Karma can only be modified in a public channel.")
		return nil
	}

	change := 0
	var user string
	user = MatchAndPull(input, `\S+\+\+`, `(\S+)\+\+`)
	if user != "" {
		change = 1
	} else {
		user = MatchAndPull(input, `\S+\-\-`, `(\S+)\-\-`)
		if user != "" {
			change = -1
		}
	}
	if change != 0 {
		if user == sender {
			conn.SendTo(channel, "I will not allow you to modify your own karma "+sender+".")
			return nil
		}
		var k Karma
		k, err = kp.FindOrCreateKarma(user)
		if err != nil {
			return errors.New("Unable to find or create karma entry:" + err.Error())
		}
		k.Points = k.Points + change
		err = kp.Update(k)
		if err != nil {
			return errors.New("Unable to update karma entry:" + err.Error())
		}
		conn.SendTo(channel, user+" now has "+strconv.Itoa(k.Points)+" karma.")
	}
	return nil
}

func (kp KarmaPlugin) Help() (texts []string) {
	texts = append(texts, "<name>++ or <name>--")
	texts = append(texts, kp.Nick+"[:] rank")
	texts = append(texts, kp.Nick+"[:] rank <user>")
	return texts
}

func (kp KarmaPlugin) FindRank(user string) (rank string, points int, err error) {
	var k Karma
	err = kp.Db.SelectOne(&k, "select * from karma where user=?", user)
	if err != nil {
		return
	}
	klist, err := kp.GetKarmaByPoints()
	if err != nil {
		return
	}

	var ranknum int
	for i, _ := range klist {
		if klist[i].User == user {
			k = klist[i]
			r := i + 1
			ranknum = r
		}
	}
	points = k.Points
	rank = addSuffix(ranknum)

	return
}

func (kp KarmaPlugin) FindOrCreateKarma(u string) (k Karma, err error) {
	err = kp.Db.SelectOne(&k, "select * from karma where user=?", u)
	if err != nil {
		if err == sql.ErrNoRows {
			k.Points = 0
			k.User = u
			err = kp.Db.Insert(&k)
			if err != nil {
				return
			}
		}
	}
	return
}

func (kp KarmaPlugin) GetKarmaByPoints() (klist []Karma, err error) {
	_, err = kp.Db.Select(&klist, "select * from karma order by points DESC")
	return
}

func (kp KarmaPlugin) Update(k Karma) (err error) {
	var rowCnt int64
	rowCnt, err = kp.Db.Update(k)
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return sql.ErrNoRows
	}
	return nil
}
