import type { ReturnPoint } from "@shared/api/responses";

const readWindowPoints = (): ReturnPoint[] =>
	typeof window === "undefined" ? [] : (window.initialReturnPoints ?? []);

class ReturnPointsStore {
	private serverPoints = $state.raw<ReturnPoint[]>(readWindowPoints());
	private locallyCreated = $state.raw<ReturnPoint[]>([]);

	public readonly all: ReturnPoint[] = $derived([...this.serverPoints, ...this.locallyCreated]);

	public readonly byId: Map<string, ReturnPoint> = $derived(
		new Map(this.all.map((point) => [point.id, point])),
	);

	public constructor() {
		if (typeof document === "undefined") return;

		document.addEventListener("astro:after-swap", () => {
			this.serverPoints = readWindowPoints();
			this.locallyCreated = [];
		});
	}

	public find(id: string | undefined): ReturnPoint | undefined {
		return id === undefined ? undefined : this.byId.get(id);
	}

	public add(point: ReturnPoint): void {
		this.locallyCreated = [...this.locallyCreated, point];
	}
}

export const returnPoints = new ReturnPointsStore();
