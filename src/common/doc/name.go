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

const FETCH_PRICE_UPDATE NOTE_FIELD = "fetch_price_updated"
const FETCH_COMPANY_UPDATE NOTE_FIELD = "fetch_company_updated"
const LINE_BOUND_UPDATE NOTE_FIELD = "line_bound_updated"
const LINE_YMXB_UPDATE NOTE_FIELD = "line_ymxb_updated"
const AGG_VOL_UPDATE NOTE_FIELD = "agg_vol_updated"

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
