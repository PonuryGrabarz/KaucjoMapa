import ky from 'ky';
import ApiClientBase from '../abstractions/api-client-base';

import type { GetMapPointsRequest } from '../requests/get-map-points-request';
import type { GetMapPointsResponse } from '../responses/get-map-points-response';

type MapPointDto = {
  return_point_id: string;
  short_name: string;
  bg_color: string;
  text_color: string;
  latitude: number;
  longitude: number;
};

export class ApiClient extends ApiClientBase {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  protected async getMapPoints(_request: GetMapPointsRequest): Promise<GetMapPointsResponse> {
    const points = await ky.get('kaucjomaty', { prefix: this.getApiUrl() }).json<MapPointDto[]>();

    return {
      points: points.map((point) => ({
        returnPointId: point.return_point_id,
        shortName: point.short_name,
        bgColor: point.bg_color,
        textColor: point.text_color,
        latitude: point.latitude,
        longitude: point.longitude,
      })),
    };
  }
}
