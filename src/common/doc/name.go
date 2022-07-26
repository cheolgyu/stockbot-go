package doc

import (
	"fmt"

	"github.com/cheolgyu/stockbot/src/common/model"
)

const DB_PUB = "test"
const DB_PUB_COLL_COMPANY = "company"
const DB_PUB_COLL_MARKET = "market"
const DB_PUB_COLL_NOTE = "note"

type NOTE_FIELD string

const PRICE_UPDATE NOTE_FIELD = "price_updated_date"
const COMPANY_UPDATE NOTE_FIELD = "company_updated_date"

func GetNoteField(country model.Country, field NOTE_FIELD) string {
	return fmt.Sprintf("%v_%v", country, field)
}

const DB_DATA = "data"
const DB_DATA_COLL_PRICE = "price"
const DB_DATA_COLL_PRICE_OPENING = "opening"
const DB_DATA_COLL_BOUND_POINT = "bound"
const DB_DATA_COLL_YMXB = "ymxb"
const DB_DATA_COLL_YMXB_QUOTE_UNIT = "quote_unit"
const DB_DATA_COLL_AGG_VOL = "agg_vol"
const DB_DATA_COLL_AGG_VOL_SUM = "agg_vol_sum"
const DB_DATA_COLL_LOGS = "logs"
