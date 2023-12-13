/// <reference path="data/user.ts" />
import { data as dataUser } from "./data/user";
/// <reference path="data/metadata/metadata.ts" />
import { metadata as dataMetadataMetadata } from "./data/metadata/metadata";
/// <reference path="data/metadata/null.ts" />
import { metadata as dataMetadataNull } from "./data/metadata/null";
/// <reference path="data/product/metadata.ts" />
import { product as dataProductMetadata } from "./data/product/metadata";
export namespace service {
    export type ProductOwnerGetRequest = { cursor: dataUser.UserCursor; params: Array<dataMetadataMetadata.Metadata>; product_metadata: dataProductMetadata.Metadata; };
    export type ProductDeleteResponse = { ok: boolean; };
    export type Product = { title: string; price: number; };
    export type Cursor = { position: number; };
    export type ProductOwnerGetResponse = { user: dataUser.User; };
    export type ProductListResponse = { productList: Array<service.Product>; };
    export type ProductListRequest = { cursor: service.Cursor; };
    export type ProductDeleteRequest = { test: number; };
    
    type ProductListParams = { version: string; };
    export function ProductList(authority: string, params: ProductListParams, body: ProductListRequest, init?: RequestInit) {
      return new Request(`${authority}/${params.version}/products`, {
        method: 'GET',
        body: JSON.stringify(body),
        ...init,
      });
    }
    
    type ProductDeleteParams = { version: string; id: string; };
    export function ProductDelete(authority: string, params: ProductDeleteParams, body: ProductDeleteRequest, init?: RequestInit) {
      return new Request(`${authority}/${params.version}/product/${params.id}`, {
        method: 'DELETE',
        body: JSON.stringify(body),
        ...init,
      });
    }
    
    type ProductOwnerGetParams = { version: string; id: string; };
    export function ProductOwnerGet(authority: string, params: ProductOwnerGetParams, body: ProductOwnerGetRequest, init?: RequestInit) {
      return new Request(`${authority}/${params.version}/product/${params.id}/owner`, {
        method: 'GET',
        body: JSON.stringify(body),
        ...init,
      });
    }
    export type NullRequest = { nullTextSample: string; };
    
    type ProductStubGetParams = { stub: string; id: string; };
    export function ProductStubGet(authority: string, params: ProductStubGetParams, body: NullRequest, init?: RequestInit) {
      return new Request(`${authority}/${params.stub}/product/${params.id}/owner`, {
        method: 'GET',
        body: JSON.stringify(body),
        ...init,
      });
    }
    
}
