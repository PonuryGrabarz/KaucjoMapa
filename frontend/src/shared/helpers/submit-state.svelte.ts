interface ActionResult<TValue> {
	data: TValue | undefined;
	error: { message: string } | undefined;
}

export class SubmitState {
	public pending = $state(false);
	public error = $state<string | undefined>(undefined);

	private token = 0;

	public cancel(): void {
		this.token += 1;

		this.pending = false;
		this.error = undefined;
	}

	public fail(message: string): void {
		this.error = message;
	}

	public async run<TValue>(
		action: () => Promise<ActionResult<TValue>>,
		fallbackMessage: string,
	): Promise<TValue | undefined> {
		this.error = undefined;
		this.pending = true;

		const token = (this.token += 1);

		const { data, error } = await action();

		if (token !== this.token) return undefined;

		this.pending = false;

		if (error || data === undefined) {
			this.error = error?.message ?? fallbackMessage;

			return undefined;
		}

		return data;
	}
}
