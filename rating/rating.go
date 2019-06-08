// Package rating implements the simple glicko-2 Rating.
// Glicko-2 Rating: see as following. (within this package, denoted as ref[1] below)
// Professor Mark E. Glickman. "Example of the Glicko-2 system" http://www.glicko.net/glicko/glicko2.pdf
// Note: The variable names in this package match the mathematical Greek letter variables of the dissertation.
package rating

import (
	"encoding/binary"
	"math"
	"sync"

	"github.com/pkg/errors"
)

const (
	//rating <=> glicko2 scale deviation convart rate
	convartRate = 173.7178
	//rating center value
	centerValue = 1500.0
	//start deviation
	startDeviation = 350.0
	startPhi       = startDeviation / convartRate
	//x ~ N(0,1), this value is z when P(-z <= x <= z) = 0.95
	zscore95 = 1.96

	//iterative algorithm truncates with this number
	iterationLimit = 100000
	//end condition of iterative algorithm
	epsiron = 0.000001

	//ScoreWin is Score when winning an opponent.
	ScoreWin = float64(1.0)
	//ScoreLose is Score when losing to an opponent.
	ScoreLose = float64(0.0)
	//ScoreDraw is Score when tied to an opponent.
	ScoreDraw = float64(0.5)
)

//Rating is a structure to evaluate the strength of a player / team.
type Rating struct {

	// mu is a number representing strength.
	// phi is a numerical value that represents the deviation of strength. as rating deviation:RD.
	// In the glicko-2 system, x ~ N (mu, phi ^ 2) and x are evaluated with a model that considers the strength to be actually exhibited in a match
	// Note: N (0,1) is a standard normal distribution
	mu  float64
	phi float64

	// sigma is a numerical value that represents the volatility of strength.
	// If there is no war during one rating period, RD will rise by the next equation in the next period.
	//
	//   phi_next = sqrt(phi^2 + sigma^2)
	//
	// start phi is 2.014, default phi is 0.2878
	// If you want to return to start phi if you have not played for 100 rating periods.
	//
	//  sigma = sqrt((2.014^2 - 0.2878^2)/100) = 0.1997 ~= 0.2
	//
	// in ref[1], this value is 0.06
	sigma float64
}

//New is a constractor for Rating
func New(strength, deviation, volatility float64) Rating {
	return Rating{
		mu:    (strength - centerValue) / convartRate,
		phi:   deviation / convartRate,
		sigma: volatility,
	}
}

// Default is return default rating for starting Player/Team.
func Default(volatility float64) Rating {
	return New(centerValue, startDeviation, volatility)
}

// Strength is return value of strength, as rating general value.
func (r Rating) Strength() float64 {
	return nthFloor(r.mu*convartRate+centerValue, 2)
}

// Deviation is return RD(Rating Deviation) as general value.
func (r Rating) Deviation() float64 {
	return nthFloor(r.phi*convartRate, 2)
}

// Volatility is return Rating volatility.
func (r Rating) Volatility() float64 {
	return nthFloor(r.sigma, 6)
}

//Interval is return value of strength 95% confidence interval.
func (r Rating) Interval() (float64, float64) {
	s := r.Strength()
	rd2 := r.Deviation() * 2
	return s - rd2, s + rd2
}

//IsDifferent is a function to check the significance of Rating
func (r Rating) IsDifferent(o Rating) bool {
	y := r.mu - o.mu
	z := y / geometricMean(r.phi, o.phi)
	if math.Abs(z) > zscore95 {
		return true
	}
	return false
}

//IsStronger is checker function. this rating r is storonger than rating o.
func (r Rating) IsStronger(o Rating) bool {
	if r.mu <= o.mu {
		return false
	}
	return r.IsDifferent(o)
}

//IsWeeker is checker function. this rating r is weeker than rating o.
func (r Rating) IsWeeker(o Rating) bool {
	if r.mu >= o.mu {
		return false
	}
	return r.IsDifferent(o)
}

