package mappers

import "companiesHandler/proto/pb"

const (
	Corporations       = "Corporations"
	NonProfit          = "NonProfit"
	Cooperative        = "Cooperative"
	SoleProprietorship = "Sole Proprietorship"
)

var TypesToString = map[pb.Types]string{
	pb.Types_Corporations:   Corporations,
	pb.Types_NonProfit:      NonProfit,
	pb.Types_Cooperative:    Cooperative,
	pb.Types_Proprietorship: SoleProprietorship,
}

var StringToTypes = map[string]pb.Types{
	Corporations:       pb.Types_Corporations,
	NonProfit:          pb.Types_NonProfit,
	Cooperative:        pb.Types_Cooperative,
	SoleProprietorship: pb.Types_Proprietorship,
}
