// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package cli

import(
	`github.com/michaeldv/donna`
	`fmt`
	`time`
)

func Repl(engine *donna.Engine) {
	var game *donna.Game
	var position *donna.Position

	setup := func() {
		if game == nil || position == nil {
			engine.Set(`cache`, 64, `movetime`, 5000) // 5s per move.
			game = donna.NewGame()
			position = game.Start()
			fmt.Printf("%s\n", position)
		}
	}

	think := func() {
		if move := game.Think(); move != 0 {
			position = position.MakeMove(move)
			fmt.Printf("%s\n", position)
		}
	}

	benchmark := func() {
		engine.Set(`cache`, 64, `depth`, 11)

		// Sergey Kaminer, 1935
		// 1.h8Q+ Kxh8 2.Ng4+ Kg7 3.Qh7+ Kf8 4.Qg8+ Ke7 5.Qe8+ Kd6 6.Qd7+ Kc5 7.Qb5+ Kd6
		// 8.Qb6+ Ke7 9.Qc7+ Kf8 10.Qc8+ Ke7 11.Qd7+ Kf8 12.Qe8+ Kg7 13.Qg8+ Kxg8 14.Kxf6+
		game := donna.NewGame(`Kd1,Qh2,Nh6,a4,g3,h7`, `Kg7,Qe4,Bf6,b7,e6,g6`)
		position := game.Start(donna.White)
		fmt.Printf("%s\n", position)
		game.Think()

		// Tylkowski vs. Wojciechowski, Poznan 1931
		// 30...Rxb2! 31. Nxb2 c3 32. Rxb6 c4!! 33. Rb4 a5! 34.Na4 axb4
		game = donna.NewGame(`Kg1,Rb7,Na4,a2,b2,f4,g2,g3`, `Kh7,Rd2,Bb6,a7,c5,c4,g7,h6`)
		position = game.Start(donna.Black)
		fmt.Printf("%s\n", position)
		game.Think()

		// Georg Rotlewi vs. Akiba Rubinstein, Lodz 1907, after 21 moves.
		// http://www.chessgames.com/perl/chessgame?gid=1119679
		game = donna.NewGame(`Kh1,Qe2,Ra1,Rf1,Bb2,Be4,Nc3,a3,b4,e5,f4,g2,h2`, `Kg8,Qe7,Rc8,Rd8,Bb6,Bb7,Ng4,a6,b5,e6,f7,g7,h7`)
		position = game.Start(donna.Black)
		fmt.Printf("%s\n", position)
		game.Think()

		// Donald Byrne vs. Bobby Fischer, Third Rosenwald Trophy, 1956 after 16 moves.
		// http://www.chessgames.com/perl/chessgame?gid=1008361
		// 16... Rf8-e8 17. Ke1-f1 Bg4-e6!!
		game = donna.NewGame(`Ke1,Qa3,Rd1,Rh1,Bc4,Bc5,Nf3,a2,d4,f2,g2,h2`, `Kg8,Qb6,Ra8,Rf8,Bg4,Bg7,Nc3,a7,b7,c6,f7,g6,h7`)
		position = game.Start(donna.Black)
		fmt.Printf("%s\n", position)
		game.Think()

		// Bobby Fischer vs. James Sherwin, donna.New Jersey Open 1957, after 16 moves.
		// http://www.chessgames.com/perl/chessgame?gid=1008366
		// Fischer played 17. h2-h4!
		game = donna.NewGame(`Kg1,Qc2,Ra1,Re1,Bc1,Bg2,Ng5,a2,b2,c3,d4,f2,g3,h2`, `Kg8,Qd6,Ra8,Rf8,Bc8,Nd5,Ng6,a7,b6,c4,e6,f7,g7,h7`)
		position = game.Start(donna.White)
		fmt.Printf("%s\n", position)
		game.Think()

		// Mikhail Botvinnik vs. Jose Raul Capablanca, AVRO 1936, after 29 moves.
		// Botvinnik played 30. Bb2-a3!
		game = donna.NewGame(`Kg1,Qe5,Bb2,Ng3,c3,d4,e6,g2,h2`, `Kg7,Qe7,Nb3,Nf6,a7,b6,c4,d5,g6,h7`)
		position = game.Start(donna.White)
		fmt.Printf("%s\n", position)
		game.Think()
	}

	perft := func(depth int) {
		position := donna.NewGame().Start()
		start := time.Now()
		total := position.Perft(depth)
		finish := time.Since(start).Seconds()
		fmt.Printf("\n  Nodes: %d\n", total)
		fmt.Printf("Elapsed: %.2fs\n", finish)
		fmt.Printf("Nodes/s: %.2f\n", float64(total)/finish)
	}

	for command := ``; ; command = `` {
		fmt.Print(`donna> `)
		fmt.Scanf(`%s`, &command)
		switch command {
		case ``:
		case `bench`:
			benchmark()
		case `exit`, `quit`:
			return
		case `go`:
			setup()
			think()
		case `help`, `?`:
			fmt.Println("The commands are:\n\n" +
				"   bench   Run benchmark test\n" +
				"   exit    Exit the program\n" +
				"   go      Take side and make a move\n" +
				"   help    Display this help\n" +
				"   new     Start new game\n" +
				"   perft   Run perft test\n" +
				"   score   Show evaluation summary\n" +
				"   undo    Undo last move\n")
		case `new`:
			game, position = nil, nil
			setup()
		case `perft`:
			perft(5)
		case `score`:
			setup()
			_, metrics := position.EvaluateWithTrace()
			donna.Summary(metrics)
		case `undo`:
			if position != nil {
				position = position.UndoLastMove()
				fmt.Printf("%s\n", position)
			}
		default:
			setup()
			if move := donna.NewMoveFromString(position, command); move != 0 {
				if advance := position.MakeMove(move); advance != nil {
					position = advance
					think()
					continue
				}
			}
			// Invalid move (typo) or non-evasion on check.
			fmt.Printf("%s appears to be an invalid move.\n", command)
		}
	}
}