// WinProb is estimate winning probability,
// this value 1500 and 1700, both RD is 0 => P(1700 is win) = 0.76
func (r Rating) WinProb(o Rating) float64 {
	return nthFloor(fE(r.mu, o.mu, geometricMean(r.phi, o.phi)), 4)
}

func float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func byteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

const ratingBinaryVersion byte = 1

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (r Rating) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, 25)
	b = append(b, ratingBinaryVersion)
	b = append(b, float64ToByte(r.mu)...)
	b = append(b, float64ToByte(r.phi)...)
	b = append(b, float64ToByte(r.sigma)...)
	return b, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (r *Rating) UnmarshalBinary(data []byte) error {
	buf := data
	if len(buf) == 0 {
		return errors.New("Rating.UnmarshalBinary: no data")
	}

	if buf[0] != ratingBinaryVersion {
		return errors.New("Rating.UnmarshalBinary: unsupported version")
	}

	if len(buf) != 25 {
		return errors.New("Rating.UnmarshalBinary: invalid length")
	}

	r.mu = byteToFloat64(buf[1:9])
	r.phi = byteToFloat64(buf[9:17])
	r.sigma = byteToFloat64(buf[17:25])
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// The rating is a quoted string in Default format
func (r Rating) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DefaultFormat)+2)
	b = append(b, '"')
	b = r.AppendFormat(b, DefaultFormat)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The rating is expected to be a quoted string in Default format.
func (r *Rating) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	*r, err = Parse(`"`+DefaultFormat+`"`, string(data))
	return err
}

