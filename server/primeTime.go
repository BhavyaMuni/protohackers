package server

type PrimeTimeServer struct {
	request  PrimeTimeRequest
	response PrimeTimeResponse
}

type PrimeTimeRequest struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

type PrimeTimeResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func ValidRequest(request PrimeTimeRequest) bool {
	if request.Method != "isPrime" || request.Number < 0 {
		return true
	}
	return false
}

func IsPrime(n float64) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= int(n); i++ {
		if int(n)%i == 0 {
			return false
		}
	}
	return true
}

func (pt PrimeTimeServer) HandleRequest() PrimeTimeResponse {
	if ValidRequest(pt.request) {
		return PrimeTimeResponse{Method: pt.request.Method, Prime: IsPrime(pt.request.Number)}
	}
	return PrimeTimeResponse{Method: pt.request.Method, Prime: false}
}
