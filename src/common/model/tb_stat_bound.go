package model

// table public.bound
type TbStatBound struct {
	Code_id int

	Cp_X1        uint
	Cp_Y1        float32
	Cp_X2        uint
	Cp_Y2        float32
	Cp_X_tick    uint
	Cp_Y_minus   float32
	Cp_Y_Percent float32

	Op_X1        uint
	Op_Y1        float32
	Op_X2        uint
	Op_Y2        float32
	Op_X_tick    uint
	Op_Y_minus   float32
	Op_Y_Percent float32

	Lp_X1        uint
	Lp_Y1        float32
	Lp_X2        uint
	Lp_Y2        float32
	Lp_X_tick    uint
	Lp_Y_minus   float32
	Lp_Y_Percent float32

	Hp_X1        uint
	Hp_Y1        float32
	Hp_X2        uint
	Hp_Y2        float32
	Hp_X_tick    uint
	Hp_Y_minus   float32
	Hp_Y_Percent float32
}
