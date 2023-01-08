package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db}
}

func (p *ProductRepository) AddProduct(product model.Product) error {
	result := p.db.Create(&product)
	return result.Error
}

func (p *ProductRepository) ReadProducts() ([]model.Product, error) {
	var product []model.Product
	err := p.db.Where("deleted_at IS NULL").Find(&product).Error

	if err != nil {
		return nil, err
	}
	return product, nil // TODO: replace this
}

func (p *ProductRepository) DeleteProduct(id uint) error {
	result := p.db.Delete(&model.Product{}, id)
	return result.Error
}

func (p *ProductRepository) UpdateProduct(id uint, product model.Product) error {
	var oldProduct model.Product
	err := p.db.Where("id = ?", id).First(&oldProduct).Error
	if err != nil {
		return err
	}

	oldProduct.Name = product.Name
	oldProduct.Price = product.Price
	oldProduct.Stock = product.Stock
	oldProduct.Discount = product.Discount
	oldProduct.Type = product.Type

	result := p.db.Save(&oldProduct)
	return result.Error
}
