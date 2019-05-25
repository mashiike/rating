package golicko

//Glicko-2 Rating: see as http://www.glicko.net/glicko/glicko2.pdf

import "math"

const (
	valueConvateRate = 173.7178
	valueCenterValue = 1500.0
	iterationLimit   = 100000
	epsiron          = 0.000001

	ScoreWin  = Score(1)
	ScoreLose = Score(0)
	ScoreDraw = Score(0.5)
)

// Score is enum type for Win/Lose/Draw.
type Score float64

//Opponent is Opponent Player's/Team's Score.
func (s Score) Opponent() Score {
	return Score(1 - s)
}

// Rating is a variable county that represents the strength of the player
type Rating struct {
	Value      float64
	Deviation  float64
	Volatility float64
}

// Glicko2Scale is internal rating value
type Glicko2Scale struct {
	Mu    float64
	Phi   float64
	Sigma float64
}

// ToGlicko2 is convert to Glicko-2 Scale from Rating. as Step.2
func (r Rating) ToGlicko2() Glicko2Scale {
	return Glicko2Scale{
		Mu:    (r.Value - valueCenterValue) / valueConvateRate,
		Phi:   r.Deviation / valueConvateRate,
		Sigma: r.Volatility,
	}
}

// ToRating is convert to Rating from Glicko-2 Scale. as Step.8
func (s Glicko2Scale) ToRating() Rating {
	return Rating{
		Value:      s.Mu*valueConvateRate + valueCenterValue,
		Deviation:  s.Phi * valueConvateRate,
		Volatility: s.Sigma,
	}
}

func (r Rating) Equals(o Rating) bool {
	// value and deviation tolerance is to the 1st decimal point
	if math.Abs(r.Value-o.Value) > 0.1 {
		return false
	}
	if math.Abs(r.Deviation-o.Deviation) > 0.1 {
		return false
	}
	// volatility olerance is to the 4th decimal point
	if math.Abs(r.Volatility-o.Volatility) > 0.001 {
		return false
	}

	return true
}

// Result is matches result
type Result struct {
	Opponent Glicko2Scale
	Score    Score
}

// EstimatedQuantity is quantity for compute next rating.
type EstimatedQuantity struct {
	Variance    float64
	Improvement float64
}

// ComputeQuantity is Step.3 and Step.4. Compute the quantity of Estimated Improvement and Estimated Variance.
func (s Glicko2Scale) ComputeQuantity(results []Result) EstimatedQuantity {
	invEV := 0.0
	tmpImp := 0.0
	for _, r := range results {
		valg := fg(r.Opponent.Phi)
		valE := fE(s.Mu, r.Opponent.Mu, r.Opponent.Phi)
		invEV += valg * valg * valE * (1.0 - valE)
		tmpImp += valg * (float64(r.Score) - valE)
	}
	return EstimatedQuantity{
		Variance:    1.0 / invEV,
		Improvement: tmpImp / invEV,
	}
}

// ApplyQuantity is Step.5 and Step.6, Step.7. Compute the next Glicko-2 scale.
func (s Glicko2Scale) ApplyQuantity(quantity EstimatedQuantity, tau float64) Glicko2Scale {
	if math.IsInf(quantity.Variance, 0) {
		// if estimated variance is infinity, can not apply. becouse maybe no result.
		// In this case, rating value and volatility parameters remain the same, but the rating deviation increases
		return Glicko2Scale{
			Mu:    s.Mu,
			Phi:   math.Sqrt(s.Phi*s.Phi + s.Sigma*s.Sigma),
			Sigma: s.Sigma,
		}
	}

	sigmaDash := s.determineSigma(quantity, tau)
	phiAsta := math.Sqrt(s.Phi*s.Phi + sigmaDash*sigmaDash)
	phiDash := 1.0 / math.Sqrt(1.0/(phiAsta*phiAsta)+1.0/quantity.Variance)

	return Glicko2Scale{
		Mu:    s.Mu + phiDash*phiDash*quantity.Improvement/quantity.Variance,
		Phi:   phiDash,
		Sigma: sigmaDash,
	}
}

func (s Glicko2Scale) determineSigma(quantity EstimatedQuantity, tau float64) float64 {

	a := math.Log(s.Sigma * s.Sigma)
	largeA := a
	largeB := 0.0

	sqDelta := quantity.Improvement
	sqPhi := s.Phi * s.Phi

	switch {
	case sqDelta > sqPhi+quantity.Variance:
		largeB = math.Log(sqDelta - sqPhi - quantity.Variance)
	default:
		valf := 0.0
		for k := 1; k < iterationLimit+1; k++ {
			largeB = (a - float64(k)*tau)
			valf = fx(largeB, a, sqDelta, sqPhi, quantity.Variance, tau)
			if valf >= 0.0 {
				break
			}
		}
	}

	valfA := fx(largeA, a, sqDelta, sqPhi, quantity.Variance, tau)
	valfB := fx(largeB, a, sqDelta, sqPhi, quantity.Variance, tau)

	for i := 0; i < iterationLimit; i++ {
		if math.Abs(largeB-largeA) <= epsiron {
			break
		}
		largeC := largeA + ((largeA-largeB)*valfA)/(valfB-valfA)
		valfC := fx(largeC, a, sqDelta, sqPhi, quantity.Variance, tau)
		switch {
		case valfB*valfC < 0.0:
			largeA = largeB
			valfA = valfB
		default:
			valfA /= 2.0
		}
		largeB = largeC
		valfB = valfC
	}

	return math.Exp(largeA / 2.0)
}

//fE is internal function E(mu, mu_j, phi_j)
func fE(mu, muOppo, phiOppo float64) float64 {
	return 1.0 / (1.0 + math.Exp(-1.0*fg(phiOppo)*(mu-muOppo)))
}

//fg is internal function g(phi)
func fg(phi float64) float64 {
	return 1.0 / math.Sqrt(1.0+3.0*(phi*phi)/(math.Pi*math.Pi))
}

//fx is f(x) ,target function for newton-raphson method
func fx(x, a, sqDeltaMu, sqPhi, ev, tau float64) float64 {
	sumVal := sqDeltaMu + sqPhi + ev
	diffVal := sqDeltaMu - sqPhi - ev
	firstTerm := (math.Exp(x) * (diffVal - math.Exp(x))) / (2 * sumVal * sumVal)
	secondTerm := (x - a) / (tau * tau)
	return firstTerm - secondTerm
}
