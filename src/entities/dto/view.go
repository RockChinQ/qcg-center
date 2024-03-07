package dto

type ValueAmount struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type AggregationValueAmountDTO struct {
	TotalAmount int           `json:"total_amount"`
	Data        []ValueAmount `json:"data"`
}
