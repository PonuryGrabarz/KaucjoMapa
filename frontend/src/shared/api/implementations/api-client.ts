import ApiClientBase from "../abstractions/api-client-base";

import type { GetMapPointsRequest } from "../requests/get-map-points-request";
import type { GetMapPointsResponse } from "../responses/get-map-points-response";

export class ApiClient extends ApiClientBase {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  protected getMapPoints(request: GetMapPointsRequest): Promise<GetMapPointsResponse> {
    throw new Error("Method not implemented.");
  }
}
