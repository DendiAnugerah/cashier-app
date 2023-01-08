package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return CartRepository{db}
}

func (c *CartRepository) ReadCart() ([]model.JoinCart, error) {
	var joinCart []model.JoinCart
	err := c.db.Table("carts").Select("carts.id, carts.product_id, products.name, carts.quantity, carts.total_price").Joins("left join products on carts.product_id = products.id").Scan(&joinCart).Error
	if err != nil {
		return nil, err
	}
	return joinCart, nil
}

func (c *CartRepository) AddCart(product model.Product) error {
	var Cart model.Cart
	var Product model.Product
	err := c.db.Where("id = ?", product.ID).First(&Product).Error
	if err != nil {
		return err
	}

	err = c.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("product_id = ?", product.ID).First(&Cart).Error
		if err != nil {
			Cart.ProductID = product.ID
			Cart.Quantity = 1
			Cart.TotalPrice = product.Price - (product.Price * product.Discount / 100)
			c.db.Create(&Cart)
		} else {
			Cart.Quantity += 1
			Cart.TotalPrice += product.Price - (product.Price * product.Discount / 100)
			c.db.Save(&Cart)
		}

		err = tx.Model(&Product).Where("id = ?", product.ID).UpdateColumn("stock", gorm.Expr("stock - ?", 1)).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (c *CartRepository) DeleteCart(id uint, productID uint) error {
	var cart model.Cart
	err := c.db.Where("product_id = ?", productID).First(&cart).Error
	if err != nil {
		return err
	}

	err = c.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Delete(&cart).Error
		if err != nil {
			return err
		}

		err = tx.Model(&model.Product{}).Where("id = ?", productID).UpdateColumn("stock", gorm.Expr("stock + ?", cart.Quantity)).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (c *CartRepository) UpdateCart(id uint, cart model.Cart) error {
	var oldCart model.Cart
	err := c.db.Where("id = ?", id).First(&oldCart).Error
	if err != nil {
		return err
	}

	oldCart.ProductID = cart.ProductID
	oldCart.Quantity = cart.Quantity
	oldCart.TotalPrice = cart.TotalPrice

	result := c.db.Save(&oldCart)
	return result.Error
}
