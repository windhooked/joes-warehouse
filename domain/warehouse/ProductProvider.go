package warehouse

import (
	"context"
	"time"

	"github.com/averageflow/joes-warehouse/infrastructure"
)

func GetProducts() ([]infrastructure.ProductModel, error) {
	return nil, nil
}

func AddProducts(db infrastructure.ApplicationDatabase, products []infrastructure.RawProductModel) error {
	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	for i := range products {
		var productID int

		err := tx.QueryRow(
			ctx,
			addProductsQuery,
			products[i].Name,
			0,
			now,
			now,
		).Scan(&productID)

		if err != nil {
			return err
		}

		err = AddArticleProductRelation(db, productID, products[i].Articles)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func ModifyProduct(product infrastructure.ProductModel) error {
	return nil
}

func DeleteProduct(product infrastructure.ProductModel) error {
	return nil
}

func ConvertRawProduct(products []infrastructure.RawProductModel) []infrastructure.ProductModel {
	result := make([]infrastructure.ProductModel, len(products))

	for i := range products {
		result[i] = infrastructure.ProductModel{
			Name: products[i].Name,
		}
	}

	return result
}