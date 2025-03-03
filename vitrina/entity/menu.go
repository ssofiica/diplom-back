package entity

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
	Items        []FoodDTO `json:"items"`
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
	CategoryID   uint64
	RestaurantID uint64
}

type FoodDTO struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Price        uint16 `json:"price"`
	Weight       uint16 `json:"weight"`
	Img          string `json:"img_url"`
	CategoryID   uint64 `json:"category_id,omitempty"`
	RestaurantID uint64 `json:"restaurant_id,omitempty"`
}

func (dto *FoodDTO) ToFood() Food {
	return Food{
		ID:           dto.ID,
		Name:         dto.Name,
		Price:        dto.Price,
		Weight:       dto.Weight,
		Img:          dto.Img,
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
