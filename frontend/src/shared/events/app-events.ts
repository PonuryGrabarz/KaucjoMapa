import Emittery from 'emittery';

type AppEvents = {
  "addLocation:click": undefined;
  "map:request_location": undefined;
  "map:begin_location_pick": undefined;
  "map:location_picked": { longitude: number; latitude: number };
  "map:show_point_by_id": { id: string };
  "map:center_on_location": { lat: number; lng: number };
};

export const events = new Emittery<AppEvents>();