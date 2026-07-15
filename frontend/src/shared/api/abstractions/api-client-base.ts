import * as requestTypes from "../requests";
import * as responseTypes from "../responses";

type ApiRequest = | requestTypes.GetMapPointsRequest;

type ApiResponse<TRequest> =
  TRequest extends requestTypes.GetMapPointsRequest ? responseTypes.GetMapPointsResponse : never;

export default abstract class ApiClientBase {
  private readonly apiUrl: string;

  constructor(apiUrl: string) {
    this.apiUrl = apiUrl;
  }

  public send<TRequest extends ApiRequest>(request: TRequest): Promise<ApiResponse<TRequest>> {
    if (request instanceof requestTypes.GetMapPointsRequest) {
      return this.getMapPoints(request) as Promise<ApiResponse<TRequest>>;
    }

    throw new Error(`Unsupported request type: ${(request as { constructor: { name: string } }).constructor.name}`);
  }

  protected getApiUrl(): string {
    return this.apiUrl;
  }

  protected abstract getMapPoints(request: requestTypes.GetMapPointsRequest): Promise<responseTypes.GetMapPointsResponse>
}
