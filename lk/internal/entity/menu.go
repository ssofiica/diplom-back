package entity

type FoodStatus string

var (
	FoodStatusStop   FoodStatus = "stop"
	FoodStatusIn     FoodStatus = "in"
	FoodStatusDelete FoodStatus = "delete"
)

func IsFoodStatus(status string) bool {
	return status == "stop" || status == "in" || status == "delete"
}

func (f *FoodStatus) Scan(str string) {
	switch str {
	case "stop":
		*f = FoodStatusStop
	case "in", "":
		*f = FoodStatusIn
	case "delete":
		*f = FoodStatusDelete
	default:
		*f = FoodStatusStop
	}
}

type ChangeStatusRequest struct {
	Status string `json:"status"`
}

type Category struct {
	ID           uint64
	Name         string
	RestaurantID uint64
	Items        FoodList
}

type CategoryDTO struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	RestaurantID uint64    `json:"restaurant_id,omitempty"`
	Items        []FoodDTO `json:"items,omitempty"`
}

func (dto *CategoryDTO) ToCategory() Category {
	return Category{
		ID:           dto.ID,
		Name:         dto.Name,
		RestaurantID: dto.RestaurantID,
	}
}

func (c *Category) ToDTO() CategoryDTO {
	return CategoryDTO{
		ID:           c.ID,
		Name:         c.Name,
		RestaurantID: c.RestaurantID,
		Items:        c.Items.ToDTO(),
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
	CategoryName string
	RestaurantID uint64
}

type FoodDTO struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Price        uint16 `json:"price"`
	Weight       uint16 `json:"weight"`
	Img          string `json:"img_url"`
	Status       string `json:"status"`
	CategoryID   uint64 `json:"category_id,omitempty"`
	RestaurantID uint64 `json:"restaurant_id,omitempty"`
}

type EditFood struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	Price      uint16 `json:"price"`
	Weight     uint16 `json:"weight"`
	Img        string `json:"img_url"`
	Status     string `json:"status"`
	CategoryID uint64 `json:"category_id,omitempty"`
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

func (food *Food) ToDTO() FoodDTO {
	return FoodDTO{
		ID:           food.ID,
		Name:         food.Name,
		Price:        food.Price,
		Weight:       food.Weight,
		Img:          food.Img,
		Status:       string(food.Status),
		CategoryID:   food.CategoryID,
		RestaurantID: food.RestaurantID,
	}
}

type FoodList []Food
type CategoryList []Category

func (list *FoodList) ToDTO() []FoodDTO {
	res := make([]FoodDTO, len(*list))
	for i, food := range *list {
		res[i] = food.ToDTO()
	}
	return res
}
