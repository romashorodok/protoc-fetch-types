/// <reference path="../../data/product/metadata.ts" />
import { product as dataProductMetadata } from "./../../data/product/metadata";
export namespace metadata {
    export type NullRequest = { nullTextSample: string; };
    
    type ProductListParams = { version: string; };
    export function ProductList(authority: string, params: ProductListParams, body: NullRequest, init?: RequestInit) {
      return new Request(`${authority}/${params.version}/data/metadata/null`, {
        method: 'GET',
        body: JSON.stringify(body),
        ...init,
      });
    }
    
}
