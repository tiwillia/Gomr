package main

import (
	"errors"
	"regexp"
	// TODO remove log entries
	"log"
)

type FactoidPlugin struct {
	Blacklist []string
}

type Factoid struct {
	Id         int    `db:"id, primarykey, autoincrement"`
	Fact       string `db:"fact, size:100"`
	Definition string `db:"definition, size:1000"`
}

func (fp FactoidPlugin) Register() (err error) {
	// The plugin will silently ignore the following words
	fp.Blacklist = []string{"why", "where", "who", "when", "how", "now"}
	return nil
}

func (fp FactoidPlugin) Parse(sender, channel, input string, conn *Connection) (err error) {
	// Check for factoid retrieval match
	if Match(input, `^\S+\?\r$`) || Match(input, `^`+config.Nick+`:*\s+\S+\?\r$`) {
		log.Println("Matched!")
		var frgxStr string
		if Match(input, `^`+config.Nick) {
			frgxStr = `^` + config.Nick + `[:]\s+(\S+)\?\r$`
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

			var defs []string
			defs, err := getDefinitions(fact)
			if err != nil {
				return err
			}

			for i := range defs {
				conn.SendTo(channel, fact+": "+defs[i])
			}
			return nil
		}
	}

	// Check for factoid set match
	setrgxStr := `^` + config.Nick + `[:]\s+(\S+) is\s+(\S+.*)\r$`
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

			factoid := Factoid{Fact: fact, Definition: def}
			err = factoid.Create()
			if err != nil {
				return err
			}
			conn.SendTo(channel, "Ok, I'll remember "+fact)
			return nil
		}
	}
	return nil
}

func (fp FactoidPlugin) Help() (texts []string) {
	texts = append(texts, config.Nick+"[:] <fact> is <definition>")
	texts = append(texts, config.Nick+"[:] <fact>?")
	texts = append(texts, "<fact>?")
	return texts
}

func (f *Factoid) Create() (err error) {
	err = Db.Insert(f)
	return
}

func (f *Factoid) Update() (err error) {
	var rowCnt int64
	rowCnt, err = Db.Update(f)
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return ErrNoRowsUpdated
	}
	return nil
}

func getDefinitions(fact string) (defs []string, err error) {
	fs := []Factoid{}
	_, err = Db.Select(&fs, "select * from factoids where fact=?", fact)
	if err != nil {
		return
	}
	for i := range fs {
		defs = append(defs, fs[i].Definition)
	}
	return
}
