import { Result } from "../result";

export class GetMapPointsRequest {
  private readonly requestType = "get-map-points";

  private constructor() {}

  public static create(): Result<GetMapPointsRequest> {
    return Result.success(new GetMapPointsRequest());
  }
}
