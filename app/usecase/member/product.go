package usecase_member

// func (u *appUsecase) GetProductList(ctx context.Context, claim domain.JWTClaimUser, paramQuery url.Values) response.Base {
// 	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
// 	defer cancel()

// 	page, limit, offset := yurekahelpers.GetLimitOffset(paramQuery)

// 	fetchOptions := map[string]interface{}{
// 		"limit":  limit,
// 		"offset": offset,
// 	}

// 	// filtering
// 	if paramQuery.Get("sort") != "" {
// 		fetchOptions["sort"] = paramQuery.Get("sort")
// 	}

// 	if paramQuery.Get("dir") != "" {
// 		fetchOptions["dir"] = paramQuery.Get("dir")
// 	}

// 	if paramQuery.Get("companyId") != "" {
// 		fetchOptions["companyID"] = paramQuery.Get("companyId")
// 	}

// 	// count first
// 	totalDocuments := u.mongodbRepo.CountCompanyProduct(ctx, fetchOptions)

// 	if totalDocuments == 0 {
// 		return response.Success(response.List{
// 			List:  []interface{}{},
// 			Page:  page,
// 			Limit: limit,
// 			Total: totalDocuments,
// 		})
// 	}

// 	// check ticket list
// 	cur, err := u.mongodbRepo.FetchCompanyProductList(ctx, fetchOptions)

// 	if err != nil {
// 		return response.Success(response.List{
// 			List:  []interface{}{},
// 			Page:  page,
// 			Limit: limit,
// 			Total: totalDocuments,
// 		})
// 	}

// 	defer cur.Close(ctx)

// 	list := make([]interface{}, 0)
// 	for cur.Next(ctx) {
// 		row := model.CompanyProduct{}
// 		err := cur.Decode(&row)
// 		if err != nil {
// 			logrus.Error("Topup Decode ", err)
// 			return response.Success(response.List{
// 				List:  []interface{}{},
// 				Page:  page,
// 				Limit: limit,
// 				Total: totalDocuments,
// 			})
// 		}

// 		list = append(list, row)
// 	}

// 	return response.Success(response.List{
// 		List:  list,
// 		Page:  page,
// 		Limit: limit,
// 		Total: totalDocuments,
// 	})
// }
