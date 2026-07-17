import * as requestTypes from "../requests";
import * as responseTypes from "../responses";

import { Result } from "../result";

type ApiRequest =
  | requestTypes.GetMapPointsRequest
  | requestTypes.LoginRequest
  | requestTypes.RegisterRequest;

type ApiResponse<TRequest> = Result<
  TRequest extends requestTypes.GetMapPointsRequest
    ? responseTypes.GetMapPointsResponse
    : TRequest extends requestTypes.LoginRequest
      ? responseTypes.LoginResponse
      : TRequest extends requestTypes.RegisterRequest
        ? responseTypes.RegisterResponse
        : never
>;

export default abstract class ApiClientBase {
  private readonly apiUrl: string;

  constructor(apiUrl: string) {
    this.apiUrl = apiUrl;
  }

  public async send<TRequest extends ApiRequest>(
    request: TRequest,
  ): Promise<ApiResponse<TRequest>> {
    try {
      if (request instanceof requestTypes.GetMapPointsRequest) {
        const response = await this.getMapPoints(request);

        return Result.success(response) as ApiResponse<TRequest>;
      }

      if (request instanceof requestTypes.LoginRequest) {
        const response = await this.login(request);

        return Result.success(response) as ApiResponse<TRequest>;
      }

      if (request instanceof requestTypes.RegisterRequest) {
        const response = await this.register(request);

        return Result.success(response) as ApiResponse<TRequest>;
      }

      return Result.failure(
        new Error(
          `Nieobsługiwany typ żądania: ${(request as { constructor: { name: string } }).constructor.name}`,
        ),
      ) as ApiResponse<TRequest>;
    } catch (error) {
      return Result.failure(
        error instanceof Error
          ? error
          : new Error("Wystąpił nieznany błąd podczas komunikacji z API."),
      ) as ApiResponse<TRequest>;
    }
  }

  protected getApiUrl(): string {
    return this.apiUrl;
  }

  protected abstract getMapPoints(
    request: requestTypes.GetMapPointsRequest,
  ): Promise<responseTypes.GetMapPointsResponse>;

  protected abstract login(
    request: requestTypes.LoginRequest,
  ): Promise<responseTypes.LoginResponse>;

  protected abstract register(
    request: requestTypes.RegisterRequest,
  ): Promise<responseTypes.RegisterResponse>;
}
