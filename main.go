// This is a multipurpose jungletrain bot
package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/irlndts/go-discogs"
	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

// Requires ENVs to be set:
// NICK (`nick` string to assume - `bassface`)
// SERV (`server` string to connect to - `<server>:<port>`)
// JOIN (`channel` string to join e.g. `#channelName:password`)
// REPORT_TO (`nick` unquoted space-separated list to DM when taking action - `user1 user2`)
func main() {
	// set up irc session
	var NICK = os.Getenv("NICK")
	var SERV = os.Getenv("SERVER")

	hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
	}
	channels := func(bot *hbot.Bot) {
		bot.Channels = []string{os.Getenv("JOIN")}
	}
	irc, err := hbot.NewBot(SERV, NICK, hijackSession, channels)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(badwords)
	irc.AddTrigger(bassface)
	irc.AddTrigger(boh)
	irc.AddTrigger(boobs)
	irc.AddTrigger(discogsQuery)
	irc.AddTrigger(register)
	irc.AddTrigger(streams)
	irc.AddTrigger(whargwarn)
	// irc.Logger.SetHandler(log.StdoutHandler)
	// logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	// or
	// irc.Logger.SetHandler(logHandler)
	// or
	irc.Logger.SetHandler(log.StreamHandler(os.Stdout, log.JsonFormat()))

	// Start up bot (this blocks until we disconnect)
	irc.Run()
	fmt.Println("Bot shutting down.")
}

// badwords
// if message sender is an op or system, bail
// if a message contains a string from BAD_WORDS, ban message sender, send DM(s) to REPORT_TO
var badwords = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		var ignore = []string{"user1", "user2"}
		for _, i := range ignore {
			if strings.Contains(m.From, i) {
				fmt.Println("bad word ignored")
				//fmt.Println(m.Content)
				return false
			}
		}
		for _, b := range strings.Fields(os.Getenv("BAD_WORDS")) {
			if strings.Contains(strings.ToLower(m.Content), strings.ToLower(b)) {
				fmt.Println("triggered")
				fmt.Println("m.Content")
				return true
			}
		}
		return false
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		var reportTo = strings.Fields(os.Getenv("REPORT_TO"))
		for _, r := range reportTo {
			irc.Send("PRIVMSG " + r + " : " + "attempted to ban: <" + m.From + "> in channel <" + m.To + "> due to badword match")
		}
		irc.ChMode(m.From, m.To, "+b")
		return false
	},
}

// bassface hello world
var bassface = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!bassface")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		input := strings.Fields(m.Content)
		for i, f := range input {
			if f == "hello" && i == 1 { // i.e `!bassface hello`
				irc.Reply(m, "KONICHIWA, BITCHES: index <"+strconv.Itoa(i)+">")
			}
		}
		return false
	},
}

// boh to random `LIVE FROM` messages from `mc_okkie`
var boh = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.From == "mc_okkie" && strings.Contains(m.Content, "LIVE from")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		if rand.Intn(100) > 95 {
			irc.Msg(m.To, "boh \\m/")
			return true
		}
		return false
	},
}

// hehe, boobs
var boobs = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && (strings.Contains(strings.ToLower(m.Content), "boob") ||
			strings.Contains(strings.ToLower(m.Content), "bewb") ||
			strings.Contains(strings.ToLower(m.Content), "tittay"))
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		// randomly display boobs
		boobs := [11]string{
			"( . )( . )",
			"(.)(.)",
			"( ^ )( ^ )",
			"( ^ )( v )",
			"( _ )( _ )",
			"(. Y .)",
			"( o )( o )",
			"(o)(o)",
			"( . )( : )",
			"( . )( o )",
			"( . Y . Y . )",
		}
		irc.Msg(m.To, fmt.Sprint(boobs[rand.Intn(len(boobs))]))
		return false
	},
}

