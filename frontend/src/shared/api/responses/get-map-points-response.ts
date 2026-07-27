export class MapPoint {
  returnPointId!: string;
  shortName!: string;
  bgColor!: string;
  textColor!: string;
  latitude!: number;
  longitude!: number;
}

export class MapPointDetails extends MapPoint {
  name!: string;
  address!: string;
  pet!: boolean;
  can!: boolean;
  glass!: boolean;
  openingHours!: string;
  description!: string;
  currentStatus!:
    | "dostepny"
    | "zapelniony"
    | "nieczynny"
    | "awaria"
    | "brak_informacji"
    | "do_weryfikacji";
}

export class GetMapPointsResponse {
  points!: MapPoint[] | MapPointDetails[];
}
