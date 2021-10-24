package warehouse

import (
	"context"
	"time"

	"github.com/averageflow/joes-warehouse/infrastructure"
)

func GetFullProductResponse(db infrastructure.ApplicationDatabase) (map[int64]infrastructure.WebProduct, []int64, error) {
	products, sortProducts, err := GetProducts(db)
	if err != nil {
		return nil, nil, err
	}

	productIDs := CollectProductIDs(products)

	relatedArticles, err := GetArticlesForProduct(db, productIDs)
	if err != nil {
		return nil, nil, err
	}

	for i := range relatedArticles {
		wantedProduct := products[i]
		wantedProduct.Articles = relatedArticles[i]
		wantedProduct.AmountInStock = ProductAmountInStock(wantedProduct)
		products[i] = wantedProduct
	}

	return products, sortProducts, nil
}

func GetProducts(db infrastructure.ApplicationDatabase) (map[int64]infrastructure.WebProduct, []int64, error) {
	ctx := context.Background()

	rows, err := db.Query(ctx, getProductsQuery)
	if err != nil {
		return nil, nil, err
	}

	if rows.Err() != nil {
		return nil, nil, err
	}

	defer rows.Close()

	var products []infrastructure.WebProduct

	for rows.Next() {
		var product infrastructure.WebProduct

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		products = append(products, product)
	}

	result := make(map[int64]infrastructure.WebProduct, len(products))
	orderData := make([]int64, len(products))

	for i := range products {
		result[products[i].ID] = products[i]
		orderData[i] = products[i].ID
	}

	return result, orderData, nil
}

func AddProducts(db infrastructure.ApplicationDatabase, products []infrastructure.RawProduct) error {
	ctx := context.Background()

	now := time.Now().Unix()

	articleMap := make(map[int][]infrastructure.ArticleProductRelation)

	for i := range products {
		tx, err := db.Begin(ctx)
		if err != nil {
			return err
		}

		var productID int

		err = tx.QueryRow(
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

		err = tx.Commit(ctx)
		if err != nil {
			return err
		}

		articleMap[productID] = ConvertRawArticleFromProductFile(products[i].Articles)
	}

	for i := range articleMap {
		if err := AddArticleProductRelation(db, i, articleMap[i]); err != nil {
			return err
		}
	}

	return nil
}

func SellProducts(db infrastructure.ApplicationDatabase, products map[string]int) error {
	return nil
}

func ModifyProduct(product infrastructure.Product) error {
	return nil
}

func DeleteProduct(product infrastructure.Product) error {
	return nil
}
