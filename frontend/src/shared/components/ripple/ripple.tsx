import Styles from "./ripple.module.sass";

import { useCallback, useEffect, useRef, useState } from "preact/hooks";

type RippleProps = {
	disabled?: boolean;
	class?: string;
};

interface RippleInfo {
	id: number;
	x: number;
	y: number;
	size: number;
	exiting: boolean;
}

const EXIT_FALLBACK_MS = 600;
const KEYBOARD_SOURCE = "keyboard";

type RippleSource = number | typeof KEYBOARD_SOURCE;

/** The parent element must establish a containing block, usually with position: relative. */
export default function Ripple({ disabled, class: className }: RippleProps) {
	const [ripples, setRipples] = useState<RippleInfo[]>([]);

	const rippleRef = useRef<HTMLSpanElement>(null);

	const disabledRef = useRef(disabled);

	const seqRef = useRef(0);

	const activeRipplesRef = useRef(new Map<RippleSource, number>());

	const exitTimersRef = useRef(new Map<number, ReturnType<typeof setTimeout>>());

	disabledRef.current = disabled;

	const addRipple = useCallback((x: number, y: number, width: number, height: number): number => {
		const radius = Math.hypot(
			Math.max(x, width - x),
			Math.max(y, height - y),
		);
		const id = ++seqRef.current;

		setRipples((current) => [
			...current,
			{ id, x, y, size: radius * 2, exiting: false },
		]);

		return id;
	}, []);

	const removeRipple = useCallback((id: number): void => {
		const timer = exitTimersRef.current.get(id);

		if (timer) clearTimeout(timer);

		exitTimersRef.current.delete(id);

		setRipples((current) => current.filter((r) => r.id !== id));
	}, []);

	const exitRipple = useCallback((id: number): void => {
		setRipples((current) => current.map((r) =>
			r.id === id && !r.exiting ? { ...r, exiting: true } : r,
		));

		if (exitTimersRef.current.has(id)) return;

		const timer = setTimeout(() => {
			exitTimersRef.current.delete(id);
			setRipples((current) => current.filter((r) => r.id !== id));
		}, EXIT_FALLBACK_MS);

		exitTimersRef.current.set(id, timer);
	}, []);

	const onPointerDown = useCallback((e: PointerEvent): void => {
		if (disabledRef.current || activeRipplesRef.current.has(e.pointerId)) return;

		const rect = rippleRef.current?.getBoundingClientRect();

		if (!rect) return;

		const id = addRipple(
			e.clientX - rect.left,
			e.clientY - rect.top,
			rect.width,
			rect.height,
		);

		activeRipplesRef.current.set(e.pointerId, id);
	}, [addRipple]);

	const exitActiveRipple = useCallback((source: RippleSource): void => {
		const id = activeRipplesRef.current.get(source);

		if (id === undefined) return;

		activeRipplesRef.current.delete(source);

		exitRipple(id);
	}, [exitRipple]);

	const onPointerExit = useCallback((e: PointerEvent): void => {
		exitActiveRipple(e.pointerId);
	}, [exitActiveRipple]);

	const onKeyDown = useCallback((e: KeyboardEvent): void => {
		if (
			disabledRef.current
			|| e.repeat
			|| activeRipplesRef.current.has(KEYBOARD_SOURCE)
			|| (e.key !== "Enter" && e.key !== " ")
		) return;

		const rect = rippleRef.current?.getBoundingClientRect();

		if (!rect) return;

		const id = addRipple(rect.width / 2, rect.height / 2, rect.width, rect.height);

		activeRipplesRef.current.set(KEYBOARD_SOURCE, id);
	}, [addRipple]);

	const exitKeyboardRipple = useCallback((): void => {
		exitActiveRipple(KEYBOARD_SOURCE);
	}, [exitActiveRipple]);

	const onKeyUp = useCallback((e: KeyboardEvent): void => {
		if (e.key === "Enter" || e.key === " ") exitKeyboardRipple();
	}, [exitKeyboardRipple]);

	const onTransitionEnd = useCallback((e: TransitionEvent, id: number): void => {
		if (e.propertyName === "opacity" && e.target === e.currentTarget) removeRipple(id);
	}, [removeRipple]);

	useEffect(() => {
		const target = rippleRef.current?.parentElement;

		const activeRipples = activeRipplesRef.current;

		const exitTimers = exitTimersRef.current;

		target?.addEventListener("pointerdown", onPointerDown);
		target?.addEventListener("pointerup", onPointerExit);
		target?.addEventListener("pointerleave", onPointerExit);
		target?.addEventListener("pointercancel", onPointerExit);
		target?.addEventListener("keydown", onKeyDown);
		target?.addEventListener("keyup", onKeyUp);
		target?.addEventListener("blur", exitKeyboardRipple);

		return () => {
			target?.removeEventListener("pointerdown", onPointerDown);
			target?.removeEventListener("pointerup", onPointerExit);
			target?.removeEventListener("pointerleave", onPointerExit);
			target?.removeEventListener("pointercancel", onPointerExit);
			target?.removeEventListener("keydown", onKeyDown);
			target?.removeEventListener("keyup", onKeyUp);
			target?.removeEventListener("blur", exitKeyboardRipple);

			activeRipples.clear();

			for (const timer of exitTimers.values()) clearTimeout(timer);

			exitTimers.clear();
		};
	}, [exitKeyboardRipple, onKeyDown, onKeyUp, onPointerDown, onPointerExit]);

	useEffect(() => {
		if (!disabled) return;

		const ids = new Set(activeRipplesRef.current.values());

		activeRipplesRef.current.clear();

		for (const id of ids) exitRipple(id);
	}, [disabled, exitRipple]);

	const rippleClass = [Styles["ripple-container"], className].filter(Boolean).join(" ");

	return (
		<span ref={rippleRef} class={rippleClass} aria-hidden="true" data-ripple-container>
			{ripples.map((r) => (
				<span
					key={r.id}
					class={`${Styles["ripple"]} ${r.exiting ? Styles["ripple--exit"] : ""}`}
					style={`left:${r.x}px;top:${r.y}px;width:${r.size}px;height:${r.size}px`}
					onTransitionEnd={(e) => onTransitionEnd(e, r.id)}
					data-ripple
				/>
			))}
		</span>
	);
}
