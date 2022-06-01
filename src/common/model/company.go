package model

type Company struct {
	Code   Code   `bson:"inline"`
	Market string `bson:"market"`
	Detail CompanyDetail
	State  CompanyState
}

type CompanyDetail struct {
	//Company `bson:"company"`

	Full_code                 string `bson:"full_code"`
	Full_name_kr              string `bson:"full_name_kr"`
	Full_name_eng             string `bson:"full_name_eng"`
	Listing_date              string `bson:"listing_date"`
	Market                    string `bson:"market"`
	Securities_classification string `bson:"securities_classification"`
	Department                string `bson:"department"`
	Stock_type                string `bson:"stock_type"`
	Face_value                string `bson:"face_value"`
	Listed_stocks             string `bson:"listed_stocks"`
}

type CompanyState struct {
	//Company `bson:"company"`
	Stop          bool `bson:"stop"`
	Clear         bool `bson:"clear"`
	Managed       bool `bson:"managed"`
	Ventilation   bool `bson:"ventilation"`
	Unfaithful    bool `bson:"unfaithful"`
	Low_liquidity bool `bson:"low_liquidity"`
	Lack_listed   bool `bson:"lack_listed"`
	Overheated    bool `bson:"overheated"`
	Caution       bool `bson:"caution"`
	Warning       bool `bson:"warning"`
	Risk          bool `bson:"risk"`
}
