export const distanceKm = (from: [number, number], to: [number, number]): number => {
	const radians = (degrees: number) => degrees * (Math.PI / 180);

	const earthRadius = 6371;

	const latitudeDelta = radians(to[1] - from[1]);
	const longitudeDelta = radians(to[0] - from[0]);

	const a =
		Math.sin(latitudeDelta / 2) ** 2 +
		Math.cos(radians(from[1])) * Math.cos(radians(to[1])) * Math.sin(longitudeDelta / 2) ** 2;
    
	return earthRadius * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}