// discogs search
var discogsQuery = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!discogs")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		// if m.Content begins with `!discogs` and has at least 3 args
		if strings.HasPrefix(m.Content, "!discogs") && len(strings.Fields(m.Content)) < 3 {
			irc.Msg(m.To, "Usage: '!discogs <command> <query>' where '<command>' is one of: 'artist|catno|search', and '<query>' is a search string (e.g. 'rick astley', '7THGLOBEX 003' or 'promised land volume one')")
			return true
		}
		// init discogs api
		var TOKEN = os.Getenv("DISCOGS_TOKEN")
		client, err := discogs.New(&discogs.Options{
			UserAgent: "Bassface/0.1 +https://jungletrain.net",
			Token:     TOKEN,
			URL:       "https://api.discogs.com",
		})
		if err != nil {
			// notify but continue
			fmt.Printf("Discogs not initialized: %v", err)
		}
		// r[0] == '!discogs'
		// r[1] == 'command'
		// r[2:] == 'the query'
		r := strings.Fields(m.Content)
		// 'hardcore' is a tough style to match against since it also means hardcore punk rock.
		// Seems that any electronic hardcore is accompanied by 'breakbeat', so we'll put that forward in the search string
		var styles = []string{
			"Jungle",
			"Drum n Bass",
			"Breakbeat",
			"Breaks",
			"Dub",
			"Techno",
			"Acid",
			"Ambient",
		}
		switch r[1] {
		case "test":
			//master, err := client.Database.Master(718441)
			//if err != nil {
			//	fmt.Println(err)
			//}
			// fmt.Printf("%+v\n", master)
		case "artist":
			params := discogs.SearchRequest{Q: strings.Join(r[2:], " "), Type: "artist", Page: 0, PerPage: 1}
			search, err := client.Search(params)
			if err != nil {
				fmt.Printf("\n search err: %v", err)
				return true
			}
			for _, r := range search.Results {
				artist, err := client.Artist(r.ID)
				if err != nil {
					fmt.Printf("\n search err: %x", err)
					return true
				}
				irc.Msg(m.To, "artist: "+artist.Name+": https://discogs.com"+r.URI)
			}
		case "catno":
			for _, s := range styles {
				params := discogs.SearchRequest{Catno: strings.Join(r[2:], ""), Style: s, Type: "release", Page: 0, PerPage: 1}
				search, err := client.Search(params)
				if err != nil {
					fmt.Printf("\n search err: %v", err)
					return true
				}
				for _, r := range search.Results {
					irc.Msg(m.To, r.Label[0]+": "+r.Title+": https://discogs.com"+r.URI)
					return true
				}
			}
		case "search":
			for _, s := range styles {
				params := discogs.SearchRequest{Q: strings.Join(r[2:], " "), Style: s, Type: "release", Page: 0, PerPage: 1}
				search, err := client.Search(params)
				if err != nil {
					fmt.Printf("\n search err: %v", err)
					return true
				}
				for _, r := range search.Results {
					irc.Msg(m.To, r.Label[0]+": "+r.Title+": https://discogs.com"+r.URI)
					return true
				}
			}
		default:
			fmt.Printf("discogs: unknown command: %v", r[1])
		}
		return false
	},
}

// if NickServ asks for a password, ship it
// nick registration seems currently broken (no emails being sent) -- fix that, then we can register and assert password at connect
var register = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.From == "NickServ" && strings.Contains(m.Content, "This nickname is registered and protected")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		var PASSWORD = os.Getenv("PASSWORD")
		fmt.Println("REGISTERING")
		irc.Msg("NickServ", "IDENTIFY "+PASSWORD)
		return false
	},
}

// output jungletrain relay playlist for external media players
var streams = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && (strings.HasPrefix(m.Content, "!streams") || strings.HasPrefix(m.Content, "!pls"))
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Msg(m.To, "Jungletrain Relay Playlist: https://jungletrain.net/128kbps.pls")
		return false
	},
}

// whargwarn response
var whargwarn = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		re := regexp.MustCompile(`[Ww][h]?[a]+[r]?gw[a]+([r]+)?n+`)
		return m.Command == "PRIVMSG" && re.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		if rand.Intn(100) > 50 {
			irc.Msg(m.To, "Whargwarnnnnn Rudebwoys & Gyals!")
			return false
		} else {
			irc.Action(m.To, "sparks a lighter")
			return false
		}
	},
}
