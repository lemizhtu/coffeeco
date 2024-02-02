package payment

type Means string

const (
	MeansCard      = "card"
	MeansCash      = "cash"
	MeansCoffeeBux = "coffee_bux"
)

type CardDetails struct {
	cardToken string
}
