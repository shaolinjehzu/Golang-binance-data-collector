package klines

type dataKline struct {
	minTime int64
	maxTime int64
	isFirst bool
	lastKey int
}

type Klines struct {
	T    int64
	O    float64
	C    float64
	SMIN float64
	SMAX float64
	H    float64
	L    float64
	V    float64
	Q    float64
	QS   float64
	QB   float64
	N    int32
	NS   int32
	NB   int32
	Vt   float64
	Vm   float64
}

type trades struct{
	o	float64
}

