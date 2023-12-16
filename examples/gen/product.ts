/// <reference path="models/product_model.ts" />
import { models as modelsProduct_model } from "./models/product_model";
export namespace product {
    export type ProductListResponse = { productList: Array<modelsProduct_model.Product>; };
    export type ProductListRequest = { cursor: modelsProduct_model.ProductCursor; };
    
    type ProductListParams = { version: string; };
    export function ProductList(authority: string, params: ProductListParams, body: ProductListRequest, init?: RequestInit) {
      return new Request(`${authority}/${params.version}/products`, {
        method: 'GET',
        body: JSON.stringify(body),
        ...init,
      });
    }
    
}
