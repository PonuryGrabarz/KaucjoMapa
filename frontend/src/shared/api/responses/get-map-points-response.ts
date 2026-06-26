export class MapPoint {
  returnPointId!: string;
  shortName!: string;
  bgColor!: string;
  textColor!: string;
  latitude!: number;
  longitude!: number;
}

export class GetMapPointsResponse {
  points!: MapPoint[];
}
