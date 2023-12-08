export type ProductListResponsePackageService = { product_service: Array<ProductPackageService>; };
export type UserCursorPackageData = { direction: string; };
export type UserPackageData = { username: string; };
export type ProductOwnerGetRequestPackageService = { userCursor_data: UserCursorPackageData; metadata_metadata: Array<MetadataPackageMetadata>; metadata_product: MetadataPackageProduct; };
export type ProductListRequestPackageService = { cursor_service: CursorPackageService; };
export type ProductDeleteResponsePackageService = { ok: boolean; };
export type NullRequestPackageMetadata = { nullTextSample: string; };
export type ProductDeleteRequestPackageService = { };
export type MetadataPackageMetadata = { created_at: string; };
export type MetadataPackageProduct = { count: number; };
export type ProductOwnerGetResponsePackageService = { user_data: UserPackageData; };
export type CursorPackageService = { position: number; };
export type ProductPackageService = { title: string; price: number; };

type ProductStubGetParams = { stub: string; id: string; };
export function ProductStubGet(authority: string, params: ProductStubGetParams, body: NullRequestPackageMetadata, init?: RequestInit) {
  return new Request(`${authority}/${params.stub}/product/${params.id}/owner`, {
    method: 'GET',
    body: JSON.stringify(body),
    ...init,
  });
}

type ProductListParams = { version: string; };
export function ProductList(authority: string, params: ProductListParams, body: ProductListRequestPackageService, init?: RequestInit) {
  return new Request(`${authority}/${params.version}/products`, {
    method: 'GET',
    body: JSON.stringify(body),
    ...init,
  });
}

type ProductDeleteParams = { version: string; id: string; };
export function ProductDelete(authority: string, params: ProductDeleteParams, body: ProductDeleteRequestPackageService, init?: RequestInit) {
  return new Request(`${authority}/${params.version}/product/${params.id}`, {
    method: 'DELETE',
    body: JSON.stringify(body),
    ...init,
  });
}

type ProductOwnerGetParams = { version: string; id: string; };
export function ProductOwnerGet(authority: string, params: ProductOwnerGetParams, body: ProductOwnerGetRequestPackageService, init?: RequestInit) {
  return new Request(`${authority}/${params.version}/product/${params.id}/owner`, {
    method: 'GET',
    body: JSON.stringify(body),
    ...init,
  });
}
