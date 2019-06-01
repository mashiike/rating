package golicko

//Glicko-2 Rating: see as http://www.glicko.net/glicko/glicko2.pdf

import (
	"fmt"
	"math"
)

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
	//Value is a number representing the strength of the player.
	//default value is 1500.0
	Value float64 `json:"value" yaml:"value" csv:"value"`

	//Deviation is a rating deviation, which is used to express 95% confidence intervals.
	//default value is 350.0
	Deviation float64 `json:"deviation" yaml:"deviation" csv:"deviation"`

	//Volatility measure indicates the degree of expected fluctuation in a player’s rating.
	//default value is 0.06
	Volatility float64 `json:"volatility" yaml:"volatility" csv:"volatility"`
}

//DefaultRating is fill by default values
var DefaultRating Rating = Rating{
	Value:      1500.0,
	Deviation:  350.0,
	Volatility: 0.06,
}

//String is for dump. fmt.Stringer interface implements
func (r Rating) String() string {
	min, max := r.Interval()
	return fmt.Sprintf("rating:%0.2f (%0.2f-%0.2f), volatility:%0.6f", r.Value, min, max, r.Volatility)
}

//Interval return 95% confidence interval.
func (r Rating) Interval() (float64, float64) {
	return r.Value - r.Deviation*2, r.Value + r.Deviation*2
}

//IsDiff is a function to check the significance of Rating
func (r Rating) IsDiff(o Rating) bool {
	y := r.Value - o.Value
	z := y / math.Sqrt(r.Deviation*r.Deviation+o.Deviation*o.Deviation)
	if math.Abs(z) > 1.96 {
		return true
	}
	return false
}

//IsStronger is checker function. this rating r is storonger than rating o.
func (r Rating) IsStronger(o Rating) bool {
	if r.Value-o.Value <= 0.0 {
		return false
	}
	return r.IsDiff(o)
}

//IsWeeker is checker function. this rating r is weeker than rating o.
func (r Rating) IsWeeker(o Rating) bool {
	if r.Value-o.Value >= 0.0 {
		return false
	}
	return r.IsDiff(o)
}

// glicko2Scale is internal rating value
type glicko2Scale struct {
	Mu    float64
	Phi   float64
	Sigma float64
}

// ToGlicko2 is convert to Glicko-2 Scale from Rating. as Step.2
func (r Rating) toGlicko2() glicko2Scale {
	return glicko2Scale{
		Mu:    (r.Value - valueCenterValue) / valueConvateRate,
		Phi:   r.Deviation / valueConvateRate,
		Sigma: r.Volatility,
	}
}

// ToRating is convert to Rating from Glicko-2 Scale. as Step.8
func (s glicko2Scale) toRating() Rating {
	return Rating{
		Value:      s.Mu*valueConvateRate + valueCenterValue,
		Deviation:  s.Phi * valueConvateRate,
		Volatility: s.Sigma,
	}
}

// Result is matches result
type Result struct {
	Opponent Rating
	Score    Score
}

// Setting is Update setting
type Setting struct {
	Tau float64
}

var DefaultSetting Setting = Setting{Tau: 0.5}

// Update is compute new Rating from Results and Setting
func (r Rating) Update(results []Result, setting Setting) Rating {
	scale := r.toGlicko2()
	quant := scale.ComputeQuantity(results)
	newScale := scale.ApplyQuantity(quant, setting)
	return newScale.toRating()

}

// estimatedQuantity is quantity for compute next rating.
type estimatedQuantity struct {
	Variance    float64
	Improvement float64
}

// ComputeQuantity is Step.3 and Step.4. Compute the quantity of Estimated Improvement and Estimated Variance.
func (s glicko2Scale) ComputeQuantity(results []Result) estimatedQuantity {
	invEV := 0.0
	tmpImp := 0.0
	for _, r := range results {
		opponent := r.Opponent.toGlicko2()
		valg := fg(opponent.Phi)
		valE := fE(s.Mu, opponent.Mu, opponent.Phi)
		invEV += valg * valg * valE * (1.0 - valE)
		tmpImp += valg * (float64(r.Score) - valE)
	}
	return estimatedQuantity{
		Variance:    1.0 / invEV,
		Improvement: tmpImp / invEV,
	}
}

// ApplyQuantity is Step.5 and Step.6, Step.7. Compute the next Glicko-2 scale.
func (s glicko2Scale) ApplyQuantity(quantity estimatedQuantity, setting Setting) glicko2Scale {
	if math.IsInf(quantity.Variance, 0) {
		// if estimated variance is infinity, can not apply. becouse maybe no result.
		// In this case, rating value and volatility parameters remain the same, but the rating deviation increases
		return glicko2Scale{
			Mu:    s.Mu,
			Phi:   math.Sqrt(s.Phi*s.Phi + s.Sigma*s.Sigma),
			Sigma: s.Sigma,
		}
	}

	sigmaDash := s.determineSigma(quantity, setting.Tau)
	phiAsta := math.Sqrt(s.Phi*s.Phi + sigmaDash*sigmaDash)
	phiDash := 1.0 / math.Sqrt(1.0/(phiAsta*phiAsta)+1.0/quantity.Variance)

	return glicko2Scale{
		Mu:    s.Mu + phiDash*phiDash*quantity.Improvement/quantity.Variance,
		Phi:   phiDash,
		Sigma: sigmaDash,
	}
}

func (s glicko2Scale) determineSigma(quantity estimatedQuantity, tau float64) float64 {

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
