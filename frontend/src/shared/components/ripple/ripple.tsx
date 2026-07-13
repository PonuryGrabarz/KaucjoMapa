import Styles from "./ripple.module.sass";

import { useCallback, useEffect, useRef, useState } from "preact/hooks";

import type { JSX } from "preact";

export type RippleProps = {
	disabled?: boolean;
};

interface RippleInfo {
	id: number;
	x: number;
	y: number;
	size: number;
	exiting: boolean;
}

export default function Ripple({ disabled }: RippleProps): JSX.Element {
	const [ripples, setRipples] = useState<RippleInfo[]>([]);

	const rippleRef = useRef<HTMLSpanElement>(null);
	const ripplesRef = useRef(ripples);

	const disabledRef = useRef(disabled);

	const seqRef = useRef(0);

	const exitTimersRef = useRef(new Set<ReturnType<typeof setTimeout>>());

	ripplesRef.current = ripples;
	disabledRef.current = disabled;

	const onDown = useCallback((e: PointerEvent): void => {
		if (disabledRef.current) return;

		const ripple = rippleRef.current;

		if (!ripple) return;

		const rect = ripple.getBoundingClientRect();

		const x = e.clientX - rect.left;
		const y = e.clientY - rect.top;

		const radius = Math.hypot(
			Math.max(x, rect.width - x),
			Math.max(y, rect.height - y),
		);

		setRipples((current) => [
			...current,
			{ id: ++seqRef.current, x, y, size: radius * 2, exiting: false },
		]);
	}, []);

	const onExit = useCallback((): void => {
		const exitingIds = ripplesRef.current.filter((r) => !r.exiting).map((r) => r.id);

		if (exitingIds.length === 0) return;

		setRipples((current) => current.map((r) =>
			exitingIds.includes(r.id) ? { ...r, exiting: true } : r,
		));

		const timer = setTimeout(() => {
			exitTimersRef.current.delete(timer);
			setRipples((current) => current.filter((r) => !exitingIds.includes(r.id)));
		}, 600);

		exitTimersRef.current.add(timer);
	}, []);

	useEffect(() => {
		const target = rippleRef.current?.parentElement;

		target?.addEventListener("pointerdown", onDown);
		target?.addEventListener("pointerup", onExit);
		target?.addEventListener("pointerleave", onExit);
		target?.addEventListener("pointercancel", onExit);

		return () => {
			target?.removeEventListener("pointerdown", onDown);
			target?.removeEventListener("pointerup", onExit);
			target?.removeEventListener("pointerleave", onExit);
			target?.removeEventListener("pointercancel", onExit);

			for (const timer of exitTimersRef.current) clearTimeout(timer);

			exitTimersRef.current.clear();
		};
	}, [onDown, onExit]);

	return (
		<span ref={rippleRef} class={Styles["ripple-container"]} aria-hidden="true">
			{ripples.map((r) => (
				<span
					key={r.id}
					class={`${Styles["ripple"]} ${r.exiting ? Styles["ripple--exit"] : ""}`}
					style={`left:${r.x}px;top:${r.y}px;width:${r.size}px;height:${r.size}px`}
				/>
			))}
		</span>
	);
}
