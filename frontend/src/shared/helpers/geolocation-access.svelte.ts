export type GeolocationPermission = "unknown" | "prompt" | "granted" | "denied" | "unsupported";

export interface GeolocationPosition2D {
	lat: number;
	lon: number;
}

const POSITION_OPTIONS: PositionOptions = { maximumAge: 300_000, timeout: 8000 };

class GeolocationAccess {
	public permission = $state<GeolocationPermission>("unknown");
	public position = $state<GeolocationPosition2D | undefined>(undefined);

	public constructor() {
		if (typeof navigator === "undefined") return;

		void this.detect();
	}

	public get granted(): boolean {
		return this.permission === "granted";
	}

	private async detect(): Promise<void> {
		if (!("geolocation" in navigator)) {
			this.permission = "unsupported";

			return;
		}

		if (!navigator.permissions?.query) return;

		try {
			const status = await navigator.permissions.query({ name: "geolocation" });

			this.apply(status.state);

			status.addEventListener("change", () => this.apply(status.state));
		} catch {
			this.permission = "unknown";
		}
	}

	private apply(state: PermissionState): void {
		this.permission = state;

		if (state !== "granted") {
			this.position = undefined;

			return;
		}

		navigator.geolocation.getCurrentPosition(
			(reading) => {
				this.position = { lat: reading.coords.latitude, lon: reading.coords.longitude };
			},
			() => {
				this.position = undefined;
			},
			POSITION_OPTIONS,
		);
	}
}

export const geolocationAccess = new GeolocationAccess();
