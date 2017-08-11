package router

import (
	"net/http"
	"goBlog/model"
)

//地理位置/省/市/区
type LocationHandler struct {
	BaseHandler
}

func (h *LocationHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//location/province/city/district
func (*LocationHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	province := r.FormValue("province")

	if province == "" { //所有的省
		Result(w, r, model.GetProvinces())
	} else {
		city := r.FormValue("city")
		//district := r.FormValue("district")
		if city == "" { //本省所有的市
			Result(w, r, model.GetCitys(province))
		} else { //本市所有的区
			Result(w, r, model.GetDistrics(province, city))
		}
	}
}