// MarshalText implements the encoding.TextMarshaler interface.
// The rating is formatted in Default Format.
func (r Rating) MarshalText() ([]byte, error) {
	b := make([]byte, 0, len(DefaultFormat))
	return r.AppendFormat(b, DefaultFormat), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The rating is expected to be in Default Format.
func (r *Rating) UnmarshalText(data []byte) error {
	var err error
	*r, err = Parse(DefaultFormat, string(data))
	return err
}

//Update is utils function for rating update. case non sequentially update, use this function.
func (r Rating) Update(opponents []Rating, scores []float64, tau float64) (Rating, error) {
	if len(opponents) != len(scores) {
		return r, errors.New("opponents and scores length missmatch")
	}
	e := NewEstimated(r, tau)
	for i := 0; i < len(opponents); i++ {
		if err := e.ApplyMatch(opponents[i], scores[i]); err != nil {
			return r, err
		}
	}
	if err := e.Fix(); err != nil {
		return r, err
	}
	return e.Fixed, nil
}

// Estimated is a collection of Quantity related estimates that are being corrected.
// If you update the rating sequentially, use this struct to save the learning process during the current rating period.
type Estimated struct {
	sync.Mutex

	// in ref[1], this value is v^-1
	Accuracy float64
	// in ref[1], this value is delta
	Improvement float64
	// base fixed rating
	Fixed Rating
	// system paramter tau. this value for determine next volatility.
	// in ref[1] p.1:
	// "Reasonable choices are between 0.3 and 1.2,
	// though the system should be tested to decide which value results in greatest predictive accuracy. "
	Tau float64

	//following variables tmp value for determine next sigma
	a       float64
	v       float64
	sqDelta float64
	sqPhi   float64
}

//NewEstimated is initial estimated value constractor
func NewEstimated(rating Rating, tau float64) *Estimated {
	return &Estimated{
		Accuracy:    0.0,
		Improvement: 0.0,
		Fixed:       rating,
		Tau:         tau,
	}
}

// ApplyMatch reflects match results in the training estimates.
func (e *Estimated) ApplyMatch(opponent Rating, score float64) error {
	e.Lock()
	defer e.Unlock()
	if score < 0.0 || score > 1.0 {
		return errors.New("score must be 0 to 1 (win = 1, lose = 0, draw = 0.5)")
	}
	tmp := e.Improvement * e.Accuracy
	valg := fg(opponent.phi)
	valE := fE(e.Fixed.mu, opponent.mu, opponent.phi)
	e.Accuracy += valg * valg * valE * (1.0 - valE)
	tmp += valg * (score - valE)
	e.Improvement = tmp / e.Accuracy
	return nil
}

// Fix ends the rating period and determines the new rating.
func (e *Estimated) Fix() error {
	e.Lock()
	defer e.Unlock()
	if e.Tau <= 0 {
		return errors.New("tau must be a nonzero positive number")
	}
	if e.Accuracy == 0.0 {
		// if estimated accuracy is zero, can not apply. because maybe no matches.
		// In this case, rating value and volatility parameters remain the same, but the rating deviation increases
		e.Fixed.phi = geometricMean(e.Fixed.phi, e.Fixed.sigma)
		if e.Fixed.phi > startPhi {
			e.Fixed.phi = startPhi
		}
		return nil
	}
	sigmaDash := e.determineSigma()
	phiAsta := geometricMean(e.Fixed.phi, sigmaDash)
	phiDash := 1.0 / math.Sqrt(1.0/(math.Pow(phiAsta, 2))+e.Accuracy)
	e.Fixed.mu += math.Pow(phiDash, 2) * e.Improvement * e.Accuracy
	e.Fixed.phi = phiDash
	e.Fixed.sigma = sigmaDash
	return nil
}

func (e *Estimated) determineSigma() float64 {
	e.a = math.Log(math.Pow(e.Fixed.sigma, 2))
	A := e.a
	B := 0.0

	e.sqDelta = math.Pow(e.Improvement, 2)
	e.sqPhi = math.Pow(e.Fixed.phi, 2)
	e.v = 1.0 / e.Accuracy

	switch {
	case e.sqDelta > e.sqPhi+e.v:
		B = math.Log(e.sqDelta - e.sqPhi - e.v)
	default:
		valf := 0.0
		for k := 1; k < iterationLimit+1; k++ {
			B = (e.a - float64(k)*e.Tau)
			valf = e.fx(B)
			if valf >= 0.0 {
				break
			}
		}
	}

	valfA := e.fx(A)
	valfB := e.fx(B)

	for i := 0; i < iterationLimit; i++ {
		if math.Abs(B-A) <= epsiron {
			break
		}
		C := A + ((A-B)*valfA)/(valfB-valfA)
		valfC := e.fx(C)
		switch {
		case valfB*valfC < 0.0:
			A = B
			valfA = valfB
		default:
			valfA /= 2.0
		}
		B = C
		valfB = valfC
	}

	return math.Exp(A / 2.0)
}

func (e *Estimated) fx(x float64) float64 {
	sumVal := e.sqDelta + e.sqPhi + e.v
	diffVal := e.sqDelta - e.sqPhi - e.v
	firstTerm := (math.Exp(x) * (diffVal - math.Exp(x))) / (2 * sumVal * sumVal)
	secondTerm := (x - e.a) / (math.Pow(e.Tau, 2))
	return firstTerm - secondTerm
}

// sqrt(x**2 + y**2)
func geometricMean(x, y float64) float64 {
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

// Truncate at n decimal places
func nthFloor(x float64, ref float64) float64 {
	shift := math.Pow(10, ref)
	return math.Trunc(x*shift) / shift
}

//fE is internal function E( score | mu, mu_j, phi_j)
func fE(mu, muOppo, phiOppo float64) float64 {
	return 1.0 / (1.0 + math.Exp(-1.0*fg(phiOppo)*(mu-muOppo)))
}

//fg is internal function g(phi)
func fg(phi float64) float64 {
	return 1.0 / math.Sqrt(1.0+3.0*(phi*phi)/(math.Pi*math.Pi))
}
