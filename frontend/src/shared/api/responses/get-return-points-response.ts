export class Pagination {
  page!: number;
  page_size!: number;
  total_count!: number;
  total_pages!: number;
}

export const RETURN_POINT_STATUSES = [
  "dostepny",
  "nieczynny",
  "brak_informacji",
] as const;

export const RETURN_POINT_REFUND_METHODS = [
  "gotowka",
  "bon",
  "karta",
] as const;

export const RETURN_POINT_TYPES = [
  "butelkomat",
  "sklep",
  "punkt_skupu",
  "inny",
] as const;

export const RETURN_POINT_MATERIALS = [
  "pet",
  "can",
  "glass",
] as const;

export type ReturnPointStatus = (typeof RETURN_POINT_STATUSES)[number];
export type ReturnPointRefundMethod = (typeof RETURN_POINT_REFUND_METHODS)[number];
export type ReturnPointType = (typeof RETURN_POINT_TYPES)[number];
export type ReturnPointMaterial = (typeof RETURN_POINT_MATERIALS)[number];

export class OpeningHours {
  "mon-fri"?: string;
  "mon-sat"?: string;
  "mon-sun"?: string;
  sat?: string;
  sun?: string;
}

export class ReturnPoint {
  address?: string;
  bg_color!: string;
  can!: boolean;
  created_at!: string;
  current_status!: ReturnPointStatus;
  description?: string;
  distance_km?: number
  glass!: boolean;
  id!: string;
  is_verified!: boolean;
  latitude!: number;
  longitude!: number;
  name!: string;
  opening_hours!: OpeningHours;
  pet!: boolean;
  refund_method!: ReturnPointRefundMethod[];
  short_name!: string;
  text_color!: string;
  type!: ReturnPointType;
  updated_at!: string;
  user_id!: string;
}

export class GetReturnPointsResponse {
  pagination!: Pagination;
  return_points!: ReturnPoint[];
}