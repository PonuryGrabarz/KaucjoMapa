import { Result } from "../result";

export type MapPointType = "butelkomat" | "sklep" | "punkt_skupu" | "inny";

export type MapPointStatus =
  | "dostepny"
  | "zapelniony"
  | "nieczynny"
  | "awaria"
  | "brak_informacji"
  | "do_weryfikacji";

export type GetMapPointsRequestOptions = {
  mapPointType?: MapPointType;
  mapPointStatus?: MapPointStatus;
  details?: boolean;
  page?: number;
  pageSize?: number;
  query?: string;
};

export class GetMapPointsRequest {
  private readonly requestType = "get-map-points";

  readonly mapPointType?: MapPointType;
  readonly mapPointStatus?: MapPointStatus;
  readonly details: boolean;
  readonly page: number;
  readonly pageSize: number;
  readonly query?: string;

  private constructor(options: GetMapPointsRequestOptions) {
    this.mapPointType = options.mapPointType;
    this.mapPointStatus = options.mapPointStatus;
    this.details = options.details ?? false;
    this.page = options.page ?? 1;
    this.pageSize = options.pageSize ?? 10;
    this.query = options.query?.trim() || undefined;
  }

  public static create(
    options: GetMapPointsRequestOptions = {},
  ): Result<GetMapPointsRequest> {
    const errors: Error[] = [];

    const page = options.page ?? 1;
    const pageSize = options.pageSize ?? 10;

    if (!Number.isInteger(page) || page < 1) {
      errors.push(
        new Error("Numer strony musi być liczbą całkowitą większą lub równą 1."),
      );
    }

    if (!Number.isInteger(pageSize) || pageSize < 1) {
      errors.push(
        new Error(
          "Liczba punktów na stronie musi być liczbą całkowitą większą lub równą 1.",
        ),
      );
    }

    if (errors.length > 0) {
      return Result.failure<GetMapPointsRequest>(errors);
    }

    return Result.success(new GetMapPointsRequest(options));
  }
}
