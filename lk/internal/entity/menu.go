package entity

type FoodStatus string

var (
	FoodStatusStop   FoodStatus = "stop"
	FoodStatusIn     FoodStatus = "in"
	FoodStatusDelete FoodStatus = "delete"
)

type Category struct {
	ID           uint64
	Name         string
	RestaurantID uint64
}

type CategoryDTO struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	RestaurantID uint64 `json:"restaurant_id"`
}

func (dto *CategoryDTO) ToCategory() Category {
	return Category{
		ID:           dto.ID,
		Name:         dto.Name,
		RestaurantID: dto.RestaurantID,
	}
}

type Food struct {
	ID           uint64
	Name         string
	Price        uint16
	Weight       uint16
	Img          string
	Status       FoodStatus
	CategoryID   uint64
	RestaurantID uint64
}

type FoodDTO struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Price        uint16 `json:"price"`
	Weight       uint16 `json:"weight"`
	Img          string `json:"img_url"`
	Status       string `json:"status"`
	CategoryID   uint64 `json:"category_id"`
	RestaurantID uint64 `json:"restaurant_id"`
}

func (f *FoodStatus) Scan(str string) {
	switch str {
	case "stop":
		*f = FoodStatusStop
	case "in":
		*f = FoodStatusIn
	case "delete":
		*f = FoodStatusDelete
	default:
		*f = FoodStatusStop
	}
}

func (dto *FoodDTO) ToFood() Food {
	var status FoodStatus
	status.Scan(dto.Status)
	return Food{
		ID:           dto.ID,
		Name:         dto.Name,
		Price:        dto.Price,
		Weight:       dto.Weight,
		Img:          dto.Img,
		Status:       status,
		CategoryID:   dto.CategoryID,
		RestaurantID: dto.RestaurantID,
	}
}

type FoodList []Food
type CategoryList []Category
