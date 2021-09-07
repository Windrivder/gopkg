package validators

import "regexp"

var (
	// 1. China Mobile:
	//    134, 135, 136, 137, 138, 139, 150, 151, 152, 157, 158, 159, 182, 183, 184, 187, 188,
	//    178(4G), 147(Net)；
	//    172
	//
	// 2. China Unicom:
	//    130, 131, 132, 155, 156, 185, 186 ,176(4G), 145(Net), 175
	//
	// 3. China Telecom:
	//    133, 153, 180, 181, 189, 177(4G)
	//
	// 4. Satelite:
	//    1349
	//
	// 5. Virtual:
	//    170, 173
	//
	// 6. 2018:
	//    16x, 19x
	RegMobile = regexp.MustCompile(`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^16[\d]{9}$|^17[0,2,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$|^19[\d]{9}$`)

	RegPassword = regexp.MustCompile(`^[a-z0-9A-Z\p{Han}]+(_[a-z0-9A-Z\p{Han}]+)*$`)

	RegUsername = regexp.MustCompile(`^[a-zA-Z][a-z0-9A-Z]*(_[a-z0-9A-Z]+)*$`)
)
