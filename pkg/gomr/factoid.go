package gomr

import (
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/go-gorp/gorp"
)

type FactoidPlugin struct {
	// The plugin will silently ignore the following words
	Blacklist []string
	Db        *gorp.DbMap
	Nick      string
}

type Factoid struct {
	Id           int    `db:"id, primarykey, autoincrement"`
	Fact         string `db:"fact, size:100"`
	Definition   string `db:"definition, size:1000"`
	CreationDate int64  `db:"creation_date"`
}

func (fp FactoidPlugin) Register() (err error) {
	return nil
}

func (fp FactoidPlugin) Parse(sender, channel, input string, conn *Connection) (err error) {
	// Check for factoid retrieval match
	if Match(input, `^\S+\?\r$`) || Match(input, `(?i)^`+fp.Nick+`:*\s+\S+\?\r$`) {
		var frgxStr string
		if Match(input, `(?i)^`+fp.Nick) {
			frgxStr = `(?i)^` + fp.Nick + `:*\s+(\S+)\?\r$`
		} else {
			frgxStr = `^(\S+)\?\r$`
		}
		frgx := regexp.MustCompile(frgxStr)
		fmatch := frgx.FindStringSubmatch(input)

		if fmatch != nil && len(fmatch) > 1 {
			if len(fmatch) > 2 {
				return errors.New("Found more than one match for factoid, stopping processing.")
			}
			fact := fmatch[1]

			// Loop through the blacklist and end it if a match is found
			for _, blWord := range fp.Blacklist {
				if fact == blWord {
					return nil
				}
			}

			var factoids []Factoid
			factoids, err = fp.GetFactoids(fact)
			if err != nil {
				return err
			}

			for i := range factoids {
				conn.SendTo(channel, "#"+strconv.Itoa(i+1)+" "+fact+": "+factoids[i].Definition)
			}
			return nil
		}
	}

	// Check for factoid set match
	setrgxStr := `(?i)^` + fp.Nick + `:*\s+(\S+) is\s+(\S+.*)\r$`
	if Match(input, setrgxStr) {
		srgx := regexp.MustCompile(setrgxStr)
		smatch := srgx.FindStringSubmatch(input)
		if smatch != nil && len(smatch) > 2 {
			fact := smatch[1]
			def := smatch[2]

			// Loop through the blacklist and end it if a match is found
			for _, blWord := range fp.Blacklist {
				if fact == blWord {
					return nil
				}
			}

			utime := time.Now().Unix()
			factoid := Factoid{Fact: fact, Definition: def, CreationDate: utime}
			err = fp.Create(factoid)
			if err != nil {
				return err
			}
			conn.SendTo(channel, "Ok, I'll remember "+fact)
			return nil
		}
	}

	// Check for factoid forget match
	frgxStr := `(?i)^` + fp.Nick + `:*\s+forget\s+(\S+)\s*([0-9]*)\r$`
	if Match(input, frgxStr) {
		frgx := regexp.MustCompile(frgxStr)
		fmatch := frgx.FindStringSubmatch(input)
		if fmatch != nil && len(fmatch) > 1 {
			if fmatch[2] != "" {
				// id was provided
				id, _ := strconv.Atoi(fmatch[2])
				fact := fmatch[1]
				factoids, err := fp.GetFactoids(fact)
				if err != nil {
					return err
				}

				if len(factoids) == 0 {
					conn.SendTo(channel, fact+" has never been defined.")
					return nil
				}

				if len(factoids) < id {
					conn.SendTo(channel, "No definition for "+fact+" exists with ID: "+strconv.Itoa(id))
					return nil
				}

				err = fp.Delete(factoids[id-1])
				if err != nil {
					return err
				}

				conn.SendTo(channel, "Deleted definition for "+fact+" with ID: "+strconv.Itoa(id))
			} else {
				// id not provided - delete the latest
				fact := fmatch[1]
				factoids, err := fp.GetFactoids(fact)
				if err != nil {
					return err
				}

				if len(factoids) == 0 {
					conn.SendTo(channel, fact+" has never been defined.")
					return nil
				}

				err = fp.Delete(factoids[len(factoids)-1])
				if err != nil {
					return err
				}

				conn.SendTo(channel, "Deleted latest definition of "+fact)
			}
		}
	}
	return nil
}

func (fp FactoidPlugin) Help() (texts []string) {
	texts = append(texts, fp.Nick+"[:] <fact> is <definition>")
	texts = append(texts, fp.Nick+"[:] <fact>?")
	texts = append(texts, "<fact>?")
	texts = append(texts, fp.Nick+"[:] forget <fact> [n]")
	return texts
}

func (fp FactoidPlugin) Create(f Factoid) (err error) {
	err = fp.Db.Insert(&f)
	return
}

func (fp FactoidPlugin) Delete(f Factoid) (err error) {
	var rowcnt int64
	rowcnt, err = fp.Db.Delete(&f)
	if rowcnt == 0 {
		return sql.ErrNoRows
	}
	return
}

func (fp FactoidPlugin) Update(f Factoid) (err error) {
	var rowCnt int64
	rowCnt, err = fp.Db.Update(&f)
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (fp FactoidPlugin) GetFactoids(fact string) (factoids []Factoid, err error) {
	_, err = fp.Db.Select(&factoids, "select * from factoids where fact=? order by creation_date ASC", fact)
	return
}
