// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Quiescence search.
func (p *Position) searchQuiescence(alpha, beta, depth int) int {
	return p.searchQuiescenceWithFlag(alpha, beta, depth, false)
}

func (p *Position) searchQuiescenceWithFlag(alpha, beta, depth int, capturesOnly bool) (score int) {
	ply := Ply()

	if ply >= MaxPly || p.game.clock.halt {
		return p.Evaluate()
	}

	p.game.pv[ply] = p.game.pv[ply][:0]

	// Probe cache.
	cacheFlags := uint8(cacheAlpha)
	if cached := p.probeCache(); cached != nil {
		if cached.depth >= depth {
			score := cached.score
			if score > Checkmate - MaxPly && score <= Checkmate {
				score -= ply
			} else if score >= -Checkmate && score < -Checkmate + MaxPly {
				score += ply
			}

			// if cached.flags == cacheExact {
			// 	return score
			// } else if cached.flags == cacheAlpha && score <= alpha {
			// 	return alpha
			// } else if cached.flags == cacheBeta && score >= beta {
			// 	return beta
			// }
			if cached.flags == cacheExact ||
			   cached.flags == cacheAlpha && score <= alpha ||
			   cached.flags == cacheBeta && score >= beta {
				return score
			}

		}
	}

	inCheck := p.isInCheck(p.color)
	staticScore := p.Evaluate()
	if !inCheck && staticScore > alpha {
		alpha = staticScore
	}
	if alpha >= beta {
		return beta
	}

	gen := NewGen(p, ply)
	if inCheck {
		gen.generateEvasions()
	} else {
		gen.generateCaptures()
	}
	gen.quickRank()


	moveCount, bestMove := 0, Move(0)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if !inCheck && p.exchange(move) < 0 {
			continue
		}
		if position := p.MakeMove(move); position != nil {
			moveCount++
			score = -position.searchQuiescenceWithFlag(-beta, -alpha, depth, true)
			position.TakeBack(move)

			if score > alpha {
				alpha = score
				bestMove = move
				cacheFlags = cacheExact

				if alpha >= beta {
					cacheFlags = cacheBeta
					break
				}
				p.game.saveBest(ply, move)
			}
		}
	}

	if !inCheck && !capturesOnly {
		gen = NewGen(p, Ply()).generateChecks().quickRank()
		for move := gen.NextMove(); move != 0; move = gen.NextMove() {
			if p.exchange(move) < 0 {
				continue
			}
			if position := p.MakeMove(move); position != nil {
				moveCount++
				score = -position.searchQuiescenceWithFlag(-beta, -alpha, depth, false)
				position.TakeBack(move)

				if p.game.clock.halt {
					p.game.qnodes += moveCount
					//Log("searchQui at %d (%s): move %s (%d) score %d alpha %d\n", depth, C(p.color), move, moveCount, score, alpha)
					return alpha
				}

				if score > alpha {
					alpha = score
					bestMove = move
					cacheFlags = cacheExact

					if alpha >= beta {
						cacheFlags = cacheBeta
						break
					}
					p.game.saveBest(ply, move)
				}
			}
		}
	}

	p.game.qnodes += moveCount

	score = alpha
	if inCheck && moveCount == 0 {
		score = -Checkmate + ply
	}
	p.cache(bestMove, score, depth, cacheFlags)

	return
}
