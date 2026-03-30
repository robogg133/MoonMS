package tests

import (
	"fmt"
	"testing"

	"github.com/robogg133/MoonMS/mscript/interpreter/lexer"
	"github.com/robogg133/MoonMS/mscript/interpreter/parser"
)

func TestParse(t *testing.T) {

	l := lexer.New(`plugin "economy" version "1.0" {
    export {
        get_balance(player): number
        add_balance(player, amount: number)
        transfer(from, to, amount: number): bool
		// Coment
    }

    emits {
        transaction(src, dst: Player, ammount: number)
    }

    name = "EconomyPlugin"
    description = "asdalçsdasdaskldaskl"
    license = "MIT"
    authors = ["robogg133"]
    mc_version = "1.21.11"
}

require database as db
require teleport as tp

include "helpers.mscript"


let plr_money = db.get_collection("money")


// every arg will be a string because it's a command
command "send" (target, amount) { 
    let plr = get_player_by_name(target)
    let n = to_number(amount)


    let money = plr_money.get(plr.uuid)
    plr_money.set(plr.uuid, money+n)

    plr.whisper("You received "+ amount)

    if this.caller is Player {
        this.caller.whisper("Sended "+amount+" to "+target)
    }

    emit transaction(this.caller, plr, n)
}


on server_stop() {
    log("Goodbye world!")
}

on server_start {
    after 100t {
        log("100 ticks!")
    }

    after 10 * TIME_SECOND {
        log("10 seconds after 10 ticks")
    }
    dizerOi("oi")
}



fn dizerOi(s: string) {
	print("oi")
}
`)

	p := parser.New(l)

	program := p.Parse()

	fmt.Println(program.Statements[0])
}
